# Scenario

**Feature**: go-test plan at run site (gotestmap / suite / hub)

Backfill (audit compliance): production go-test path is **single-cmd only**
until Phase 2 multi-cmd path-shaped execution exists.

- **ModeWorkspaceSuite** → one plan: `cd <gen> && go test ./__workspace/suite`
- **ModeHubSuite** → one plan: `cd <…/__hub> && go test ./suite`
- **ModePathShaped** → pure TranslatePath / Plan contract only (may return
  multiple Cmds); **not** exercised via multi-cmd `finishWorkspaceGoTestCmds`
  in production CLI until Phase 2.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/session"
)

func gotestPlanOut(resp *Response) string {
	if resp == nil {
		return ""
	}
	return resp.Stdout + "\n" + resp.Stderr
}

// countCdGoTestPlanLines counts production plan lines printed as
// `cd <dir> && go test …` (finishWorkspaceGoTest). Not package FAIL/ok lines.
func countCdGoTestPlanLines(out string) int {
	n := 0
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "&& go test") {
			n++
		}
	}
	return n
}

// assertExactlyOneGoTestPlanFamily requires a single cd…&& go test plan line
// that contains sub (e.g. __workspace/suite or __hub). Locks single-cmd suite/hub.
func assertExactlyOneGoTestPlanFamily(t *testing.T, out, sub string) {
	t.Helper()
	n := countCdGoTestPlanLines(out)
	if n != 1 {
		t.Fatalf("want exactly 1 cd…&& go test plan line (single-cmd suite/hub), got %d:\n%s", n, out)
	}
	assertContainsGoTestLine(t, out, sub)
}

// createSingleModTwoTrees: one go.mod, two DOCTEST roots → single-gen workspace suite.
func createSingleModTwoTrees(t *testing.T) string {
	t.Helper()
	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module gotestplan\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := createTestTree(proj, "alpha"); err != nil {
		t.Fatal(err)
	}
	if err := createTestTree(proj, "beta"); err != nil {
		t.Fatal(err)
	}
	return proj
}

// createChildModuleProject: outer + nested go.mod with child module path → multi-gen hub.
func createChildModuleProject(t *testing.T) string {
	t.Helper()
	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module gotestplan\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := createTestTree(proj, "parent_test"); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(proj, "sub")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "go.mod"), []byte("module gotestplan/sub\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := createTestTree(nested, "child_test"); err != nil {
		t.Fatal(err)
	}
	return proj
}

func assertContainsGoTestLine(t *testing.T, out, sub string) {
	t.Helper()
	if !strings.Contains(out, "go test") {
		t.Fatalf("missing go test line:\n%s", out)
	}
	if !strings.Contains(out, sub) {
		t.Fatalf("go test plan missing %q:\n%s", sub, out)
	}
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Timeout < 2*time.Minute {
		req.Timeout = 2 * time.Minute
	}
	return nil
}
```
