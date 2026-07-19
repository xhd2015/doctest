# Scenario

**Feature**: internal import scan triggers temp compile under parent module

```
# internal import detected in assembled Go
doctest test/build <tree> -> scan imports -> .doctest_run_* under moduleRoot -> go test

# optional gen-dir dump (review copy, not compile root)
--gen-dir DIR -> copy generated files to DIR (no nested go.mod)

# no internal import: legacy nested module testcase
public imports only -> cache/outside gen-dir -> module testcase + replace
```

## Preconditions

- The doctest module root is three levels above this test tree (`DOCTEST_ROOT/../../..`).
- Each leaf creates a temporary Go module with `internal/greet` or public `pkg/greet`.
- `GOWORK=off` is set in subprocess env so workspace mode does not mask module behavior.
- `req.WorkDir` is set to the temp module root for doctest invocations.
- Tests do not use `--in-module` or `--nested-module` (removed from design).

## Steps

1. Build the doctest binary from the module root.
2. Create a temp module and doctest tree per leaf scenario.
3. Execute the doctest binary with leaf-specific args and capture output.

## Context

- Package-level vars `moduleRoot`, `genDir`, and `testDir` are set by shared helpers.
- Feature leaves assert import-scan + temp-compile behavior (RED before implementation).
- `legacy/nested-module-unchanged` verifies unchanged public-import legacy path.

