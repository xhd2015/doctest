# Invalid --label expression — L2 parse

## Version
0.0.2

**Layer L2 in-process** — `runner.ParseTestOptions` rejects trailing `&&` in label EXPR.
No product binary, no `label: heavy`.

## How to Run

```sh
doctest vet ./tests/test/label-filter/cli/parse-error/trailing-and/
doctest test ./tests/test/label-filter/cli/parse-error/trailing-and/
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

func Run(t *testing.T, req *Request) (*Response, error) {
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
