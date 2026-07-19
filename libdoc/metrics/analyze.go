package metrics

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DefaultRunRetention is how many newest run files prune keeps.
const DefaultRunRetention = 30

// DefaultTopN is the default row limit for metrics top.
const DefaultTopN = 10

// DefaultSummaryLast is the default --last N for metrics summary.
const DefaultSummaryLast = 5

// Event is a loosely typed metrics JSONL event.
type Event map[string]any

// RunFile describes one runs/*.jsonl file.
type RunFile struct {
	Name string // basename including .jsonl
	Path string
	ID   string // stem without .jsonl
}

// LeafRow is a ranked leaf from leaf_end (+ labels from leaf_start when present).
type LeafRow struct {
	Path      string   `json:"path"`
	ElapsedNs int64    `json:"elapsed_ns"`
	Result    string   `json:"result,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	Cached    bool     `json:"cached,omitempty"`
}

// PhaseRow is one pipeline phase event (type=phase).
type PhaseRow struct {
	Scope      string `json:"scope"` // suite | tree | nested
	Phase      string `json:"phase"` // discover | materialize | generate | go_test | …
	Tree       string `json:"tree,omitempty"`
	ParentLeaf string `json:"parent_leaf,omitempty"` // set when scope=nested
	ElapsedNs  int64  `json:"elapsed_ns"`
}

// NestedPhaseRow ranks nested re-entry cost (scope=nested) by parent leaf.
type NestedPhaseRow struct {
	ParentLeaf string `json:"parent_leaf"`
	Phase      string `json:"phase"`
	ElapsedNs  int64  `json:"elapsed_ns"`
	Tree       string `json:"tree,omitempty"`
}

// PhaseTotals aggregates phase wall times (sum of tree walls; may exceed suite wall when parallel).
type PhaseTotals struct {
	ByPhase map[string]int64 `json:"by_phase"` // phase name -> summed elapsed_ns
	// Ordered known phases first for stable display.
	Order []string `json:"order,omitempty"`
}

// PhaseAnalysis is the result of ExtractPhases + totals for one run.
type PhaseAnalysis struct {
	RunID     string         `json:"run_id"`
	WallNs    int64          `json:"wall_ns"` // suite wall from run_end when present
	Phases    []PhaseRow     `json:"phases"`
	Totals    PhaseTotals    `json:"totals"`
	TopTrees  []PhaseTreeRow `json:"top_trees,omitempty"` // by go_test
}

// PhaseTreeRow ranks trees by one phase (typically go_test).
type PhaseTreeRow struct {
	Tree      string `json:"tree"`
	ElapsedNs int64  `json:"elapsed_ns"`
	Phase     string `json:"phase"`
}

// RunSummary is a human/machine summary of one run file.
type RunSummary struct {
	RunID         string    `json:"run_id"`
	File          string    `json:"file"`
	DefaultSuite  bool      `json:"default_suite"`
	Passed        int       `json:"passed"`
	Total         int       `json:"total"`
	Skipped       int       `json:"skipped"`
	WallNs        int64     `json:"wall_ns"`
	ExitOK        bool      `json:"exit_ok"`
	LeafCount     int       `json:"leaf_count"`
	Slowest       []LeafRow `json:"slowest,omitempty"`
	HasRunEnd     bool      `json:"has_run_end"`
	HasRunStart   bool      `json:"has_run_start"`
}

// AggregateSummary covers multiple runs.
type AggregateSummary struct {
	Runs       int      `json:"runs"`
	RunIDs     []string `json:"run_ids"`
	Passed     int      `json:"passed"`
	Total      int      `json:"total"`
	Skipped    int      `json:"skipped"`
	WallNs     int64    `json:"wall_ns"`
	LeafCount  int      `json:"leaf_count"`
	DefaultOnly bool    `json:"default_only,omitempty"`
}

// ListRunFiles returns *.jsonl basenames under runsDir sorted ascending (oldest first).
func ListRunFiles(runsDir string) ([]RunFile, error) {
	ents, err := os.ReadDir(runsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []RunFile
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".jsonl") {
			continue
		}
		out = append(out, RunFile{
			Name: name,
			Path: filepath.Join(runsDir, name),
			ID:   strings.TrimSuffix(name, ".jsonl"),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// NewestRun returns the last file in lexicographic order, or nil if empty.
func NewestRun(files []RunFile) *RunFile {
	if len(files) == 0 {
		return nil
	}
	f := files[len(files)-1]
	return &f
}

// FindRunByID matches stem or full basename (with/without .jsonl).
func FindRunByID(files []RunFile, id string) *RunFile {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil
	}
	want := strings.TrimSuffix(id, ".jsonl")
	for i := range files {
		if files[i].ID == want || files[i].Name == id {
			return &files[i]
		}
	}
	return nil
}

// ReadEvents loads newline-delimited JSON objects from a run file.
func ReadEvents(path string) ([]Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var events []Event
	sc := bufio.NewScanner(f)
	// large lines
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 4*1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var ev Event
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue // skip malformed
		}
		events = append(events, ev)
	}
	if err := sc.Err(); err != nil {
		return events, err
	}
	return events, nil
}

func eventType(ev Event) string {
	if t, ok := ev["type"].(string); ok {
		return t
	}
	return ""
}

func eventString(ev Event, key string) string {
	if v, ok := ev[key].(string); ok {
		return v
	}
	return ""
}

func eventInt64(ev Event, key string) int64 {
	v, ok := ev[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int64(n)
	case int64:
		return n
	case int:
		return int64(n)
	case json.Number:
		i, _ := n.Int64()
		return i
	default:
		return 0
	}
}

func eventBool(ev Event, key string) bool {
	v, ok := ev[key]
	if !ok || v == nil {
		return false
	}
	switch b := v.(type) {
	case bool:
		return b
	default:
		return false
	}
}

func eventLabels(ev Event) []string {
	v, ok := ev["labels"]
	if !ok || v == nil {
		return nil
	}
	switch arr := v.(type) {
	case []any:
		var out []string
		for _, x := range arr {
			if s, ok := x.(string); ok {
				out = append(out, s)
			}
		}
		return out
	case []string:
		return arr
	default:
		return nil
	}
}

// IsDefaultSuite reports mode.default_suite from run_start events.
func IsDefaultSuite(events []Event) bool {
	for _, ev := range events {
		if eventType(ev) != "run_start" {
			continue
		}
		mode, ok := ev["mode"].(map[string]any)
		if !ok {
			// missing mode: treat as not default-suite for filter purposes
			return false
		}
		if b, ok := mode["default_suite"].(bool); ok {
			return b
		}
		return false
	}
	return false
}

// ExtractLeaves builds leaf rows from leaf_end, attaching labels from leaf_start by path.
func ExtractLeaves(events []Event) []LeafRow {
	labelsByPath := map[string][]string{}
	for _, ev := range events {
		if eventType(ev) != "leaf_start" {
			continue
		}
		p := eventString(ev, "path")
		if p == "" {
			continue
		}
		labelsByPath[p] = eventLabels(ev)
	}
	var rows []LeafRow
	for _, ev := range events {
		if eventType(ev) != "leaf_end" {
			continue
		}
		p := eventString(ev, "path")
		if p == "" {
			continue
		}
		labels := labelsByPath[p]
		if labels == nil {
			labels = eventLabels(ev)
		}
		rows = append(rows, LeafRow{
			Path:      p,
			ElapsedNs: eventInt64(ev, "elapsed_ns"),
			Result:    eventString(ev, "result"),
			Labels:    labels,
			Cached:    eventBool(ev, "cached"),
		})
	}
	return rows
}

// RankLeaves sorts by elapsed_ns descending; optional unlabeled-only filter; limit n (0 = no limit / DefaultTopN when applied by caller).
func RankLeaves(rows []LeafRow, unlabeledOnly bool, n int) []LeafRow {
	var filtered []LeafRow
	for _, r := range rows {
		if unlabeledOnly && len(r.Labels) > 0 {
			continue
		}
		filtered = append(filtered, r)
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		if filtered[i].ElapsedNs == filtered[j].ElapsedNs {
			return filtered[i].Path < filtered[j].Path
		}
		return filtered[i].ElapsedNs > filtered[j].ElapsedNs
	})
	if n > 0 && len(filtered) > n {
		filtered = filtered[:n]
	}
	return filtered
}

// SummarizeRun builds a RunSummary from events and file metadata.
func SummarizeRun(rf RunFile, events []Event) RunSummary {
	s := RunSummary{
		RunID:        rf.ID,
		File:         rf.Name,
		DefaultSuite: IsDefaultSuite(events),
	}
	// Prefer run_id from run_start if present
	for _, ev := range events {
		switch eventType(ev) {
		case "run_start":
			s.HasRunStart = true
			if id := eventString(ev, "run_id"); id != "" {
				s.RunID = id
			}
		case "run_end":
			s.HasRunEnd = true
			s.Passed = int(eventInt64(ev, "passed"))
			s.Total = int(eventInt64(ev, "total"))
			s.Skipped = int(eventInt64(ev, "skipped"))
			s.WallNs = eventInt64(ev, "wall_ns")
			s.ExitOK = eventBool(ev, "exit_ok")
		}
	}
	leaves := ExtractLeaves(events)
	s.LeafCount = len(leaves)
	s.Slowest = RankLeaves(leaves, false, 5)
	return s
}

// knownPhaseOrder is display order for common pipeline phases.
var knownPhaseOrder = []string{
	"cold_cache", "resolve", "discover", "materialize", "generate", "go_test", "post",
}

// ExtractPhases collects type=phase events.
func ExtractPhases(events []Event) []PhaseRow {
	var rows []PhaseRow
	for _, ev := range events {
		if eventType(ev) != "phase" {
			continue
		}
		rows = append(rows, PhaseRow{
			Scope:      eventString(ev, "scope"),
			Phase:      eventString(ev, "phase"),
			Tree:       eventString(ev, "tree"),
			ParentLeaf: eventString(ev, "parent_leaf"),
			ElapsedNs:  eventInt64(ev, "elapsed_ns"),
		})
	}
	return rows
}

// ExtractNestedPhases returns only scope=nested phase rows.
func ExtractNestedPhases(events []Event) []PhaseRow {
	var out []PhaseRow
	for _, r := range ExtractPhases(events) {
		if r.Scope == "nested" {
			out = append(out, r)
		}
	}
	return out
}

// RankNestedByParentLeaf sums nested phase elapsed per parent_leaf (all nested
// phases), slowest first. limit n (0 = no limit). Non-nested rows are ignored.
func RankNestedByParentLeaf(rows []PhaseRow, n int) []NestedPhaseRow {
	by := map[string]int64{}
	for _, r := range rows {
		if r.Scope != "nested" {
			continue
		}
		key := r.ParentLeaf
		if key == "" {
			key = "(unknown)"
		}
		by[key] += r.ElapsedNs
	}
	var out []NestedPhaseRow
	for leaf, ns := range by {
		out = append(out, NestedPhaseRow{ParentLeaf: leaf, Phase: "all", ElapsedNs: ns})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].ElapsedNs == out[j].ElapsedNs {
			return out[i].ParentLeaf < out[j].ParentLeaf
		}
		return out[i].ElapsedNs > out[j].ElapsedNs
	})
	if n > 0 && len(out) > n {
		out = out[:n]
	}
	return out
}

// RankNestedGoTest ranks nested go_test phases by elapsed, slowest first.
func RankNestedGoTest(rows []PhaseRow, n int) []NestedPhaseRow {
	var out []NestedPhaseRow
	for _, r := range rows {
		if r.Scope != "nested" || r.Phase != "go_test" {
			continue
		}
		out = append(out, NestedPhaseRow{
			ParentLeaf: r.ParentLeaf,
			Phase:      r.Phase,
			ElapsedNs:  r.ElapsedNs,
			Tree:       r.Tree,
		})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].ElapsedNs == out[j].ElapsedNs {
			return out[i].ParentLeaf < out[j].ParentLeaf
		}
		return out[i].ElapsedNs > out[j].ElapsedNs
	})
	if n > 0 && len(out) > n {
		out = out[:n]
	}
	return out
}

// SumPhases sums elapsed_ns by phase name (all scopes).
func SumPhases(rows []PhaseRow) PhaseTotals {
	by := map[string]int64{}
	seen := map[string]bool{}
	for _, r := range rows {
		if r.Phase == "" {
			continue
		}
		by[r.Phase] += r.ElapsedNs
		seen[r.Phase] = true
	}
	var order []string
	for _, p := range knownPhaseOrder {
		if seen[p] {
			order = append(order, p)
			delete(seen, p)
		}
	}
	var rest []string
	for p := range seen {
		rest = append(rest, p)
	}
	sort.Strings(rest)
	order = append(order, rest...)
	return PhaseTotals{ByPhase: by, Order: order}
}

// TopTreesByPhase ranks trees by summed elapsed for a given phase name.
func TopTreesByPhase(rows []PhaseRow, phase string, n int) []PhaseTreeRow {
	byTree := map[string]int64{}
	for _, r := range rows {
		if r.Phase != phase || r.Tree == "" {
			continue
		}
		byTree[r.Tree] += r.ElapsedNs
	}
	var out []PhaseTreeRow
	for tree, ns := range byTree {
		out = append(out, PhaseTreeRow{Tree: tree, ElapsedNs: ns, Phase: phase})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].ElapsedNs == out[j].ElapsedNs {
			return out[i].Tree < out[j].Tree
		}
		return out[i].ElapsedNs > out[j].ElapsedNs
	})
	if n > 0 && len(out) > n {
		out = out[:n]
	}
	return out
}

// AnalyzePhases builds a full phase view for one run file.
func AnalyzePhases(rf RunFile, events []Event) PhaseAnalysis {
	rows := ExtractPhases(events)
	a := PhaseAnalysis{
		RunID:  rf.ID,
		Phases: rows,
		Totals: SumPhases(rows),
	}
	for _, ev := range events {
		if eventType(ev) == "run_end" {
			a.WallNs = eventInt64(ev, "wall_ns")
			break
		}
	}
	a.TopTrees = TopTreesByPhase(rows, "go_test", 10)
	return a
}

// SelectRun chooses a run file given --run id ("last" or empty => newest).
// When defaultOnly, prefer newest file whose run_start has default_suite true.
func SelectRun(files []RunFile, runSel string, defaultOnly bool) (*RunFile, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no runs found")
	}
	runSel = strings.TrimSpace(runSel)
	if runSel == "" || runSel == "last" {
		if !defaultOnly {
			return NewestRun(files), nil
		}
		// walk newest-first
		for i := len(files) - 1; i >= 0; i-- {
			evs, err := ReadEvents(files[i].Path)
			if err != nil {
				continue
			}
			if IsDefaultSuite(evs) {
				f := files[i]
				return &f, nil
			}
		}
		return nil, fmt.Errorf("no default-suite runs found")
	}
	rf := FindRunByID(files, runSel)
	if rf == nil {
		return nil, fmt.Errorf("run not found: %s", runSel)
	}
	if defaultOnly {
		evs, err := ReadEvents(rf.Path)
		if err != nil {
			return nil, err
		}
		if !IsDefaultSuite(evs) {
			return nil, fmt.Errorf("run %s is not a default-suite run", rf.ID)
		}
	}
	return rf, nil
}

// LastNRuns returns up to n newest files (n<=0 uses DefaultSummaryLast).
// When defaultOnly, only default-suite runs are considered (newest among those).
func LastNRuns(files []RunFile, n int, defaultOnly bool) ([]RunFile, error) {
	if n <= 0 {
		n = DefaultSummaryLast
	}
	var candidates []RunFile
	if !defaultOnly {
		candidates = files
	} else {
		for _, f := range files {
			evs, err := ReadEvents(f.Path)
			if err != nil {
				continue
			}
			if IsDefaultSuite(evs) {
				candidates = append(candidates, f)
			}
		}
	}
	if len(candidates) == 0 {
		return nil, nil
	}
	// newest last; take last n
	if len(candidates) > n {
		candidates = candidates[len(candidates)-n:]
	}
	return candidates, nil
}

// AggregateRuns sums end stats across selected runs.
func AggregateRuns(files []RunFile) (AggregateSummary, error) {
	agg := AggregateSummary{
		RunIDs: make([]string, 0, len(files)),
	}
	for _, f := range files {
		evs, err := ReadEvents(f.Path)
		if err != nil {
			return agg, err
		}
		s := SummarizeRun(f, evs)
		agg.Runs++
		agg.RunIDs = append(agg.RunIDs, s.RunID)
		agg.Passed += s.Passed
		agg.Total += s.Total
		agg.Skipped += s.Skipped
		agg.WallNs += s.WallNs
		agg.LeafCount += s.LeafCount
	}
	return agg, nil
}

// PruneRuns deletes oldest files beyond keep (newest keep retained by name sort).
// Returns number removed.
func PruneRuns(runsDir string, keep int) (int, error) {
	if keep <= 0 {
		keep = DefaultRunRetention
	}
	files, err := ListRunFiles(runsDir)
	if err != nil {
		return 0, err
	}
	if len(files) <= keep {
		return 0, nil
	}
	// remove oldest prefix
	toRemove := files[:len(files)-keep]
	removed := 0
	for _, f := range toRemove {
		if err := os.Remove(f.Path); err != nil && !os.IsNotExist(err) {
			return removed, err
		}
		removed++
	}
	return removed, nil
}

// FormatDurationNS formats nanoseconds for human display.
func FormatDurationNS(ns int64) string {
	if ns < 1_000_000 {
		return fmt.Sprintf("%dns", ns)
	}
	if ns < 1_000_000_000 {
		return fmt.Sprintf("%.2fms", float64(ns)/1e6)
	}
	return fmt.Sprintf("%.2fs", float64(ns)/1e9)
}