```go
import (
"github.com/xhd2015/doctest/session"
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

const modPath = "example.com/app"
var (
	moduleRoot	string
	genDir		string
	testDir		string
)
var bt = string([]byte{96, 96, 96})
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Timeout = 120 * time.Second

	req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", "..", ".."))
	return nil
}
func runDoctestInterruptedDuringWriteCases(t *testing.T, req *Request) (*Response, error) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	defer cancel()

	bin := req.Bin
	if bin == "" {
		return nil, fmt.Errorf("req.Bin is not set")
	}
	triggerLeaf := req.InterruptTriggerLeaf
	if triggerLeaf <= 0 {
		triggerLeaf = 15
	}
	trigger := fmt.Sprintf("leaf%02d_test.go", triggerLeaf)

	cmd := exec.CommandContext(ctx, bin, req.Args...)
	cmd.Dir = req.WorkDir
	cmd.Env = append(os.Environ(), req.Env...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start doctest: %w", err)
	}

	var stdoutBuf, stderrBuf strings.Builder
	interrupted := false
	stderrDone := make(chan struct{})
	go func() {
		defer close(stderrDone)
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			stderrBuf.WriteString(line)
			stderrBuf.WriteByte('\n')
			if !interrupted && strings.Contains(line, trigger) {
				interrupted = true
				_ = cmd.Process.Signal(os.Interrupt)
			}
		}
	}()
	go func() {
		_, _ = io.Copy(&stdoutBuf, stdoutPipe)
	}()

	waitErr := cmd.Wait()
	<-stderrDone

	resp := &Response{
		Stdout:	stdoutBuf.String(),
		Stderr:	stderrBuf.String(),
		Err:	waitErr,
	}
	if waitErr != nil {
		var exitErr *exec.ExitError
		if errors.As(waitErr, &exitErr) {
			resp.ExitCode = exitErr.ExitCode()
		}
	}
	if !interrupted {
		return resp, fmt.Errorf("never sent SIGINT: trigger %q not seen in stderr:\n%s", trigger, resp.Stderr)
	}
	return resp, nil
}
func doctestGoBlock(code string) string {
	return "## Test\n\n" + bt + "go\n" + code + bt + "\n"
}
func createDoctestRoot(dir string, extraImports string, runCode string) error {
	goBody := "import (\n" +
		"\t\"testing\"\n" +
		extraImports +
		")\n\n" +
		"type Request struct{}\n" +
		"type Response struct{ Message string }\n\n" +
		runCode
	if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(goBody)), 0644); err != nil {
		return err
	}
	return nil
}
func createDoctestLeaf(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	leafSetup := doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }")
	leafAssert := doctestGoBlock(
		"import \"testing\"\n" +
			"func Assert(t *testing.T, req *Request, resp *Response, err error) {\n" +
			"\tif err != nil { t.Fatal(err) }\n" +
			"\tif resp.Message != \"hi\" { t.Fatalf(\"expected hi, got %q\", resp.Message) }\n" +
			"}",
	)
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(leafSetup), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "ASSERT.md"), []byte(leafAssert), 0644); err != nil {
		return err
	}
	return nil
}
func copyDir(dst string, src string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(src, path)
		if relErr != nil {
			return relErr
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		return os.WriteFile(target, data, 0644)
	})
}
func copyInternalModuleFixture(t *testing.T, d *session.Doctest, fixtureName string) {
	t.Helper()

	fixtureSrc := filepath.Join(d.DOCTEST_ROOT, "testdata", fixtureName)
	moduleRoot = t.TempDir()
	if err := copyDir(moduleRoot, fixtureSrc); err != nil {
		t.Fatalf("copy fixture %s: %v", fixtureName, err)
	}
	testDir = filepath.Join(moduleRoot, "tests")
	genDir = filepath.Join(moduleRoot, "_gen")
}
func createInternalModuleProject(t *testing.T, d *session.Doctest) {
	t.Helper()
	copyInternalModuleFixture(t, d, "internal-module")
}
func createInternalModuleProjectWithLeaves(t *testing.T, d *session.Doctest) {
	t.Helper()

	fixtureSrc := filepath.Join(d.DOCTEST_CASE, "testdata")
	moduleRoot = t.TempDir()
	if err := copyDir(moduleRoot, fixtureSrc); err != nil {
		t.Fatalf("copy fixture %s: %v", fixtureSrc, err)
	}
	testDir = filepath.Join(moduleRoot, "tests")
	genDir = filepath.Join(moduleRoot, "_gen")
}
func createPublicModuleProject(t *testing.T) {
	t.Helper()

	moduleRoot = t.TempDir()
	if err := os.WriteFile(
		filepath.Join(moduleRoot, "go.mod"),
		[]byte("module "+modPath+"\n\ngo 1.21\n"),
		0644,
	); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	greetDir := filepath.Join(moduleRoot, "pkg", "greet")
	if err := os.MkdirAll(greetDir, 0755); err != nil {
		t.Fatalf("mkdir pkg/greet: %v", err)
	}
	if err := os.WriteFile(
		filepath.Join(greetDir, "greet.go"),
		[]byte("package greet\n\nfunc Hello() string { return \"hi\" }\n"),
		0644,
	); err != nil {
		t.Fatalf("write greet.go: %v", err)
	}

	testDir = filepath.Join(moduleRoot, "tests")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("mkdir tests: %v", err)
	}
	extraImports := "\t\"" + modPath + "/pkg/greet\"\n"
	runCode := "func Run(t *testing.T, req *Request) (*Response, error) {\n" +
		"\treturn &Response{Message: greet.Hello()}, nil\n" +
		"}"
	if err := createDoctestRoot(testDir, extraImports, runCode); err != nil {
		t.Fatalf("create doctest root: %v", err)
	}
	if err := createDoctestLeaf(filepath.Join(testDir, "leaf")); err != nil {
		t.Fatalf("create doctest leaf: %v", err)
	}

	genDir = filepath.Join(moduleRoot, "_gen")
}
func setupModuleEnv(t *testing.T, req *Request) {
	t.Helper()
	req.WorkDir = moduleRoot
	req.Env = append(req.Env, "GOWORK=off")
}
func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %s to exist: %v", path, err)
	}
}
func assertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("expected file %s to not exist", path)
	}
}
func assertNestedGoMod(t *testing.T, dir string) {
	t.Helper()
	nestedGoMod := filepath.Join(dir, "go.mod")
	assertFileExists(t, nestedGoMod)
	goModData, readErr := os.ReadFile(nestedGoMod)
	if readErr != nil {
		t.Fatalf("read nested go.mod: %v", readErr)
	}
	goMod := string(goModData)
	if !strings.Contains(goMod, "module testcase") {
		t.Fatalf("expected nested module testcase, got:\n%s", goMod)
	}
	if !strings.Contains(goMod, "replace "+modPath+" =>") {
		t.Fatalf("expected replace directive for parent module, got:\n%s", goMod)
	}
}
func generatedLeafTestPath(genRoot string) string {
	return filepath.Join(genRoot, "tests", "leaf", "leaf_test.go")
}
func findDoctestRunDirs(root string) ([]string, error) {
	var found []string
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), ".doctest_run_") {
			found = append(found, filepath.Join(root, entry.Name()))
		}
	}
	return found, nil
}
func assertNoDoctestRunDirs(t *testing.T, root string) {
	t.Helper()
	dirs, err := findDoctestRunDirs(root)
	if err != nil {
		t.Fatalf("scan for .doctest_run_* dirs: %v", err)
	}
	if len(dirs) > 0 {
		t.Fatalf("expected no .doctest_run_* dirs under %s, found: %v", root, dirs)
	}
}
func assertStderrUsesTempCompile(t *testing.T, resp *Response) {
	t.Helper()
	combined := resp.Stdout + resp.Stderr
	if !strings.Contains(combined, ".doctest_run_") {
		t.Fatalf("expected stderr/stdout to reference .doctest_run_ temp compile dir, got:\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
}
func assertDumpHasInternalImport(t *testing.T, dumpRoot string) {
	t.Helper()
	genTest := generatedLeafTestPath(dumpRoot)
	assertFileExists(t, genTest)
	testData, readErr := os.ReadFile(genTest)
	if readErr != nil {
		t.Fatalf("read dump test file: %v", readErr)
	}
	if !strings.Contains(string(testData), modPath+"/internal/greet") {
		t.Fatalf("expected dump to import internal/greet, got:\n%s", string(testData))
	}
}
func assertDumpNoNestedGoMod(t *testing.T, dumpRoot string) {
	t.Helper()
	assertFileNotExists(t, filepath.Join(dumpRoot, "go.mod"))
}
```
