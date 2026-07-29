## Expected

- `Plan(ModeWorkspaceSuite)` → `len(cmds)==1`, pattern `__workspace/suite` (or given SuitePattern).
- `Plan(ModeHubSuite)` → `len(cmds)==1`, pattern `./suite` (default).
- Empty RunDir errors (suite/hub require RunDir).
- ModePathShaped may return multi cmds (contract only; execution is Phase 2) — assert mid+nested len==2 without claiming CLI finish.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/gotestmap"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ExitCode != 0 {
		t.Fatalf("--help exit=%d\n%s\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}

	ws, err := gotestmap.Plan(gotestmap.PlanInput{
		Mode:         gotestmap.ModeWorkspaceSuite,
		RunDir:       "/tmp/gen",
		SuitePattern: "./__workspace/suite",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ws) != 1 {
		t.Fatalf("ModeWorkspaceSuite: want len==1, got %d (%v)", len(ws), ws)
	}
	if ws[0].Dir != "/tmp/gen" || ws[0].Pattern != "./__workspace/suite" {
		t.Fatalf("ModeWorkspaceSuite cmd: %+v", ws[0])
	}

	// Default pattern when SuitePattern omitted.
	wsDef, err := gotestmap.Plan(gotestmap.PlanInput{
		Mode:   gotestmap.ModeWorkspaceSuite,
		RunDir: "/tmp/gen",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(wsDef) != 1 || wsDef[0].Pattern != "./__workspace/suite" {
		t.Fatalf("ModeWorkspaceSuite default pattern: %v", wsDef)
	}

	hub, err := gotestmap.Plan(gotestmap.PlanInput{
		Mode:   gotestmap.ModeHubSuite,
		RunDir: "/tmp/gen/__hub",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hub) != 1 {
		t.Fatalf("ModeHubSuite: want len==1, got %d (%v)", len(hub), hub)
	}
	if hub[0].Pattern != "./suite" {
		t.Fatalf("ModeHubSuite default pattern: %+v", hub[0])
	}

	_, err = gotestmap.Plan(gotestmap.PlanInput{Mode: gotestmap.ModeWorkspaceSuite})
	if err == nil || !strings.Contains(err.Error(), "RunDir") {
		t.Fatalf("workspace without RunDir: err=%v", err)
	}
	_, err = gotestmap.Plan(gotestmap.PlanInput{Mode: gotestmap.ModeHubSuite})
	if err == nil || !strings.Contains(err.Error(), "RunDir") {
		t.Fatalf("hub without RunDir: err=%v", err)
	}

	// Path-shaped multi-cmd is a pure plan contract (Phase 2 execution).
	pathCmds, err := gotestmap.Plan(gotestmap.PlanInput{
		Mode:    gotestmap.ModePathShaped,
		UserArg: "./tree/mid/...",
		Layout:  gotestmap.Layout{ModuleRoots: []string{".", "tree/mid/nestedmod"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(pathCmds) != 2 {
		t.Fatalf("ModePathShaped mid+nested: want len==2 plan cmds (not suite/hub), got %d %v", len(pathCmds), pathCmds)
	}
}
```
