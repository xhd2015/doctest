# Scenario

**Feature**: relative input under cwd is absolutized then displayed as `"./rel"`

```
# formatter pipeline
caller path string -> DisplayPath -> Abs normalize -> cwd/home rules -> display string
```

## Steps

1. Set `req.Path` to the relative path `"child"` (not absolute).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Path = "child"
	return nil
}
```