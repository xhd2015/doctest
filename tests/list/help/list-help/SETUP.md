# Scenario

**Feature**: `doctest list --help` documents usage, patterns, L2:L3, color flags

```
Harness -> cli.RunWithWriters(["list","--help"])
  -> usage text on stdout, exit 0
```

## Preconditions

- No fixture tree required.

## Steps

1. Set Args to `list --help`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"list", "--help"}
	return nil
}
```
