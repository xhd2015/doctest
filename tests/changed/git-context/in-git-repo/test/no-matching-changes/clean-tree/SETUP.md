# Scenario

**Feature**: empty changed list yields zero selected leaves

```
# no modifications
ChangedFiles=[] -> ChangedCount 0 ; Announce false
```

## Steps

1. Create flat two-leaf tree.
2. Leave `ChangedFiles` empty.
3. Run filter policy.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	applyPolicyBase(req, fx)
	req.ChangedFiles = nil
	return nil
}
```
