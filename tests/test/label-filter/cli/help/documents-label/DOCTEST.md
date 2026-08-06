# test --help documents --label — L2 CLI

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **system under test** — behavior covered by this tree.

### Behaviors
- **run** — executes the scenarios in this suite.

**Layer L2 in-process** — `cli.RunWithWriter` for help text.
No product binary, no `label: heavy`.

## How to Run

```sh
doctest vet ./tests/test/label-filter/cli/help/documents-label/
doctest test ./tests/test/label-filter/cli/help/documents-label/
```

```go
import (
	"bytes"
	"testing"

	"github.com/xhd2015/doctest/libdoc/cli"
)

type Request struct {
	Args []string
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

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
