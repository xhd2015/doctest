# Unknown runner flag — L2 parse

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **system under test** — behavior covered by this tree.

### Behaviors
- **run** — executes the scenarios in this suite.

**Layer L2 in-process** — `runner.ParseTestOptions` rejects unrecognized flags.
No product binary, no `testbin`, no `label: heavy`.

# DSN

- **ParseTestOptions** — doctest `test` flag parser (args after subcommand).
- **Unknown flag** — less-flags returns `unrecognized flag: …`.

## How to Run

```sh
doctest vet ./tests/test/unknown-flag/
doctest test ./tests/test/unknown-flag/
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/runner"
)

// Request holds args for ParseTestOptions (no "test" subcommand prefix).
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

// Run parses flags in-process; map parse failure to ExitCode/Stderr for assert parity.
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_, _, err := runner.ParseTestOptions(req.Args)
	resp := &Response{}
	if err != nil {
		resp.ExitCode = 1
		resp.Stderr = err.Error() + "\n"
		resp.ParseErr = err.Error()
		resp.Err = err
		return resp, nil
	}
	return resp, nil
}
```
