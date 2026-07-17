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
}

type ValidationError struct {
	Path string
	Msg  string
}
