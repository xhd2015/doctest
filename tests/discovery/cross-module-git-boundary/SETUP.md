# Scenario

**Feature**: cross-module git-aware `./...` discovery

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- The doctest binary is built from the module root (`DOCTEST_ROOT/../../..`).
- Each leaf creates a self-contained temp project and runs `doctest test -v ./...`.

## Steps
1. Build temp project layout (go.mod boundaries, optional git repos, doctest trees).
2. Run `doctest test -v ./...` from the project root.
3. Assert discovery, warnings, and exit code per git/module-path scenario.

## Context
- Non-child module paths (`testproj2/sub` vs parent `testproj`) trigger git comparison.
- Child module paths (`testproj/sub`) keep existing walk behavior.
- Warning format: `warning: skipping module <nested> at <dir>: not a child of <ancestor> (<reason>)`

```go
import (
"github.com/xhd2015/doctest/session"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

type gitMode int
type moduleProjectConfig struct {
	parentModulePath	string
	childDir		string
	childModulePath		string
	childTestName		string
	parentTestName		string
	git			gitMode
}
const (
	gitNone	gitMode	= iota
	gitSingleRepo
	gitParentOnly
	gitChildOnly
	gitSeparateRepos
)
var bt = "`" + "`" + "`"
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.UseCLI = true // true e2e product binary
	req.Timeout = 120 * time.Second

	req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", "..", ".."))
	return nil
}
func doctestGoBlock(code string) string {
	return "## Test\n\n" + bt + "go\n" + code + bt + "\n"
}
func doctestBody() string {
	return strings.Join([]string{
		"import \"testing\"",
		"",
		"type Request struct{ Args []string; Env []string; WorkDir string }",
		"type Response struct{ ExitCode int; Stdout string; Stderr string }",
		"",
		"func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }",
	}, "\n")
}
func rootSetupContent() string {
	return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error {\n    t.Logf(\"setup\")\n    return nil\n}")
}
func leafSetupContent() string {
	return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error {\n    t.Logf(\"setup\")\n    return nil\n}")
}
func leafAssertContent() string {
	return doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}")
}
func createTestTree(parent string, name string) error {
	root := filepath.Join(parent, name)
	if err := os.MkdirAll(filepath.Join(root, "simple"), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(root, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(doctestBody())), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(root, "SETUP.md"), []byte(rootSetupContent()), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(root, "simple", "SETUP.md"), []byte(leafSetupContent()), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(root, "simple", "ASSERT.md"), []byte(leafAssertContent()), 0644); err != nil {
		return err
	}
	return nil
}
func initGitRepo(dir string) error {
	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		return fmt.Errorf("git init: %v\n%s", err, string(out))
	}
	return nil
}
func writeGoMod(dir string, modulePath string) error {
	content := "module " + modulePath + "\ngo 1.21\n"
	return os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0644)
}
func createModuleProject(t *testing.T, cfg moduleProjectConfig) string {
	t.Helper()
	tmp := t.TempDir()

	if err := writeGoMod(tmp, cfg.parentModulePath); err != nil {
		t.Fatalf("write parent go.mod: %v", err)
	}

	if cfg.parentTestName != "" {
		if err := createTestTree(tmp, cfg.parentTestName); err != nil {
			t.Fatalf("create parent test tree %s: %v", cfg.parentTestName, err)
		}
	}

	childRoot := filepath.Join(tmp, cfg.childDir)
	if err := os.MkdirAll(childRoot, 0755); err != nil {
		t.Fatalf("mkdir child: %v", err)
	}
	if err := writeGoMod(childRoot, cfg.childModulePath); err != nil {
		t.Fatalf("write child go.mod: %v", err)
	}
	if err := createTestTree(childRoot, cfg.childTestName); err != nil {
		t.Fatalf("create child test tree %s: %v", cfg.childTestName, err)
	}

	switch cfg.git {
	case gitNone:
	case gitSingleRepo:
		if err := initGitRepo(tmp); err != nil {
			t.Fatalf("init single git repo: %v", err)
		}
	case gitParentOnly:
		if err := initGitRepo(tmp); err != nil {
			t.Fatalf("init parent git repo: %v", err)
		}
	case gitChildOnly:
		if err := initGitRepo(childRoot); err != nil {
			t.Fatalf("init child git repo: %v", err)
		}
	case gitSeparateRepos:
		if err := initGitRepo(tmp); err != nil {
			t.Fatalf("init parent git repo: %v", err)
		}
		if err := initGitRepo(childRoot); err != nil {
			t.Fatalf("init child git repo: %v", err)
		}
	}

	return tmp
}
func createLifelogMirrorProject(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()

	if err := writeGoMod(tmp, "github.com/xhd2015/lifelog/tools"); err != nil {
		t.Fatalf("write root go.mod: %v", err)
	}

	cliRoot := filepath.Join(tmp, "lifelog-cli")
	if err := os.MkdirAll(cliRoot, 0755); err != nil {
		t.Fatalf("mkdir lifelog-cli: %v", err)
	}
	if err := writeGoMod(cliRoot, "github.com/xhd2015/lifelog/lifelog-cli"); err != nil {
		t.Fatalf("write lifelog-cli go.mod: %v", err)
	}

	if err := createTestTree(filepath.Join(cliRoot, "tests"), "skill-cli"); err != nil {
		t.Fatalf("create skill-cli test tree: %v", err)
	}

	if err := initGitRepo(tmp); err != nil {
		t.Fatalf("init lifelog mirror git repo: %v", err)
	}

	return tmp
}
```
