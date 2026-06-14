## Preconditions
- A --gen-dir is specified to control where output goes.

## Steps
1. Add --gen-dir to the build arguments so the generated output is written to a known location.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

var genDir string

func Setup(t *testing.T, req *Request) error {
	genDir = filepath.Join(t.TempDir(), "generated")
	if err := os.MkdirAll(genDir, 0755); err != nil {
		t.Fatalf("mkdir gen dir: %v", err)
	}
	req.Args = append(req.Args, "--gen-dir", genDir)
	return nil
}
```
