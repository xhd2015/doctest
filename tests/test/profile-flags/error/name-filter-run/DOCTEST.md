# Reject -run — L2 parse

## Version
0.0.1

## DSN (Domain Specific Notion)

### Participants
- **system under test** — behavior covered by this tree.

### Behaviors
- **run** — executes the scenarios in this suite.

**Layer L2 in-process** — `runner.ParseTestOptions` rejects name-based go test
filters (`-run`, `-skip`, `-bench`, …). Use path or `--label` instead.

## How to Run

```sh
doctest test ./tests/test/profile-flags/error/name-filter-run/
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
