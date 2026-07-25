# Scenario

**Feature**: multi-arg plan has package-only per-arg trees and merged bookkeeping

```
arg[1/2] / arg[2/2]: no repeated go.mod each
merged: bookkeeping + both trees [+ __workspace]
```

## Preconditions

- Parent multi-arg fixture ready.

## Steps

1. Assert per-arg and merged markers on stderr.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = req
	return nil
}
```
