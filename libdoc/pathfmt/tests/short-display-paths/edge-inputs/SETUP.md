# Scenario

**Feature**: edge inputs for normalization and fallback behavior

```
# formatter pipeline
caller path string -> DisplayPath -> Abs normalize -> cwd/home rules -> display string

# fallback
otherwise -> absolute unchanged
```

## Preconditions

- Cwd is saved and restored; leaves configure edge-case paths.

## Steps

1. Save and restore the original cwd.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	saveAndRestoreCwd(t)
	return nil
}
```