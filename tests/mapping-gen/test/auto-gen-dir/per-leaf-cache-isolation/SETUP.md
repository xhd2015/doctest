## Preconditions
- A project with 2 leaves exists.
- The auto gen dir is used for both runs.
- The Run function handles two executions: first run (warmup), modify one leaf, second run.

## Steps
1. Create a project with 2 leaves: `leaf_a` and `leaf_b`.
2. Run once to warm up the cache.
3. Modify `leaf_a`'s ASSERT.md content.
4. Run again.
5. Verify `leaf_b` is cached (no recompilation), `leaf_a` is not.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	genDirPath string
	firstResp  *Response
	secondResp *Response
)

func Setup(t *testing.T, req *Request) error {
	genDirPath = ""
	firstResp = nil
	secondResp = nil
	_ = fmt.Sprintf
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	if req.Bin == "" {
		t.Fatalf("req.Bin is not set")
	}

	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	testDir := filepath.Join(proj, "tests")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("mkdir test dir: %v", err)
	}
	runCode := "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }"
	if err := createDoctestRoot(testDir, runCode); err != nil {
		t.Fatalf("create doctest root: %v", err)
	}
	if err := createDoctestLeaf(filepath.Join(testDir, "leaf_a")); err != nil {
		t.Fatalf("create leaf_a: %v", err)
	}
	if err := createDoctestLeaf(filepath.Join(testDir, "leaf_b")); err != nil {
		t.Fatalf("create leaf_b: %v", err)
	}

	baseArgs := []string{"test", testDir}

	resp1 := doRun(t, req.Bin, baseArgs)
	genDirPath = parseGenDir(resp1.Stderr)
	if genDirPath == "" {
		t.Fatalf("could not parse gen dir from first run stderr:\n%s", resp1.Stderr)
	}

	leafAAssertPath := filepath.Join(testDir, "leaf_a", "ASSERT.md")
	data, err := os.ReadFile(leafAAssertPath)
	if err != nil {
		t.Fatalf("read leaf_a ASSERT.md: %v", err)
	}
	modified := strings.Replace(string(data), "func Assert", "// modified\nfunc Assert", 1)
	if err := os.WriteFile(leafAAssertPath, []byte(modified), 0644); err != nil {
		t.Fatalf("modify leaf_a ASSERT.md: %v", err)
	}

	resp2 := doRun(t, req.Bin, baseArgs)

	firstResp = resp1
	secondResp = resp2
	return resp2, nil
}

func doRun(t *testing.T, bin string, args []string) *Response {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Env = os.Environ()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	resp := &Response{}
	err := cmd.Run()
	resp.Stdout = stdoutBuf.String()
	resp.Stderr = stderrBuf.String()
	if err == nil {
		return resp
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		resp.ExitCode = exitErr.ExitCode()
		return resp
	}
	resp.Err = err
	return resp
}
```
