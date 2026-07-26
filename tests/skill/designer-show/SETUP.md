# Scenario

**Feature**: the designer prompt specification exists under `libdoc/designer`

```
# expose embedded designer prompt
doctest skill designer --show -> embedded PROMPT.md -> stdout
```

## Preconditions
- The designer prompt is registered as skill `designer`.

## Steps
1. Run `doctest skill designer --show`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = []string{"skill", "designer", "--show"}
    return nil
}
```
