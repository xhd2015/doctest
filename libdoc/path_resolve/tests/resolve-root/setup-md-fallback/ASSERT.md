## Expected
- `resp.RootOkResult` is `true`.
- `resp.RootResult` is the absolute path of the directory containing `SETUP.md`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if !resp.RootOkResult {
		t.Fatal("expected ok == true")
	}
	fi, statErr := os.Stat(resp.RootResult)
	if statErr != nil {
		t.Fatalf("root %s not accessible: %v", resp.RootResult, statErr)
	}
	if !fi.IsDir() {
		t.Fatalf("root %s is not a directory", resp.RootResult)
	}
	if _, statErr := os.Stat(filepath.Join(resp.RootResult, "SETUP.md")); statErr != nil {
		t.Fatalf("root %s does not contain SETUP.md: %v", resp.RootResult, statErr)
	}
}
```
