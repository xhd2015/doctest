---
label: heavy
---

## Expected

- `DiscoverTreeCases` returns **nil** error.
- Exactly **1** case (`own_leaf`); nested DOCTEST under intermediate is not a parent case.
- DiscoverErr must not mention `must have a Go code block` for intermediate SETUP.

## Errors

- Today: full discover walks intermediate/ and reports missing SETUP Go block.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	_ = req
	if err != nil {
		t.Fatalf("harness error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.DiscoverErr != "" {
		t.Fatalf("DiscoverTreeCases want nil error for pure nested parent, got:\n%s", resp.DiscoverErr)
	}
	if resp.CaseCount != 1 {
		t.Fatalf("want 1 parent case (own_leaf), got %d", resp.CaseCount)
	}
	if len(resp.Cases) == 1 {
		p := resp.Cases[0].Path
		if p != "own_leaf" && !strings.HasSuffix(p, "own_leaf") {
			t.Fatalf("expected sole case path own_leaf, got %q", p)
		}
	}
}
```
