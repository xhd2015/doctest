package registry

import "testing"

type Entry struct {
	Path string
	Fn   func(*testing.T)
}

var entries []Entry

func Register(e Entry) { entries = append(entries, e) }
func All() []Entry     { return entries }
