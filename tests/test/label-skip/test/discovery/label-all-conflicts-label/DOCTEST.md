# --label-all vs --label mutual exclusion — L2 parse

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **system under test** — behavior covered by this tree.

### Behaviors
- **run** — executes the scenarios in this suite.

**Layer L2 in-process** — `runner.ParseTestOptions` rejects `--label-all` with `--label`.
No product binary, no `label: heavy`. Other `test/discovery/*` leaves stay L3 (skip summary format).

## How to Run

```sh
doctest vet ./tests/test/label-skip/test/discovery/label-all-conflicts-label/
doctest test ./tests/test/label-skip/test/discovery/label-all-conflicts-label/
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
