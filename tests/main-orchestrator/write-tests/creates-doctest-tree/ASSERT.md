## Expected
- DOCTEST.md, SETUP.md, basic/SETUP.md, and basic/ASSERT.md all exist.
- The root DOCTEST.md contains Request and Response types and a Run function.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	dir := req.SessionHome
	if dir == "" {
		t.Fatal("req.SessionHome (fixture tree) not set by SETUP")
	}

	for _, name := range []string{
		"DOCTEST.md",
		"SETUP.md",
		filepath.Join("basic", "SETUP.md"),
		filepath.Join("basic", "ASSERT.md"),
	} {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("missing file %s: %v", p, err)
		}
	}

	doctestMD, err := os.ReadFile(filepath.Join(dir, "DOCTEST.md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(doctestMD)
	if !strings.Contains(content, "type Request struct") {
		t.Fatal("root DOCTEST.md missing Request type")
	}
	if !strings.Contains(content, "func Run") {
		t.Fatal("root DOCTEST.md missing Run function")
	}
}
```
