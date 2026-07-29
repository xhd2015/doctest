# Scenario

**Feature**: invalid `--go-cmd` values are rejected with a clear error

```
ParseTestOptions([--go-cmd=foo, .])
  -> error mentioning go-cmd and invalid/allowed values
  -> non-zero exit mapping
```

## Preconditions

- Parse-time rejection preferred (same pattern as invalid `-timeout`).
- Must not be confused with a silent default to auto.

## Steps

1. Leaves set `ParseOnly` and bad `--go-cmd` args.
2. Assert non-zero exit and clear message.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.ParseOnly = true
	req.Detect = false
	req.CheckAvailable = false
	return nil
}
```
