# Scenario

**Feature**: agent commands must be able to use fake Codex instead of a real LLM

```
# agent reads requirement, invokes Fake Codex, writes output
doctest agent <cmd> --requirement req.md -> Fake Codex -> generated code

# session state tracked in event files
doctest <- Fake Codex (session id, events, progress)
```

## Preconditions
- Agent commands must be able to use fake Codex instead of a real LLM.

## Steps
1. Lookup `fake-codex` from PATH; skip if not installed.
2. Write a deterministic fake Codex script.
3. Configure the doctest command to use fake Codex.

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Shared helpers only. Short-path leaves stay L2 (UseCLI=false).
	// True agent e2e leaves set UseCLI=true and call requireFakeCodex.
	req.Timeout = 60 * time.Second
	req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, ".."))
	tmp := t.TempDir()
	if fakeCodex, err := exec.LookPath("fake-codex"); err == nil {
		script := filepath.Join(tmp, "fake-codex-script.json")
		content := `{"events":[{"type":"item.completed","item":{"id":"m1","type":"message","text":"fake doctest agent completed","status":"completed"}}]}`
		if err := os.WriteFile(script, []byte(content), 0644); err != nil {
			t.Fatalf("write fake script: %v", err)
		}
		req.Env = append(req.Env,
			"AGENT_RUNNER_FAKE_CODEX_PATH="+fakeCodex,
			"FAKE_CODEX_SCRIPT="+script,
		)
	}
	return nil
}

// requireFakeCodex skips when fake-codex is not on PATH (true agent e2e leaves).
func requireFakeCodex(t *testing.T, req *Request) {
	t.Helper()
	for _, e := range req.Env {
		if len(e) > len("AGENT_RUNNER_FAKE_CODEX_PATH=") && e[:len("AGENT_RUNNER_FAKE_CODEX_PATH=")] == "AGENT_RUNNER_FAKE_CODEX_PATH=" {
			return
		}
	}
	t.Skip("fake-codex not in PATH; install via: go install github.com/xhd2015/agent-pro/cmd/fake-codex@latest")
}
```
