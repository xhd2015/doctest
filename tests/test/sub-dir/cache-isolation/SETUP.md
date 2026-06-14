## Preconditions
- The root Setup has built the doctest binary and set req.Bin.
- Tests need extended timeout for Go compilation (two doctest invocations).

## Steps
1. Define shared configuration and state for multi-run test leaves.
2. Provide a multi-run `Run` function that creates a test tree with groups and runs doctest twice.
3. Each leaf sets the scenario via `cisoCfg.Scenario`.

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

type CacheIsoCfg struct {
	Scenario string
}

var cisoCfg CacheIsoCfg

type CacheIsoState struct {
	FirstRun  *Response
	SecondRun *Response
}

var isoState CacheIsoState

var bt = "`" + "`" + "`"

func doctestGoBlock(code string) string {
	return bt + "go\n" + code + bt + "\n"
}

func createMultiGroupTree(t *testing.T, dir string) {
	t.Helper()

	rootSetupCode := "import \"testing\"\n\ntype Request struct{}\ntype Response struct{}\n\nfunc Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }\n"
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(doctestGoBlock(rootSetupCode)), 0644); err != nil {
		t.Fatalf("write root SETUP.md: %v", err)
	}

	groups := []struct {
		name  string
		leaves []string
	}{
		{"group-a", []string{"leaf-1", "leaf-2"}},
		{"group-b", []string{"leaf-3"}},
	}

	for _, g := range groups {
		gDir := filepath.Join(dir, g.name)
		if err := os.MkdirAll(gDir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", g.name, err)
		}
		groupSetupCode := "import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n"
		if err := os.WriteFile(filepath.Join(gDir, "SETUP.md"), []byte(doctestGoBlock(groupSetupCode)), 0644); err != nil {
			t.Fatalf("write %s SETUP.md: %v", g.name, err)
		}

		for _, leaf := range g.leaves {
			lDir := filepath.Join(gDir, leaf)
			if err := os.MkdirAll(lDir, 0755); err != nil {
				t.Fatalf("mkdir %s/%s: %v", g.name, leaf, err)
			}
			leafSetupCode := "import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n"
			if err := os.WriteFile(filepath.Join(lDir, "SETUP.md"), []byte(doctestGoBlock(leafSetupCode)), 0644); err != nil {
				t.Fatalf("write %s/%s SETUP.md: %v", g.name, leaf, err)
			}
			leafAssertCode := "import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}\n"
			if err := os.WriteFile(filepath.Join(lDir, "ASSERT.md"), []byte(doctestGoBlock(leafAssertCode)), 0644); err != nil {
				t.Fatalf("write %s/%s ASSERT.md: %v", g.name, leaf, err)
			}
		}
	}
}

func doRun(t *testing.T, bin string, args []string, timeout time.Duration) *Response {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Env = os.Environ()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	resp := &Response{}
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			resp.ExitCode = exitErr.ExitCode()
		} else {
			resp.Err = err
		}
	}
	resp.Stdout = stdoutBuf.String()
	resp.Stderr = stderrBuf.String()
	return resp
}

func Run(t *testing.T, req *Request) (*Response, error) {
	if req.Bin == "" {
		t.Fatalf("req.Bin is not set")
	}

	treeRoot := t.TempDir()
	createMultiGroupTree(t, treeRoot)

	var firstArgs, secondArgs []string
	switch cisoCfg.Scenario {
	case "subdir_after_full":
		firstArgs = []string{"test", treeRoot}
		secondArgs = []string{"test", filepath.Join(treeRoot, "group-a")}
	case "full_after_subdir":
		firstArgs = []string{"test", filepath.Join(treeRoot, "group-a")}
		secondArgs = []string{"test", treeRoot}
	case "subdir_after_subdir":
		firstArgs = []string{"test", filepath.Join(treeRoot, "group-a")}
		secondArgs = []string{"test", filepath.Join(treeRoot, "group-b")}
	default:
		t.Fatalf("unknown scenario: %s", cisoCfg.Scenario)
	}

	runTimeout := 120 * time.Second
	isoState.FirstRun = doRun(t, req.Bin, firstArgs, runTimeout)
	isoState.SecondRun = doRun(t, req.Bin, secondArgs, runTimeout)

	return isoState.SecondRun, nil
}

func Setup(t *testing.T, req *Request) error {
	_ = fmt.Sprintf
	isoState.FirstRun = nil
	isoState.SecondRun = nil
	req.Timeout = 150 * time.Second
	return nil
}
```
