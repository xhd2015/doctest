## Expected

- Exit 0; stdout includes `Usage: doctest test` and `--label`.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
	}
	for _, want := range []string{"Usage: doctest test", "--label", "!"} {
		if !strings.Contains(resp.Stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
		}
	}
}
```
