package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/metrics"
)

// EnvMetricsRoot overrides the metrics cache root when set (tests inject
// MetricsRoot via this env for the CLI path).
const EnvMetricsRoot = "DOCTEST_METRICS_ROOT"

// RunTest runs doctest tests for a single tree root with metrics recording
// and optional default-suite slow WARNING on stderr.
func RunTest(dir string, opts core.Options) error {
	if opts.Stderr == nil {
		opts.Stderr = os.Stderr
	}
	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}

	start := time.Now()
	defaultSuite := !opts.LabelAll && len(opts.LabelExprs) == 0

	rec, err := openRunRecorder(dir, opts)
	if err != nil {
		// Metrics failures are non-fatal for the suite; log and continue.
		fmt.Fprintf(opts.Stderr, "doctest: metrics: %v\n", err)
		rec = nil
	}
	if rec != nil {
		defer rec.close()
		_ = rec.writeRunStart(dir, opts, defaultSuite)
	}

	// Prefer leaf-level discovery paths for metrics before running.
	cases, discoverErr := core.DiscoverTreeCases(dir)
	if discoverErr == nil && opts.SubDir != "" {
		cases = core.FilterBySubDir(cases, dir, opts.SubDir)
	}
	if discoverErr == nil {
		cases, _ = core.FilterCasesByLabel(cases, opts)
	}

	leafStarts := time.Now()
	if rec != nil && discoverErr == nil {
		for _, c := range cases {
			_ = rec.writeLeafStart(c, dir)
		}
	}

	stats, runErr := build.TestWithStats(dir, opts)
	elapsed := time.Since(start)
	if stats.Elapsed == 0 {
		stats.Elapsed = elapsed
	}

	if rec != nil {
		// Map results: pass first stats.Passed cases, fail the rest of Total.
		// Skipped leaves get leaf_end-only (no prior start when not in cases).
		n := stats.Total
		if n == 0 && len(cases) > 0 {
			n = len(cases)
		}
		passLeft := stats.Passed
		for i, c := range cases {
			if i >= n && n > 0 {
				break
			}
			result := "fail"
			if passLeft > 0 {
				result = "pass"
				passLeft--
			}
			_ = rec.writeLeafEnd(c, dir, leafStarts, time.Now(), result, false)
		}
		// Extra skipped (label-filtered) as leaf_end skip only.
		for _, sk := range stats.Skipped {
			_ = rec.writeLeafEndSkipped(sk, dir, time.Now())
		}
		warns := []string{}
		if metrics.ShouldWarnDefaultSuiteSlow(defaultSuite, stats.Total, stats.Elapsed, metrics.DefaultSuiteWarnThreshold) {
			warns = append(warns, "default_suite_slow")
		}
		_ = rec.writeRunEnd(stats, runErr == nil && stats.Passed >= stats.Total && stats.Total > 0, warns)
		_ = rec.close()
		rec = nil // avoid double-close from defer
	}

	if metrics.ShouldWarnDefaultSuiteSlow(defaultSuite, stats.Total, stats.Elapsed, metrics.DefaultSuiteWarnThreshold) {
		msg := metrics.FormatDefaultSuiteSlowWarning()
		// Already ends with \n
		fmt.Fprint(opts.Stderr, msg)
	}

	if runErr != nil {
		return runErr
	}
	if stats.Total == 0 {
		if stats.NoTestsChanged || len(stats.Skipped) > 0 {
			return nil
		}
		return fmt.Errorf("%s: no runnable test cases found", dir)
	}
	if stats.Passed < stats.Total {
		return fmt.Errorf("%d of %d tests passed", stats.Passed, stats.Total)
	}
	return nil
}

type runRecorder struct {
	w      *metrics.Writer
	path   string
	closed bool
}

func openRunRecorder(dir string, opts core.Options) (*runRecorder, error) {
	if opts.NoMetrics {
		return nil, nil
	}
	cacheDir := resolveMetricsRoot(opts)
	if cacheDir == "" {
		return nil, fmt.Errorf("empty metrics root")
	}
	absRoot, err := filepath.Abs(dir)
	if err != nil {
		absRoot = dir
	}
	projectID := projectIDForDir(absRoot)
	path, err := metrics.CreateRunFile(cacheDir, projectID, time.Now(), "")
	if err != nil {
		return nil, err
	}
	w, err := metrics.OpenWriter(path)
	if err != nil {
		return nil, err
	}
	return &runRecorder{w: w, path: path}, nil
}

