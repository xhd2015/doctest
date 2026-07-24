# Scenario

**Feature**: CasesImportAssertPackage detects assert import path in case Go blocks

```
# parsed import path exact match
"github.com/xhd2015/doctest/assert" -> true (alias name ignored)
```

## Preconditions

- Siblings cover direct import, alias import, and absent import.

## Steps

1. Descendant sets `runKind = "detect"` and builds `req.Cases`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.RunKind = "detect"
	req.ModPath = "example.com/app"
	return nil
}
```