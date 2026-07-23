package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xhd2015/doctest/libdoc/metrics"
)

const metricsUsage = `Usage: doctest metrics <subcommand> [options]

Analyze recorded doctest metrics for the current project.

Subcommands:
  path      Print the project metrics directory
  last      Summarize the newest run
  top       Rank slowest leaves in a run
  phases    Break down pipeline phases (discover/generate/go_test/…)
  summary   Aggregate stats across recent runs
  show      Dump one run (latest or by run id)
  prune     Delete oldest run files beyond retention (keep 30)

Environment:
  DOCTEST_METRICS_ROOT   Override metrics cache root

Options for top:
  --n N                Limit ranked rows (default 10)
  --unlabeled-only     Only leaves with empty labels
  --default-only       Prefer / use default-suite runs
  --json               Machine-readable JSON on stdout
  --run last|ID        Select run (default: last)

Options for phases:
  --run last|ID        Select run (default: last)
  --n N                Limit top trees by go_test (default 10)
  --json               Machine-readable JSON on stdout

Options for summary:
  --last N             Aggregate last N runs (default 5)
  --default-only       Only default-suite runs
  --json               Machine-readable JSON on stdout

Examples:
  doctest metrics path
  doctest metrics last
  doctest metrics top --n 5 --unlabeled-only
  doctest metrics phases --run last
  doctest metrics summary --last 2 --json
  doctest metrics show <run-id>
  doctest metrics prune
`

func runMetrics(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprint(cliStdout(), metricsUsage)
		return nil
	}
	sub := args[0]
	rest := args[1:]
	switch sub {
	case "path":
		return metricsPath(rest)
	case "last":
		return metricsLast(rest)
	case "top":
		return metricsTop(rest)
	case "phases":
		return metricsPhases(rest)
	case "summary":
		return metricsSummary(rest)
	case "show":
		return metricsShow(rest)
	case "prune":
		return metricsPrune(rest)
	case "help":
		fmt.Fprint(cliStdout(), metricsUsage)
		return nil
	default:
		return fmt.Errorf("unknown metrics subcommand: %s\n\n%s", sub, metricsUsage)
	}
}

type metricsCtx struct {
	Root      string
	ProjectID string
	MetricsDir string
	RunsDir   string
}

func loadMetricsCtx() (*metricsCtx, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	root := metrics.ResolveMetricsRoot("")
	if root == "" {
		return nil, fmt.Errorf("cannot resolve metrics root")
	}
	pid := metrics.ProjectIDForDir(cwd)
	return &metricsCtx{
		Root:       root,
		ProjectID:  pid,
		MetricsDir: metrics.ProjectMetricsDir(root, pid),
		RunsDir:    metrics.ProjectRunsDir(root, pid),
	}, nil
}

func metricsPath(args []string) error {
	if helpArgs(args) {
		fmt.Print("Usage: doctest metrics path\n\nPrint the absolute project metrics directory.\n")
		return nil
	}
	ctx, err := loadMetricsCtx()
	if err != nil {
		return err
	}
	abs, err := filepath.Abs(ctx.MetricsDir)
	if err != nil {
		abs = ctx.MetricsDir
	}
	fmt.Println(abs)
	return nil
}

func metricsLast(args []string) error {
	if helpArgs(args) {
		fmt.Print("Usage: doctest metrics last\n\nSummarize the newest metrics run file.\n")
		return nil
	}
	ctx, err := loadMetricsCtx()
	if err != nil {
		return err
	}
	files, err := metrics.ListRunFiles(ctx.RunsDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no runs found in %s", ctx.RunsDir)
	}
	rf := metrics.NewestRun(files)
	evs, err := metrics.ReadEvents(rf.Path)
	if err != nil {
		return err
	}
	s := metrics.SummarizeRun(*rf, evs)
	printRunSummaryHuman(s)
	return nil
}

