## Preconditions
- The `Logf` function is in `libdoc/subagent` and writes timestamped output to `os.Stdout`.
- Each leaf provides a format string and optional args via environment variables.

## Steps
1. Read `LOGF_FORMAT` from `req.Env` as the format string (default: `"default"`).
2. Read `req.Args` as variadic format arguments (all strings).
3. Redirect `os.Stdout` to a pipe, call `subagent.Logf(format, args...)`.
4. Restore `os.Stdout`, read captured output, return as `resp.Stdout`.

```go
import (
    "bytes"
    "fmt"
    "os"
    "strings"
    "testing"

    "github.com/xhd2015/doctest/libdoc/subagent"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=logf")
    return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
    fmtStr := "default"
    for _, e := range req.Env {
        if strings.HasPrefix(e, "LOGF_FORMAT=") {
            fmtStr = e[len("LOGF_FORMAT="):]
            break
        }
    }

    var args []any
    for _, a := range req.Args {
        args = append(args, a)
    }

    old := os.Stdout
    r, w, err := os.Pipe()
    if err != nil {
        return nil, fmt.Errorf("create pipe: %w", err)
    }
    os.Stdout = w

    if len(args) > 0 {
        subagent.Logf("%s", fmt.Sprintf(fmtStr, args...))
    } else {
        subagent.Logf("%s", fmtStr)
    }

    w.Close()
    os.Stdout = old

    var buf bytes.Buffer
    if _, readErr := buf.ReadFrom(r); readErr != nil {
        return nil, fmt.Errorf("read pipe: %w", readErr)
    }

    return &Response{Stdout: buf.String()}, nil
}
```
