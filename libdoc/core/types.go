package core

import (
	"io"
	"time"
)

type ImportSpec struct {
	Name string
	Path string
}

type GoBlock struct {
	SourcePath string
	Code       string

	Imports   []ImportSpec
	TypeDecls []string
	VarDecls  []string
	Consts    []string
	Helpers   []FuncSnippet
	// Methods are package-level methods (have receivers). Emitted outside the test
	// function so they can implement interfaces and reference local-named types.
	Methods []FuncSnippet
	Setup   *FuncSnippet
	Run     *FuncSnippet
	Assert  *FuncSnippet

	Types map[string]bool
}

type FuncSnippet struct {
	Name    string
	// Recv is the receiver field list string (e.g. "f *fakeRunner"). Empty for
	// plain functions. When set, the snippet is a method and must be emitted at
	// package level (not as a func literal).
	Recv    string
	Params  string
	Results string
	// ResultTypes holds return types only (no names), for valid func-literal signatures.
	ResultTypes string
	// ClosureResults holds the result signature rendered for a func literal:
	// names preserved (so bodies that assign to named returns compile) and
	// parenthesized whenever there is more than one result or any named result,
	// because bare multi-name results like "port, alt int" are invalid syntax
	// outside parentheses.
	ClosureResults string
	Body           string
}

type SetupDocument struct {
	Path    string
	GoBlock *GoBlock
}

type AssertDocument struct {
	Path        string
	GoBlock     GoBlock
	Frontmatter AssertFrontmatter
}

type TreeCase struct {
	Name        string
	Path        string
	SetupFiles  []SetupDocument
	AssertFile  AssertDocument
	Labels      []string
	Explanation string
}

// SkippedCase records a labeled leaf omitted from a discovery-mode test run.
type SkippedCase struct {
	Name        string
	Path        string
	Labels      []string
	Explanation string
	DisplayPath string
	Reason      string
}

type ColorMode int

const (
	ColorAuto ColorMode = iota
	ColorAlways
	ColorNever
)

type Options struct {
	GenDir                string
	Verbose               bool
	Stderr                io.Writer
	Stdout                io.Writer
	RemoveTemp            bool
	Count                 int
	Timeout               time.Duration
	SubDir                string
	Color                 ColorMode
	SuppressResultSummary bool
	ChangedOnly           bool
	ExplicitLeaf          bool
	LabelExprs            []string
	// LabelAll, when true, runs labeled leaves under discovery (tree root or
	// ./...). Mutually exclusive with LabelExprs at the CLI. Explicit leaf
	// paths already run labeled leaves without this flag.
	LabelAll bool
	// MetricsOn enables suite metrics JSONL recording (CLI: --metrics-on).
	// Metrics are off by default (opt-in).
	MetricsOn bool
	// MetricsRoot is the cache root for metrics files
	// ($MetricsRoot/doctest/metrics/<project_id>/runs/*.jsonl). Empty means
	// use os.UserCacheDir() (or DOCTEST_METRICS_ROOT when set by the CLI).
	MetricsRoot string
	// MetricsNestSink is the nest JSONL path for this invocation (optional).
	// When set, nested RunTest writes phases here and go test children inherit
	// DOCTEST_METRICS_NEST_SINK via cmd.Env (no mid-run process Setenv).
	MetricsNestSink string
	// MetricsParentLeaf attributes nested phases to this outer leaf path.
	// Prefer explicit Options over process env (parallel-safe).
	MetricsParentLeaf string

	// GenerateOnly stops after writing generated packages (no go test).
	// Used to prepare trees for a multi-root workspace suite run.
	GenerateOnly bool

	// ColdCache enables doctest test --cold-cache: wipe mapping gen root on
	// startup, force -count=1 when unset, and isolate GOCACHE for the run.
	ColdCache bool
	// GoCache is an isolated GOCACHE directory for go test (set by --cold-cache).
	// Empty means inherit the process GOCACHE.
	GoCache string

	// ForceWithFlagA is CLI -a: disable programmatic leaf-cache skip and forward -a
	// to go test (force rebuild of packages that are already up-to-date).
	ForceWithFlagA bool
	// NoLeafCache disables programmatic leaf-cache skip (CLI: --no-leaf-cache).
	// Pass recording (PutPass on success) still occurs.
	NoLeafCache bool

	// Go test profiling / cover flags (forwarded to go test).
	// Path fields are abs-resolved at CLI parse time when relative.
	// Rate fields use *int so zero is distinguishable from unset.
	CPUProfile           string
	MemProfile           string
	MemProfileRate       *int
	BlockProfile         string
	BlockProfileRate     *int
	MutexProfile         string
	MutexProfileFraction *int
	Trace                string
	OutputDir            string
	CoverProfile         string
	Cover                bool
}

type ValidationError struct {
	Path string
	Msg  string
}
