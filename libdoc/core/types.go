package core

import "io"

type GoBlock struct {
	SourcePath string
	Code       string

	Imports   []string
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

type Options struct {
	GenDir     string
	Verbose    bool
	Stderr     io.Writer
	RemoveTemp bool
	Count      int
	SubDir     string
}

type ValidationError struct {
	Path string
	Msg  string
}