func resolveMetricsRoot(opts core.Options) string {
	if opts.MetricsRoot != "" {
		return opts.MetricsRoot
	}
	if v := os.Getenv(EnvMetricsRoot); v != "" {
		return v
	}
	base, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return base
}

func projectIDForDir(absRoot string) string {
	origin := gitRemoteOrigin(absRoot)
	if id := metrics.ProjectIDFromOrigin(origin); id != "" {
		return id
	}
	return metrics.ProjectIDFallback(absRoot)
}

func gitRemoteOrigin(dir string) string {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func gitBranchCommit(dir string) (branch, commit string) {
	bcmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	bcmd.Dir = dir
	if out, err := bcmd.Output(); err == nil {
		branch = strings.TrimSpace(string(out))
	}
	ccmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	ccmd.Dir = dir
	if out, err := ccmd.Output(); err == nil {
		commit = strings.TrimSpace(string(out))
	}
	return branch, commit
}

func (r *runRecorder) close() error {
	if r == nil || r.closed || r.w == nil {
		return nil
	}
	r.closed = true
	return r.w.Close()
}

func (r *runRecorder) writeRunStart(dir string, opts core.Options, defaultSuite bool) error {
	absRoot, _ := filepath.Abs(dir)
	cwd, _ := os.Getwd()
	projectID := projectIDForDir(absRoot)
	branch, commit := gitBranchCommit(absRoot)
	sessionID := core.DoctestSessionIDForRun()
	ev := map[string]any{
		"type":           "run_start",
		"schema_version": metrics.SchemaVersion,
		"ts":             time.Now().UTC().Format(time.RFC3339Nano),
		"project_id":     projectID,
		"cwd":            cwd,
		"argv":           []string{"doctest", "test", dir},
		"session_id":     sessionID,
		"mode": map[string]any{
			"default_suite": defaultSuite,
			"label_all":     opts.LabelAll,
			"label_exprs":   opts.LabelExprs,
		},
	}
	if branch != "" {
		ev["git_branch"] = branch
	}
	if commit != "" {
		ev["git_commit"] = commit
	}
	return r.w.Write(ev)
}

func (r *runRecorder) writeLeafStart(c core.TreeCase, root string) error {
	rel := c.Path
	if root != "" {
		if r, err := filepath.Rel(root, c.Path); err == nil {
			rel = r
		}
	}
	return r.w.Write(map[string]any{
		"type":           "leaf_start",
		"schema_version": metrics.SchemaVersion,
		"path":           rel,
		"root":           root,
		"ts_start":       time.Now().UTC().Format(time.RFC3339Nano),
		"labels":         c.Labels,
	})
}

func (r *runRecorder) writeLeafEnd(c core.TreeCase, root string, start, end time.Time, result string, cached bool) error {
	rel := c.Path
	if root != "" {
		if rr, err := filepath.Rel(root, c.Path); err == nil {
			rel = rr
		}
	}
	return r.w.Write(map[string]any{
		"type":           "leaf_end",
		"schema_version": metrics.SchemaVersion,
		"path":           rel,
		"ts_start":       start.UTC().Format(time.RFC3339Nano),
		"ts_end":         end.UTC().Format(time.RFC3339Nano),
		"elapsed_ns":     end.Sub(start).Nanoseconds(),
		"result":         result,
		"cached":         cached,
	})
}

func (r *runRecorder) writeLeafEndSkipped(sk core.SkippedCase, root string, end time.Time) error {
	rel := sk.Path
	if root != "" {
		if rr, err := filepath.Rel(root, sk.Path); err == nil {
			rel = rr
		}
	}
	return r.w.Write(map[string]any{
		"type":           "leaf_end",
		"schema_version": metrics.SchemaVersion,
		"path":           rel,
		"ts_end":         end.UTC().Format(time.RFC3339Nano),
		"elapsed_ns":     int64(0),
		"result":         "skip",
		"cached":         false,
	})
}

func (r *runRecorder) writeRunEnd(stats build.TestRunStats, exitOK bool, warnings []string) error {
	if warnings == nil {
		warnings = []string{}
	}
	return r.w.Write(map[string]any{
		"type":           "run_end",
		"schema_version": metrics.SchemaVersion,
		"wall_ns":        stats.Elapsed.Nanoseconds(),
		"passed":         stats.Passed,
		"total":          stats.Total,
		"skipped":        len(stats.Skipped),
		"exit_ok":        exitOK,
		"warnings":       warnings,
	})
}

