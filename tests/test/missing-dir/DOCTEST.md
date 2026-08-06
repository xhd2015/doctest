# Missing test directory — L2 in-process CLI

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **system under test** — behavior covered by this tree.

### Behaviors
- **run** — executes the scenarios in this suite.

**Layer L2 in-process** — `cli.RunWithWriter` / `runner.Test` path requires `<dir>`.
No product binary, no `testbin`, no `label: e2e`.

# DSN

- **CLI `test`** — without a directory operand, fails with `test requires <dir>`.

## How to Run

```sh
doctest vet ./tests/test/missing-dir/
doctest test ./tests/test/missing-dir/
```

```go
import (
	"bytes"
	"testing"

	"github.com/xhd2015/doctest/libdoc/cli"
)

// Request drives one CLI scenario. Leaves set Args only.
type Request struct {
	Args []string // e.g. ["test"]
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

// Run dispatches in-process via cli.RunWithWriter (captures help/stdout text).
// Errors from runner map to ExitCode/Stderr like the product binary main.
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	var buf bytes.Buffer
	err := cli.RunWithWriter(&buf, req.Args)
	resp := &Response{
		Stdout: buf.String(),
		Err:    err,
	}
	if err != nil {
		resp.ExitCode = 1
		resp.Stderr = err.Error() + "\n"
		return resp, nil
	}
	return resp, nil
}
```
