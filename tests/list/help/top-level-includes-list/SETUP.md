# Scenario

**Feature**: top-level `doctest --help` lists the `list` command

```
Harness -> cli.RunWithWriters(["--help"])
  -> command list includes list
```

## Preconditions

- No fixture tree required.

## Steps

1. Set Args to `--help`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--help"}
	return nil
}
```