func printRunSummaryHuman(s metrics.RunSummary) {
	fmt.Printf("run_id: %s\n", s.RunID)
	fmt.Printf("file: %s\n", s.File)
	fmt.Printf("default_suite: %v\n", s.DefaultSuite)
	if s.HasRunEnd {
		fmt.Printf("passed: %d  total: %d  skipped: %d\n", s.Passed, s.Total, s.Skipped)
		fmt.Printf("wall: %s\n", metrics.FormatDurationNS(s.WallNs))
		fmt.Printf("exit_ok: %v\n", s.ExitOK)
	}
	fmt.Printf("leaf_count: %d\n", s.LeafCount)
	if len(s.Slowest) > 0 {
		fmt.Println("slowest:")
		for _, row := range s.Slowest {
			fmt.Printf("  %s  %s  %s\n", row.Path, metrics.FormatDurationNS(row.ElapsedNs), row.Result)
		}
	}
}

func metricsPhases(args []string) error {
	if helpArgs(args) {
		fmt.Print(`Usage: doctest metrics phases [options]

Break down pipeline phase costs for a recorded run (discover, materialize,
generate, go_test, …). Phase totals are summed tree walls and may exceed suite
wall when trees run in parallel.

Options:
  --run last|ID
  --n N
  --json
`)
		return nil
	}
	asJSON := false
	runSel := "last"
	n := 10
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-h" || a == "--help":
			return metricsPhases([]string{"--help"})
		case a == "--json":
			asJSON = true
		case a == "--n" || a == "-n":
			if i+1 >= len(args) {
				return fmt.Errorf("--n requires a value")
			}
			i++
			v, err := strconv.Atoi(args[i])
			if err != nil || v < 0 {
				return fmt.Errorf("invalid --n: %s", args[i])
			}
			n = v
		case strings.HasPrefix(a, "--n="):
			v, err := strconv.Atoi(strings.TrimPrefix(a, "--n="))
			if err != nil || v < 0 {
				return fmt.Errorf("invalid --n: %s", a)
			}
			n = v
		case a == "--run":
			if i+1 >= len(args) {
				return fmt.Errorf("--run requires a value")
			}
			i++
			runSel = args[i]
		case strings.HasPrefix(a, "--run="):
			runSel = strings.TrimPrefix(a, "--run=")
		default:
			return fmt.Errorf("unknown flag for metrics phases: %s", a)
		}
	}

	ctx, err := loadMetricsCtx()
	if err != nil {
		return err
	}
	files, err := metrics.ListRunFiles(ctx.RunsDir)
	if err != nil {
		return err
	}
	rf, err := metrics.SelectRun(files, runSel, false)
	if err != nil {
		return err
	}
	evs, err := metrics.ReadEvents(rf.Path)
	if err != nil {
		return err
	}
	a := metrics.AnalyzePhases(*rf, evs)
	a.TopTrees = metrics.TopTreesByPhase(a.Phases, "go_test", n)
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(a)
	}
	fmt.Printf("run_id: %s\n", a.RunID)
	if a.WallNs > 0 {
		fmt.Printf("suite_wall: %s\n", metrics.FormatDurationNS(a.WallNs))
	}
	fmt.Println("phase totals (summed tree wall; may exceed suite wall when parallel):")
	if len(a.Totals.ByPhase) == 0 {
		fmt.Println("  (no phase events — re-run with --metrics-on on a binary that emits phases)")
	} else {
		for _, name := range a.Totals.Order {
			ns := a.Totals.ByPhase[name]
			fmt.Printf("  %-14s %s\n", name, metrics.FormatDurationNS(ns))
		}
	}
	if len(a.TopTrees) > 0 {
		fmt.Println("top trees by go_test:")
		for i, row := range a.TopTrees {
			fmt.Printf("  %d. %s  %s\n", i+1, row.Tree, metrics.FormatDurationNS(row.ElapsedNs))
		}
	}
	return nil
}

