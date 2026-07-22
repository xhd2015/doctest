# Scenario

**Feature**: sibling-branch intermediate SETUP edit does not bust peer leaf-cache

```
# two leaves on different branches
mid-a/leaf-x  other/leaf-z

# edit only other/SETUP.md (on leaf-z spine, not leaf-x)
# leaf-x key unchanged -> leaf-cache skip still hits for peer
doctest test <fixture> -> Cached > 0 (peer leaf-x stays cached)
```

## Preconditions
- mid-a/leaf-x and other/leaf-z; ASSERT observes WorkDir on both leaves.
- Multi-run harness from parent.
- Leaf key is spine-only: sibling SETUP is not part of the peer key.

## Steps
1. Edit only `other/SETUP.md` (sibling branch intermediate).
2. Expect peer `mid-a/leaf-x` to remain leaf-cached (`Cached` > 0), not a full suite `0 Cached`.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    testDir := createTempCustomProjectOpts(t, "mytest", treeOpts{ObserveWorkDir: true},
        []string{"mid-a", "other"},
        []string{"mid-a/leaf-x", "other/leaf-z"},
    )
    cfg := multiRunCfg{
        TestDir:       testDir,
        ModifyFile:    "other/SETUP.md",
        ModifyContent: modifiedSetupContent("modified-sibling-branch"),
    }
    doMultiRun(t, req, cfg)
    return nil
}
```
