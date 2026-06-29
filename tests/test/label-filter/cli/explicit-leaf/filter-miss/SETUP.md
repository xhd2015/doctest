# Scenario

**Feature**: explicit leaf that fails label filter is skipped

## Steps

1. Run `test <mod>/slow --label heavy`.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	mod := writeLabelFilterMod(t)
	req.Args = []string{"test", filepath.Join(mod, "slow"), "--label", "heavy"}
	return nil
}
```