func metricsTop(args []string) error {
	if helpArgs(args) {
		fmt.Print(`Usage: doctest metrics top [options]

Rank leaf_end events by elapsed_ns (slowest first).

Options:
  --n N
  --unlabeled-only
  --default-only
  --json
  --run last|ID
`)
		return nil
	}
	n := metrics.DefaultTopN
	unlabeledOnly := false
	defaultOnly := false
	asJSON := false
	runSel := "last"
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-h" || a == "--help":
			return metricsTop([]string{"--help"})
		case a == "--unlabeled-only":
			unlabeledOnly = true
		case a == "--default-only":
			defaultOnly = true
		case a == "--json":
			asJSON = true
		case a == "--n" || a == "-n":
			if i+1 >= len(args) {
				return fmt.Errorf("--n requires a value")
			}
			i++
			v, err := strconv.Atoi(args[i])
			if err != nil || v < 0 {
				return fmt.Errorf("invalid --n: %s", args[i])
			}
			n = v
		case strings.HasPrefix(a, "--n="):
			v, err := strconv.Atoi(strings.TrimPrefix(a, "--n="))
			if err != nil || v < 0 {
				return fmt.Errorf("invalid --n: %s", a)
			}
			n = v
		case a == "--run":
			if i+1 >= len(args) {
				return fmt.Errorf("--run requires a value")
			}
			i++
			runSel = args[i]
		case strings.HasPrefix(a, "--run="):
			runSel = strings.TrimPrefix(a, "--run=")
		default:
			return fmt.Errorf("unknown flag for metrics top: %s", a)
		}
	}

	ctx, err := loadMetricsCtx()
	if err != nil {
		return err
	}
	files, err := metrics.ListRunFiles(ctx.RunsDir)
	if err != nil {
		return err
	}
	rf, err := metrics.SelectRun(files, runSel, defaultOnly)
	if err != nil {
		return err
	}
	evs, err := metrics.ReadEvents(rf.Path)
	if err != nil {
		return err
	}
	rows := metrics.RankLeaves(metrics.ExtractLeaves(evs), unlabeledOnly, n)
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(rows)
	}
	fmt.Printf("run_id: %s\n", rf.ID)
	for i, row := range rows {
		fmt.Printf("%d. %s  %s", i+1, row.Path, metrics.FormatDurationNS(row.ElapsedNs))
		if row.Result != "" {
			fmt.Printf("  %s", row.Result)
		}
		if len(row.Labels) > 0 {
			fmt.Printf("  labels=%s", strings.Join(row.Labels, ","))
		}
		fmt.Println()
	}
	return nil
}

func metricsSummary(args []string) error {
	if helpArgs(args) {
		fmt.Print(`Usage: doctest metrics summary [options]

Aggregate recent run files.

Options:
  --last N
  --default-only
  --json
`)
		return nil
	}
	lastN := metrics.DefaultSummaryLast
	defaultOnly := false
	asJSON := false
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-h" || a == "--help":
			return metricsSummary([]string{"--help"})
		case a == "--default-only":
			defaultOnly = true
		case a == "--json":
			asJSON = true
		case a == "--last":
			if i+1 >= len(args) {
				return fmt.Errorf("--last requires a value")
			}
			i++
			v, err := strconv.Atoi(args[i])
			if err != nil || v < 0 {
				return fmt.Errorf("invalid --last: %s", args[i])
			}
			lastN = v
		case strings.HasPrefix(a, "--last="):
			v, err := strconv.Atoi(strings.TrimPrefix(a, "--last="))
			if err != nil || v < 0 {
				return fmt.Errorf("invalid --last: %s", a)
			}
			lastN = v
		default:
			return fmt.Errorf("unknown flag for metrics summary: %s", a)
		}
	}

	ctx, err := loadMetricsCtx()
	if err != nil {
		return err
	}
	files, err := metrics.ListRunFiles(ctx.RunsDir)
	if err != nil {
		return err
	}
	selected, err := metrics.LastNRuns(files, lastN, defaultOnly)
	if err != nil {
		return err
	}
	if len(selected) == 0 {
		return fmt.Errorf("no runs found")
	}
	agg, err := metrics.AggregateRuns(selected)
	if err != nil {
		return err
	}
	agg.DefaultOnly = defaultOnly
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(agg)
	}
	fmt.Printf("runs: %d\n", agg.Runs)
	fmt.Printf("run_ids: %s\n", strings.Join(agg.RunIDs, ", "))
	fmt.Printf("passed: %d  total: %d  skipped: %d\n", agg.Passed, agg.Total, agg.Skipped)
	fmt.Printf("wall: %s\n", metrics.FormatDurationNS(agg.WallNs))
	fmt.Printf("leaf_count: %d\n", agg.LeafCount)
	if defaultOnly {
		fmt.Println("default_only: true")
	}
	return nil
}

