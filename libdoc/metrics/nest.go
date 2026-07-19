package metrics

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// EnvMetricsNestSink is the path where nested RunTest invocations append
// nest-scoped phase timing events for the outer metrics run to merge.
// Set by the outer recorder process; inherited by suite go test children.
const EnvMetricsNestSink = "DOCTEST_METRICS_NEST_SINK"

// EnvMetricsParentLeaf is set by the unified suite wrapper around each leaf
// (stdlib os.Setenv only — suite must not import this package).
const EnvMetricsParentLeaf = "DOCTEST_METRICS_PARENT_LEAF"

// parentLeaf holds an optional in-process parent leaf (tests / non-suite callers).
var parentLeaf atomic.Pointer[string]

// SetParentLeaf records which leaf is running (in-process). Prefer env
// DOCTEST_METRICS_PARENT_LEAF from the unified suite wrapper for nest timing.
// Clear with ClearParentLeaf in a defer.
func SetParentLeaf(path string) {
	if path == "" {
		parentLeaf.Store(nil)
		return
	}
	p := path
	parentLeaf.Store(&p)
}

// ClearParentLeaf clears the in-process parent leaf path.
func ClearParentLeaf() {
	parentLeaf.Store(nil)
}

// ParentLeaf returns the current suite leaf path, or "" if unset.
// Prefers DOCTEST_METRICS_PARENT_LEAF env (set by generated suite), then
// in-process SetParentLeaf.
func ParentLeaf() string {
	if v := os.Getenv(EnvMetricsParentLeaf); v != "" {
		return v
	}
	p := parentLeaf.Load()
	if p == nil {
		return ""
	}
	return *p
}

// NestSinkPath returns DOCTEST_METRICS_NEST_SINK, or "".
func NestSinkPath() string {
	return os.Getenv(EnvMetricsNestSink)
}

// NestPhaseEvent is one nest-scoped phase span written to the nest sink.
type NestPhaseEvent struct {
	Type          string         `json:"type"`
	SchemaVersion int            `json:"schema_version"`
	Scope         string         `json:"scope"` // always "nested"
	Phase         string         `json:"phase"`
	ParentLeaf    string         `json:"parent_leaf,omitempty"`
	Tree          string         `json:"tree,omitempty"`
	ElapsedNs     int64          `json:"elapsed_ns"`
	TsEnd         string         `json:"ts_end"`
	TsStart       string         `json:"ts_start,omitempty"`
	NestDepth     int            `json:"nest_depth,omitempty"`
	Detail        map[string]any `json:"detail,omitempty"`
}

var nestSinkMu sync.Mutex

// AppendNestPhase writes one nest-scoped phase event to the nest sink file.
// Safe for concurrent nested leaves within one suite process.
func AppendNestPhase(sink, phase, parentLeaf, tree string, elapsedNs int64, detail map[string]any) error {
	if sink == "" || phase == "" {
		return nil
	}
	end := time.Now().UTC()
	ev := NestPhaseEvent{
		Type:          "phase",
		SchemaVersion: SchemaVersion,
		Scope:         "nested",
		Phase:         phase,
		ParentLeaf:    parentLeaf,
		Tree:          tree,
		ElapsedNs:     elapsedNs,
		TsEnd:         end.Format(time.RFC3339Nano),
		NestDepth:     1,
		Detail:        detail,
	}
	if elapsedNs > 0 {
		ev.TsStart = end.Add(-time.Duration(elapsedNs)).Format(time.RFC3339Nano)
	}
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	data = append(data, '\n')

	nestSinkMu.Lock()
	defer nestSinkMu.Unlock()
	f, err := os.OpenFile(sink, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	_, werr := f.Write(data)
	cerr := f.Close()
	if werr != nil {
		return werr
	}
	return cerr
}

// AppendNestPhases writes all phase timings for one nested pipeline.
func AppendNestPhases(sink, parentLeaf, tree string, phases []struct {
	Name      string
	ElapsedNs int64
}, cases int) {
	if sink == "" {
		return
	}
	detail := map[string]any{}
	if cases > 0 {
		detail["cases"] = cases
	}
	for _, p := range phases {
		if p.Name == "" {
			continue
		}
		_ = AppendNestPhase(sink, p.Name, parentLeaf, tree, p.ElapsedNs, detail)
	}
}

// ReadNestSinkEvents reads JSONL objects from a nest sink file.
// Missing file returns nil, nil.
func ReadNestSinkEvents(path string) ([]map[string]any, error) {
	if path == "" {
		return nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var out []map[string]any
	sc := bufio.NewScanner(f)
	// Nested phase lines are small; raise limit for safety.
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal(line, &m); err != nil {
			continue
		}
		out = append(out, m)
	}
	return out, sc.Err()
}
