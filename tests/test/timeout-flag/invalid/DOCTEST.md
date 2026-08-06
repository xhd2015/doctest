# Invalid -timeout — L2 parse

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **system under test** — behavior covered by this tree.

### Behaviors
- **run** — executes the scenarios in this suite.

**Layer L2 in-process** — `runner.ParseTestOptions` rejects invalid `-timeout` values.
No product binary, no `label: e2e`. Sibling `forward/` stays L3 e2e (go test command line).

## How to Run

```sh
doctest vet ./tests/test/timeout-flag/invalid/
doctest test ./tests/test/timeout-flag/invalid/
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/runner"
)

type Request struct {
	Args []string
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	ParseErr string
	Err      error
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_, _, err := runner.ParseTestOptions(req.Args)
	resp := &Response{}
	if err != nil {
		resp.ParseErr = err.Error()
		resp.ExitCode = 1
		resp.Stderr = err.Error() + "\n"
		resp.Err = err
		return resp, nil
	}
	return resp, nil
}
```