func metricsShow(args []string) error {
	if helpArgs(args) {
		fmt.Print("Usage: doctest metrics show [run-id]\n\nDump events for one run (latest if run-id omitted).\n")
		return nil
	}
	runID := ""
	for _, a := range args {
		if a == "-h" || a == "--help" {
			return metricsShow([]string{"--help"})
		}
		if strings.HasPrefix(a, "-") {
			return fmt.Errorf("unknown flag for metrics show: %s", a)
		}
		if runID == "" {
			runID = a
		} else {
			return fmt.Errorf("metrics show accepts at most one run-id")
		}
	}

	ctx, err := loadMetricsCtx()
	if err != nil {
		return err
	}
	files, err := metrics.ListRunFiles(ctx.RunsDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no runs found")
	}
	var rf *metrics.RunFile
	if runID == "" {
		rf = metrics.NewestRun(files)
	} else {
		rf = metrics.FindRunByID(files, runID)
		if rf == nil {
			return fmt.Errorf("run not found: %s", runID)
		}
	}
	evs, err := metrics.ReadEvents(rf.Path)
	if err != nil {
		return err
	}
	fmt.Printf("run_id: %s\n", rf.ID)
	fmt.Printf("file: %s\n", rf.Name)
	for _, ev := range evs {
		typ, _ := ev["type"].(string)
		switch typ {
		case "run_start":
			fmt.Printf("event: run_start run_id=%v\n", ev["run_id"])
			if mode, ok := ev["mode"].(map[string]any); ok {
				fmt.Printf("  mode.default_suite=%v\n", mode["default_suite"])
			}
		case "leaf_start":
			fmt.Printf("event: leaf_start path=%v labels=%v\n", ev["path"], ev["labels"])
		case "leaf_end":
			fmt.Printf("event: leaf_end path=%v elapsed_ns=%v result=%v\n", ev["path"], ev["elapsed_ns"], ev["result"])
		case "run_end":
			fmt.Printf("event: run_end passed=%v total=%v wall_ns=%v exit_ok=%v\n",
				ev["passed"], ev["total"], ev["wall_ns"], ev["exit_ok"])
		default:
			b, _ := json.Marshal(ev)
			fmt.Printf("event: %s\n", string(b))
		}
	}
	return nil
}

func metricsPrune(args []string) error {
	if helpArgs(args) {
		fmt.Printf("Usage: doctest metrics prune\n\nKeep the newest %d run files; delete older ones.\n", metrics.DefaultRunRetention)
		return nil
	}
	ctx, err := loadMetricsCtx()
	if err != nil {
		return err
	}
	removed, err := metrics.PruneRuns(ctx.RunsDir, metrics.DefaultRunRetention)
	if err != nil {
		return err
	}
	fmt.Printf("prune: removed %d file(s); retention=%d\n", removed, metrics.DefaultRunRetention)
	return nil
}

func helpArgs(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" {
			return true
		}
	}
	return len(args) == 1 && (args[0] == "-h" || args[0] == "--help")
}
