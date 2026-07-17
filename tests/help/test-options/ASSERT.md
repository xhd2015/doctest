## Expected
- The command succeeds.
- stdout includes test runner options.

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
        t.Fatalf("exit code = %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    for _, want := range []string{
        "Usage: doctest test",
        "-v", "--verbose", "--rm", "-count", "--timeout", "--color", "--no-color",
        // Go-style profiling / cover flags forwarded to go test
        "-cpuprofile", "-memprofile", "-memprofilerate",
        "-blockprofile", "-blockprofilerate",
        "-mutexprofile", "-mutexprofilefraction",
        "-trace", "-outputdir",
        "-coverprofile", "-cover",
    } {
        if !strings.Contains(resp.Stdout, want) {
            t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
        }
    }
}
```
