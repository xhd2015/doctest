package gotestmap

import (
	"reflect"
	"testing"
)

// MECE coverage of ideal go-test translation for doctest path args.
// Translate() is still broken → these tests must FAIL (RED) until IdealTranslate is wired.

func TestTranslate_idealMECE(t *testing.T) {
	// Shared layout: outer + sibling shape + same-mod path + nested module under mid.
	// Module roots only (go.mod locations); packages are implied by patterns.
	layoutMid := Layout{ModuleRoots: []string{".", "tree/mid/nestedmod"}}

	tests := []struct {
		name string
		arg  string
		lay  Layout
		want []Cmd
	}{
		// --- recursion axis: no /... ---
		{
			name: "leaf_or_pkg_no_dots",
			arg:  "./tree/mid/two",
			lay:  layoutMid,
			want: []Cmd{{Dir: ".", Pattern: "./tree/mid/two"}},
		},
		{
			name: "mid_branch_no_dots",
			arg:  "./tree/mid",
			lay:  layoutMid,
			want: []Cmd{{Dir: ".", Pattern: "./tree/mid"}},
		},

		// --- recursion axis: /... mid prefix (core CTF shape) ---
		{
			name: "mid_dotdotdot_outer_only_no_nested_mod",
			arg:  "./tree/mid/...",
			lay:  Layout{ModuleRoots: []string{"."}},
			// must NOT be ./tree/... (sibling would be included)
			want: []Cmd{{Dir: ".", Pattern: "./tree/mid/..."}},
		},
		{
			name: "mid_dotdotdot_with_nested_gomod_BC",
			arg:  "./tree/mid/...",
			lay:  layoutMid,
			want: []Cmd{
				{Dir: ".", Pattern: "./tree/mid/..."},
				{Dir: "tree/mid/nestedmod", Pattern: "./..."},
			},
		},
		{
			name: "mid_dotdotdot_two_nested_gomods",
			arg:  "./tree/mid/...",
			lay:  Layout{ModuleRoots: []string{".", "tree/mid/mod_a", "tree/mid/mod_b"}},
			want: []Cmd{
				{Dir: ".", Pattern: "./tree/mid/..."},
				{Dir: "tree/mid/mod_a", Pattern: "./..."},
				{Dir: "tree/mid/mod_b", Pattern: "./..."},
			},
		},

		// --- nested module under mid must not pull sibling of mid ---
		// (encoded as pattern never ./tree/... )
		{
			name: "mid_dotdotdot_never_whole_tree",
			arg:  "./tree/mid/...",
			lay:  layoutMid,
			want: []Cmd{
				{Dir: ".", Pattern: "./tree/mid/..."},
				{Dir: "tree/mid/nestedmod", Pattern: "./..."},
			},
		},

		// --- full tree root ---
		{
			name: "whole_tree_dotdotdot",
			arg:  "./tree/...",
			lay:  layoutMid,
			// nestedmod is under tree/ → included
			want: []Cmd{
				{Dir: ".", Pattern: "./tree/..."},
				{Dir: "tree/mid/nestedmod", Pattern: "./..."},
			},
		},

		// --- module-wide ---
		{
			name: "dot_dotdotdot",
			arg:  "./...",
			lay:  layoutMid,
			want: []Cmd{
				{Dir: ".", Pattern: "./..."},
				{Dir: "tree/mid/nestedmod", Pattern: "./..."},
			},
		},

		// --- arg already inside nested module ---
		{
			name: "inside_nested_mod_dotdotdot",
			arg:  "./tree/mid/nestedmod/suite/...",
			lay:  layoutMid,
			want: []Cmd{
				{Dir: "tree/mid/nestedmod", Pattern: "./suite/..."},
			},
		},
		{
			name: "inside_nested_mod_no_dots",
			arg:  "./tree/mid/nestedmod/suite",
			lay:  layoutMid,
			want: []Cmd{
				{Dir: "tree/mid/nestedmod", Pattern: "./suite"},
			},
		},

		// --- nested go.mod outside selected mid: not included ---
		{
			name: "mid_excludes_nested_mod_outside_prefix",
			arg:  "./tree/mid/...",
			lay:  Layout{ModuleRoots: []string{".", "tree/other/nestedmod"}},
			want: []Cmd{
				{Dir: ".", Pattern: "./tree/mid/..."},
				// tree/other/nestedmod NOT under mid → absent
			},
		},

		// --- no ./ prefix ---
		{
			name: "mid_dotdotdot_no_dot_slash",
			arg:  "tree/mid/...",
			lay:  Layout{ModuleRoots: []string{"."}},
			want: []Cmd{{Dir: ".", Pattern: "./tree/mid/..."}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Lock ideal in the test name; call product Translate (still broken).
			ideal, err := IdealTranslate(tt.arg, tt.lay)
			if err != nil {
				t.Fatalf("IdealTranslate: %v", err)
			}
			if !reflect.DeepEqual(ideal, tt.want) {
				t.Fatalf("IdealTranslate mismatch (test table wrong):\n got %v\nwant %v", cmds(ideal), cmds(tt.want))
			}

			got, err := Translate(tt.arg, tt.lay)
			if err != nil {
				t.Fatalf("Translate: %v", err)
			}
			if reflect.DeepEqual(got, tt.want) {
				// If someone wires IdealTranslate, this goes green — OK.
				return
			}
			t.Fatalf("Translate RED (want ideal path-shaped go test):\n  arg=%q\n  got  %v\n  want %v",
				tt.arg, cmds(got), cmds(tt.want))
		})
	}
}

func TestTranslate_bareEllipsisError(t *testing.T) {
	_, err := IdealTranslate("...", Layout{ModuleRoots: []string{"."}})
	if err == nil {
		t.Fatal("IdealTranslate(...): want error")
	}
	_, err = Translate("...", Layout{ModuleRoots: []string{"."}})
	if err == nil {
		t.Fatal("Translate(...): want error")
	}
}

func TestCmdString(t *testing.T) {
	if g, w := (Cmd{Dir: ".", Pattern: "./tree/mid/..."}).String(), "go test ./tree/mid/..."; g != w {
		t.Fatalf("got %q want %q", g, w)
	}
	if g, w := (Cmd{Dir: "tree/mid/nestedmod", Pattern: "./..."}).String(),
		"(cd tree/mid/nestedmod && go test ./...)"; g != w {
		t.Fatalf("got %q want %q", g, w)
	}
}

func cmds(c []Cmd) []string {
	out := make([]string, len(c))
	for i := range c {
		out[i] = c[i].String()
	}
	return out
}
