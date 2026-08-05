# Scenario

**Feature**: `--color` and `--no-color` together is an error

```
Harness -> list --color --no-color
  -> non-zero; cannot be specified together
```

## Steps

1. Args = `list --color --no-color` (no patterns needed for flag conflict).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"list", "--color", "--no-color"}
	return nil
}
```
