# Scenario

**Feature**: explicit labeled leaf runs when expression matches

## Steps

1. Run `test <mod>/slow --label slow`.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	mod := writeLabelFilterMod(t)
	req.Args = []string{"test", filepath.Join(mod, "slow"), "--label", "slow"}
	return nil
}
```