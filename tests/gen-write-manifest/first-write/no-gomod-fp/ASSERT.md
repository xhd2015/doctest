## Expected

- `Run` succeeds.
- `doctest.gomod-fp` does **not** exist under gen root.
- Unified `doctest.gen-manifest` **does** exist (skip key moved here).

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteGoMod failed: %v", err)
	}
	if resp.GomodFpExists {
		t.Fatalf("legacy %s must not be written after unified manifest lands", gomodFpName)
	}
	if !resp.ManifestExists {
		t.Fatalf("expected %s (replacement skip index) at gen root", genManifestName)
	}
}
```
