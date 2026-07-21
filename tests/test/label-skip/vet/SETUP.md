# Scenario

**Feature**: `doctest vet` validates ASSERT.md frontmatter

```
# vet walks all leaves including labeled
doctest vet <tree> -> fail on malformed YAML
```

## Steps

1. Configure `req.Args` with `doctest vet` and a temp tree.

```go
import (
	"testing"
	"time"
)

func Setup(t *testing.T, req *Request) error {
	if req.Timeout == 0 {
		req.Timeout = 120 * time.Second
	}
	return nil
}
```