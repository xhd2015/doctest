# Scenario

**Feature**: empty path input is returned unchanged

```
# normalize
empty or Abs error -> return input unchanged
```

## Steps

1. Set `req.Path` to `""`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Path = ""
	return nil
}
```