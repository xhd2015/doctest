## Preconditions
- The doctest binary is built from the module root (`DOCTEST_ROOT/..`).
- `fake-codex` is expected in PATH.
- Tests are executed by the doc-style test runner from this tree.
- Agent invocations use `fake-codex` to avoid real LLM calls.

## Steps
1. Build `doctest` into a temp binary.
2. Lookup `fake-codex` from PATH; skip if not installed.
3. Copy doctest as `yield-pending-questions` for ypq tests.
4. Configure env vars so agent implement uses fake-codex.
5. Provide helper functions for writing mock configs and inspecting session metadata.

## Context
- Session directories are created under `~/.doctest/implementer/sessions/`.
- Mock config format follows the `fake-codex` conventions.

```go
import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
    "time"
)

func needsBuildVCSFlag(dir string) bool {
    git, err := exec.LookPath("git")
    if err != nil {
        return true
    }
    cmd := exec.Command(git, "-C", dir, "rev-parse", "--is-inside-work-tree")
    out, err := cmd.Output()
    if err != nil {
        return true
    }
    return strings.TrimSpace(string(out)) != "true"
}

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 60 * time.Second

    tmp := t.TempDir()

    doctestBin := filepath.Join(tmp, "doctest")
    buildDir := filepath.Join(DOCTEST_ROOT, "..")
    buildArgs := []string{"build", "-o", doctestBin}
    if needsBuildVCSFlag(buildDir) {
        buildArgs = append(buildArgs, "-buildvcs=false")
    }
    buildArgs = append(buildArgs, "./cmd/doctest")
    build := exec.Command("go", buildArgs...)
    build.Dir = buildDir
    if out, err := build.CombinedOutput(); err != nil {
        t.Fatalf("build doctest: %v\n%s", err, string(out))
    }

    fakeCodex, err := exec.LookPath("fake-codex")
    if err != nil {
        t.Skip("fake-codex not in PATH; install via: go install github.com/xhd2015/agent-pro/cmd/fake-codex@latest")
        return nil
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

func sessionsDir() string {
    if v := os.Getenv("DOCTEST_DEBUG_SESSION_HOME"); v != "" {
        return v
    }
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".doctest", "implementer", "sessions")
}

func findSessionMeta(t *testing.T, field string, value string) string {
    t.Helper()
    base := sessionsDir()
    today := time.Now()
    for i := 0; i < 7; i++ {
        dateDir := today.AddDate(0, 0, -i).Format("2006/01/02")
        dayPath := filepath.Join(base, dateDir)
        entries, err := os.ReadDir(dayPath)
        if err != nil {
            continue
        }
        for _, entry := range entries {
            if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "sess_") {
                continue
            }
            metaPath := filepath.Join(dayPath, entry.Name(), "meta.json")
            data, err := os.ReadFile(metaPath)
            if err != nil {
                continue
            }
            var m map[string]any
            if err := json.Unmarshal(data, &m); err != nil {
                continue
            }
            if v, ok := m[field]; ok {
                if s, ok := v.(string); ok && s == value {
                    return metaPath
                }
            }
        }
    }
    return ""
}

func readMetaJSON(t *testing.T, path string) map[string]any {
    t.Helper()
    data, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("read meta.json %s: %v", path, err)
    }
    var m map[string]any
    if err := json.Unmarshal(data, &m); err != nil {
        t.Fatalf("parse meta.json %s: %v", path, err)
    }
    return m
}

func countSessionDirs(t *testing.T, field string, value string) int {
    t.Helper()
    base := sessionsDir()
    count := 0
    today := time.Now()
    for i := 0; i < 7; i++ {
        dateDir := today.AddDate(0, 0, -i).Format("2006/01/02")
        dayPath := filepath.Join(base, dateDir)
        entries, err := os.ReadDir(dayPath)
        if err != nil {
            continue
        }
        for _, entry := range entries {
            if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "sess_") {
                continue
            }
            metaPath := filepath.Join(dayPath, entry.Name(), "meta.json")
            data, err := os.ReadFile(metaPath)
            if err != nil {
                continue
            }
            var m map[string]any
            if err := json.Unmarshal(data, &m); err != nil {
                continue
            }
            if v, ok := m[field]; ok {
                if s, ok := v.(string); ok && s == value {
                    count++
                }
            }
        }
    }
    return count
}

func getSessionDir(t *testing.T, field string, value string) (string, string) {
    t.Helper()
    base := sessionsDir()
    today := time.Now()
    for i := 0; i < 7; i++ {
        dateDir := today.AddDate(0, 0, -i).Format("2006/01/02")
        dayPath := filepath.Join(base, dateDir)
        entries, err := os.ReadDir(dayPath)
        if err != nil {
            continue
        }
        for _, entry := range entries {
            if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "sess_") {
                continue
            }
            metaPath := filepath.Join(dayPath, entry.Name(), "meta.json")
            data, err := os.ReadFile(metaPath)
            if err != nil {
                continue
            }
            var m map[string]any
            if err := json.Unmarshal(data, &m); err != nil {
                continue
            }
            if v, ok := m[field]; ok {
                if s, ok := v.(string); ok && s == value {
                    return filepath.Join(dayPath, entry.Name()), metaPath
                }
            }
        }
    }
    return "", ""
}

func readEventsJSONL(t *testing.T, sessDir string) []map[string]any {
    t.Helper()
    eventsPath := filepath.Join(sessDir, "events.jsonl")
    data, err := os.ReadFile(eventsPath)
    if err != nil {
        t.Fatalf("cannot read events.jsonl: %v", err)
    }
    lines := strings.Split(strings.TrimSpace(string(data)), "\n")
    var events []map[string]any
    for _, line := range lines {
        if strings.TrimSpace(line) == "" {
            continue
        }
        var ev map[string]any
        if err := json.Unmarshal([]byte(line), &ev); err != nil {
            t.Fatalf("invalid JSON line in events.jsonl: %v\nline: %s", err, line)
        }
        events = append(events, ev)
    }
    return events
}
```
