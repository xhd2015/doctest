# Scenario

**Feature**: the `doctest`, `fake-codex`, and `yield-pending-questions` binaries are available

```
# full TDD cycle: design -> RED -> seal -> implement -> GREEN
orchestrator -> design agent -> writes tests -> RED (all fail)

# seal tests, hand off to implementer
orchestrator -> git add tests/ -> implement agent -> writes code -> GREEN (all pass)

# question/answer loop
user <--questions-- implement agent <--yields-- orchestrator -> resume
```

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
"github.com/xhd2015/doctest/session"
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

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Shared helpers only. Do not force UseCLI — short-path leaves stay L2
	// (in-process). True e2e leaves set UseCLI=true and skip without fake-codex.
	req.Timeout = 60 * time.Second
	tmp := t.TempDir()

	doctestBin := testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, ".."))
	req.Bin = doctestBin

	yieldPQ := filepath.Join(tmp, "yield-pending-questions")
	if out, err := exec.Command("cp", doctestBin, yieldPQ).CombinedOutput(); err != nil {
		t.Fatalf("copy yield-pending-questions: %v\n%s", err, string(out))
	}
	sessionHome := filepath.Join(tmp, "sessions")
	req.SessionHome = sessionHome
	req.YieldPQBin = yieldPQ

	// Stash parent-side paths; e2e leaves also push them onto child Env.
	req.Env = append(req.Env,
		"DOCTEST_BIN="+doctestBin,
		"YIELD_PQ_BIN="+yieldPQ,
		"DOCTEST_DEBUG_SESSION_HOME="+sessionHome,
	)

	if fakeCodex, err := exec.LookPath("fake-codex"); err == nil {
		req.Env = append(req.Env, "AGENT_RUNNER_FAKE_CODEX_PATH="+fakeCodex)
	}
	// Child-only env via req.Env — never parent os.Setenv.
	return nil
}

// requireFakeCodex skips when fake-codex is not on PATH (true agent e2e leaves).
func requireFakeCodex(t *testing.T, req *Request) {
	t.Helper()
	for _, e := range req.Env {
		if strings.HasPrefix(e, "AGENT_RUNNER_FAKE_CODEX_PATH=") {
			return
		}
	}
	t.Skip("fake-codex not in PATH; install via: go install github.com/xhd2015/agent-pro/cmd/fake-codex@latest")
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

    goBody := "import (\n" +
        "    \"fmt\"\n" +
        "    \"testing\"\n" +
        ")\n\n" +
        "type Request struct {\n" +
        "    Name string\n" +
        "}\n\n" +
        "type Response struct {\n" +
        "    Greeting string\n" +
        "}\n\n" +
        "func Run(t *testing.T, req *Request) (*Response, error) {\n" +
        "    " + runBody + "\n" +
        "}"

    testtree.WriteFile(t, dir, "DOCTEST.md", testtree.MinimalDOCTEST(goBody))
    bt := "\x60\x60\x60"
    testtree.WriteFile(t, dir, "SETUP.md", "## Preconditions\n- This is a test tree for the greet feature.\n\n## Steps\n1. Set the default name value.\n\n"+bt+"go\nimport \"testing\"\n\nfunc Setup(t *testing.T, req *Request) error {\n    req.Name = \"world\"\n    return nil\n}\n"+bt+"\n")
    testtree.WriteFile(t, dir, "basic/SETUP.md", "## Preconditions\n- The request name is \"world\".\n\n## Steps\n1. Set req.Name to \"world\".\n\n"+bt+"go\nfunc Setup(t *testing.T, req *Request) error {\n    req.Name = \"world\"\n    return nil\n}\n"+bt+"\n")
    testtree.WriteFile(t, dir, "basic/ASSERT.md", "## Expected\n- The greeting is \"Hello, world!\".\n\n"+bt+"go\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {\n    if err != nil {\n        t.Fatal(err)\n    }\n    want := \"Hello, world!\"\n    if resp.Greeting != want {\n        t.Fatalf(\"expected %q, got %q\", want, resp.Greeting)\n    }\n}\n"+bt+"\n")
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
