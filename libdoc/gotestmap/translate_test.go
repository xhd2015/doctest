package gotestmap

import (
	"reflect"
	"testing"
)

func TestPlan_workspaceAndHub(t *testing.T) {
	ws, err := Plan(PlanInput{
		Mode:         ModeWorkspaceSuite,
		RunDir:       "/tmp/gen",
		SuitePattern: "./__workspace/suite",
	})
	if err != nil {
		t.Fatal(err)
	}
	wantWS := []Cmd{{Dir: "/tmp/gen", Pattern: "./__workspace/suite"}}
	if !reflect.DeepEqual(ws, wantWS) {
		t.Fatalf("workspace: got %v want %v", cmds(ws), cmds(wantWS))
	}

	hub, err := Plan(PlanInput{
		Mode:         ModeHubSuite,
		RunDir:       "/tmp/gen/__hub",
		SuitePattern: "./suite",
	})
	if err != nil {
		t.Fatal(err)
	}
	wantHub := []Cmd{{Dir: "/tmp/gen/__hub", Pattern: "./suite"}}
	if !reflect.DeepEqual(hub, wantHub) {
		t.Fatalf("hub: got %v want %v", cmds(hub), cmds(wantHub))
	}
}

func TestPlan_pathShaped(t *testing.T) {
	got, err := Plan(PlanInput{
		Mode:    ModePathShaped,
		UserArg: "./tree/mid/...",
		Layout:  Layout{ModuleRoots: []string{".", "tree/mid/nestedmod"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []Cmd{
		{Dir: ".", Pattern: "./tree/mid/..."},
		{Dir: "tree/mid/nestedmod", Pattern: "./..."},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", cmds(got), cmds(want))
	}
}

func TestNeedsPathShaped(t *testing.T) {
	layout := Layout{ModuleRoots: []string{".", "tree/mid/nestedmod"}}
	cases := []struct {
		arg  string
		lay  Layout
		want bool
	}{
		{"./tests/foo", Layout{ModuleRoots: []string{"."}}, false},
		{"./tests/foo/...", Layout{ModuleRoots: []string{"."}}, false},
		{"./...", Layout{ModuleRoots: []string{"."}}, false},
		{"./tree/mid/...", layout, true},  // nestedmod under mid
		{"./tree/mid", layout, false},     // no dots → suite shape
		{"./...", layout, true},           // nested modules
		{"./tree/...", layout, true},      // nestedmod under tree
		{"./tree/mid/two", Layout{ModuleRoots: []string{"."}}, false},
		{"./tree/mid/nestedmod/suite", layout, true}, // inside nested module
	}
	for _, tc := range cases {
		if g := NeedsPathShaped(tc.arg, tc.lay); g != tc.want {
			t.Errorf("NeedsPathShaped(%q)=%v want %v", tc.arg, g, tc.want)
		}
	}
}

// MECE path-shaped cases (fixture-relative).
func TestTranslatePath_MECE(t *testing.T) {
	layoutMid := Layout{ModuleRoots: []string{".", "tree/mid/nestedmod"}}

	tests := []struct {
		name string
		arg  string
		lay  Layout
		want []Cmd
	}{
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
		{
			name: "mid_dotdotdot_outer_only_no_nested_mod",
			arg:  "./tree/mid/...",
			lay:  Layout{ModuleRoots: []string{"."}},
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
		{
			name: "whole_tree_dotdotdot",
			arg:  "./tree/...",
			lay:  layoutMid,
			want: []Cmd{
				{Dir: ".", Pattern: "./tree/..."},
				{Dir: "tree/mid/nestedmod", Pattern: "./..."},
			},
		},
		{
			name: "dot_dotdotdot",
			arg:  "./...",
			lay:  layoutMid,
			want: []Cmd{
				{Dir: ".", Pattern: "./..."},
				{Dir: "tree/mid/nestedmod", Pattern: "./..."},
			},
		},
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
		{
			name: "mid_excludes_nested_mod_outside_prefix",
			arg:  "./tree/mid/...",
			lay:  Layout{ModuleRoots: []string{".", "tree/other/nestedmod"}},
			want: []Cmd{
				{Dir: ".", Pattern: "./tree/mid/..."},
			},
		},
		{
			name: "mid_dotdotdot_no_dot_slash",
			arg:  "tree/mid/...",
			lay:  Layout{ModuleRoots: []string{"."}},
			want: []Cmd{{Dir: ".", Pattern: "./tree/mid/..."}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TranslatePath(tt.arg, tt.lay)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v want %v", cmds(got), cmds(tt.want))
			}
			// Plan ModePathShaped must match.
			viaPlan, err := Plan(PlanInput{Mode: ModePathShaped, UserArg: tt.arg, Layout: tt.lay})
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(viaPlan, tt.want) {
				t.Fatalf("Plan: got %v want %v", cmds(viaPlan), cmds(tt.want))
			}
		})
	}
}

func TestTranslate_bareEllipsisError(t *testing.T) {
	_, err := TranslatePath("...", Layout{ModuleRoots: []string{"."}})
	if err == nil {
		t.Fatal("want error")
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
	if g, w := (Cmd{Dir: "/tmp/gen", Pattern: "./__workspace/suite"}).String(),
		"(cd /tmp/gen && go test ./__workspace/suite)"; g != w {
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
