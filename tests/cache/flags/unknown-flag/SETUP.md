# Scenario

**Feature**: unknown flag on `cache` fails with non-zero exit and error hint

```
Harness -> cli.RunWithWriters(["cache","--not-a-real-flag"])
  -> non-zero; stderr usage/error
```

## Preconditions

- No fixture tree required.

## Steps

1. Set Args to `cache --not-a-real-flag`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"cache", "--not-a-real-flag"}
	return nil
}
```
