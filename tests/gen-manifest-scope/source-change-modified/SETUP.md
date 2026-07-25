# Scenario

**Feature**: changing source leaf content re-emits gen file as `# modified`

```
run1: test tree
rewrite tree/leaf/ASSERT.md (source)
run2: same + gen-plan
  -> tree/leaf/leaf.go # modified (or modified>=1)
```

```go
import (
	"strings"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	prepareSingleTreeModule(t, req)
	args := baseArgs(req, "tree")
	req.ArgsFull = args
	req.ArgsSubset = append([]string(nil), args...)
	req.DebugEnv = debugGenPlanBypass
	// Build ASSERT.md without embedding markdown fences in this SETUP (vet).
	// Use string(rune) to form fence at runtime so the go block stays last.
	tick := string([]byte{0x60, 0x60, 0x60}) // three backticks
	assertBody := strings.Join([]string{
		"## Expected",
		"- changed-expected-marker",
		"",
		tick + "go",
		"import (",
		"\t\"testing\"",
		"",
		"\t\"github.com/xhd2015/doctest/session\"",
		")",
		"",
		"func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {",
		"\t_ = d",
		"\t// changed-expected-marker",
		"}",
		tick,
		"",
	}, "\n")
	req.RewriteSourceRels = map[string]string{
		"tree/leaf/ASSERT.md": assertBody,
	}
	return nil
}
```
