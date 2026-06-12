## Preconditions
- The `doctest`, `fake-codex`, and `yield-pending-questions` binaries are available.
- The orchestrator tests verify the full TDD workflow using these binaries.

## Steps
1. Lookup `fake-codex` from PATH; skip if not installed.
2. Build `doctest` into a temporary executable.
3. Copy doctest as `yield-pending-questions`.
4. Configure env vars so agent implement uses fake-codex.

```go
import (
    "errors"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
    "time"
)

func Setup(t *testing.T, req *Request) error {
    tmp := t.TempDir()

    fakeCodex, err := exec.LookPath("fake-codex")
    if err != nil {
        t.Skip("fake-codex not in PATH; install via: go install github.com/xhd2015/agent-pro/cmd/fake-codex@latest")
        return nil
    }

    doctestBin := filepath.Join(tmp, "doctest")
    buildDT := exec.Command("go", "build", "-o", doctestBin, "./cmd/doctest")
    buildDT.Dir = filepath.Join(DOCTEST_ROOT, "..")
    if out, err := buildDT.CombinedOutput(); err != nil {
        t.Fatalf("build doctest: %v\n%s", err, string(out))
    }

    yieldPQ := filepath.Join(tmp, "yield-pending-questions")
    if out, err := exec.Command("cp", doctestBin, yieldPQ).CombinedOutput(); err != nil {
        t.Fatalf("copy yield-pending-questions: %v\n%s", err, string(out))
    }

    sessionHome := filepath.Join(tmp, "sessions")
    req.Env = append(req.Env,
        "DOCTEST_BIN="+doctestBin,
        "AGENT_RUNNER_FAKE_CODEX_PATH="+fakeCodex,
        "YIELD_PQ_BIN="+yieldPQ,
        "DOCTEST_DEBUG_SESSION_HOME="+sessionHome,
    )
    req.Bin = doctestBin
    os.Setenv("YIELD_PQ_BIN", yieldPQ)
    os.Setenv("DOCTEST_DEBUG_SESSION_HOME", sessionHome)
    req.Timeout = 60 * time.Second
    return nil
}

func writeMockConfig(t *testing.T, req *Request, body string) string {
    t.Helper()
    path := filepath.Join(t.TempDir(), "mock.json")
    if err := os.WriteFile(path, []byte(body), 0644); err != nil {
        t.Fatalf("write mock config: %v", err)
    }
    req.Env = append(req.Env, "FAKE_CODEX_MOCK_CONFIG="+path)
    return path
}

func createDoctestTree(t *testing.T, dir string, stub bool) {
    t.Helper()
    if err := os.MkdirAll(filepath.Join(dir, "basic"), 0755); err != nil {
        t.Fatal(err)
    }

    runBody := "return &Response{Greeting: \"Hello, \" + req.Name + \"!\"}, nil"
    if stub {
        runBody = "return nil, fmt.Errorf(\"error not implemented\")"
    }

    bt := "`"
    goOpen := bt + bt + bt + "go"
    goClose := bt + bt + bt

    rootSetup := "## Preconditions\n- This is a test tree for the greet feature.\n\n## Steps\n1. Set the default name value.\n2. Run invokes the Greet function and returns the greeting.\n\n"
    rootSetup += goOpen + "\n"
    rootSetup += "import (\n"
    rootSetup += "    \"fmt\"\n"
    rootSetup += "    \"testing\"\n"
    rootSetup += ")\n"
    rootSetup += "\n"
    rootSetup += "type Request struct {\n"
    rootSetup += "    Name string\n"
    rootSetup += "}\n"
    rootSetup += "\n"
    rootSetup += "type Response struct {\n"
    rootSetup += "    Greeting string\n"
    rootSetup += "}\n"
    rootSetup += "\n"
    rootSetup += "func Setup(t *testing.T, req *Request) error {\n"
    rootSetup += "    req.Name = \"world\"\n"
    rootSetup += "    return nil\n"
    rootSetup += "}\n"
    rootSetup += "\n"
    rootSetup += "func Run(t *testing.T, req *Request) (*Response, error) {\n"
    rootSetup += "    " + runBody + "\n"
    rootSetup += "}\n"
    rootSetup += goClose + "\n"

    leafSetup := "## Preconditions\n- The request name is \"world\".\n\n## Steps\n1. Set req.Name to \"world\".\n\n"
    leafSetup += goOpen + "\n"
    leafSetup += "func Setup(t *testing.T, req *Request) error {\n"
    leafSetup += "    req.Name = \"world\"\n"
    leafSetup += "    return nil\n"
    leafSetup += "}\n"
    leafSetup += goClose + "\n"

    leafAssert := "## Expected\n- The greeting is \"Hello, world!\".\n\n"
    leafAssert += goOpen + "\n"
    leafAssert += "func Assert(t *testing.T, req *Request, resp *Response, err error) {\n"
    leafAssert += "    if err != nil {\n"
    leafAssert += "        t.Fatal(err)\n"
    leafAssert += "    }\n"
    leafAssert += "    want := \"Hello, world!\"\n"
    leafAssert += "    if resp.Greeting != want {\n"
    leafAssert += "        t.Fatalf(\"expected %q, got %q\", want, resp.Greeting)\n"
    leafAssert += "    }\n"
    leafAssert += "}\n"
    leafAssert += goClose + "\n"

    _write := func(p, c string) {
        if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
            t.Fatal(err)
        }
        if err := os.WriteFile(p, []byte(c), 0644); err != nil {
            t.Fatalf("write %s: %v", p, err)
        }
    }
    _write(filepath.Join(dir, "DOCTEST.md"), "doctest test ./")
    _write(filepath.Join(dir, "SETUP.md"), rootSetup)
    _write(filepath.Join(dir, "basic", "SETUP.md"), leafSetup)
    _write(filepath.Join(dir, "basic", "ASSERT.md"), leafAssert)
}

func writeFile(t *testing.T, path string, content string) {
    t.Helper()
    if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
        t.Fatal(err)
    }
    if err := os.WriteFile(path, []byte(content), 0644); err != nil {
        t.Fatalf("write %s: %v", path, err)
    }
}

func runCmd(t *testing.T, dir string, env []string, name string, args ...string) (string, string, int) {
    t.Helper()
    cmd := exec.Command(name, args...)
    cmd.Dir = dir
    cmd.Env = append(os.Environ(), env...)
    var stdout, stderr strings.Builder
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()
    code := 0
    if err != nil {
        var exitErr *exec.ExitError
        if errors.As(err, &exitErr) {
            code = exitErr.ExitCode()
        } else {
            t.Fatalf("run %s: %v", name, err)
        }
    }
    return stdout.String(), stderr.String(), code
}

func assertContains(t *testing.T, got string, want string) {
    t.Helper()
    if !strings.Contains(got, want) {
        t.Fatalf("missing %q in:\n%s", want, got)
    }
}

func assertExitCode(t *testing.T, code int, want int) {
    t.Helper()
    if code != want {
        t.Fatalf("exit code = %d, want %d", code, want)
    }
}
```
