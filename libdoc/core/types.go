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
	Setup     *FuncSnippet
	Run       *FuncSnippet
	Assert    *FuncSnippet

	Types map[string]bool
}

type FuncSnippet struct {
	Name    string
	Params  string
	Results string
	Body    string
}

type SetupDocument struct {
	Path    string
	GoBlock *GoBlock
}

type AssertDocument struct {
	Path    string
	GoBlock GoBlock
}

type TreeCase struct {
	Name       string
	Path       string
	SetupFiles []SetupDocument
	AssertFile AssertDocument
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
	RemoveTemp            bool
	Count                 int
	Timeout               time.Duration
	SubDir                string
	Color                 ColorMode
	SuppressResultSummary bool
}

type ValidationError struct {
	Path string
	Msg  string
}
