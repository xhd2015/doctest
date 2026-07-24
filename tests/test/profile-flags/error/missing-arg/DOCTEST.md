# Missing -cpuprofile value — L2 parse

## Version
0.0.2

**Layer L2 in-process** — `runner.ParseTestOptions` rejects `-cpuprofile` without a value.
No product binary, no `label: heavy`. Forward/side-effect siblings stay L3 e2e.

## How to Run

```sh
doctest vet ./tests/test/profile-flags/error/missing-arg/
doctest test ./tests/test/profile-flags/error/missing-arg/
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
