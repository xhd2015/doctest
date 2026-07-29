# Scenario

**Feature**: go-test plan at run site (gotestmap / suite / hub)

Backfill: verbose `cd … && go test …` lines after `doctest test`, plus pure
TranslatePath checks where path-shaped rules are the contract.

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
