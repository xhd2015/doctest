## Expected
- Exit code 0.
- `./local_subpath/...` from `workdir/` finds only `workdir/local_subpath/`.
- Does NOT find the module-root `parent_subpath/` tree.
- Stderr contains `local_subpath`, does NOT contain `parent_subpath`.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "local_subpath") {
        t.Fatalf("stderr missing local_subpath:\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "parent_subpath") {
        t.Fatalf("stderr should not contain parent_subpath (module-root sibling, not under working dir):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "PASS") {
        t.Fatalf("stdout missing PASS:\n%s", resp.Stdout)
    }
}
```
