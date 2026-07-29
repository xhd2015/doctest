# Scenario

**Feature**: path-shaped TranslatePath contract (mid + nested go.mod)

**Phase 1 (locked here):** pure `gotestmap.TranslatePath` / `Plan(ModePathShaped)`
rules — mid pattern + nested go.mod cmds, no widen to parent siblings.

**Phase 2 (not this leaf):** path-shaped *execution* via multi-cmd
`finishWorkspaceGoTestCmds`. Production CLI still uses single workspace/hub
suite until that executor exists — do not require multi-cmd finish here.

CLI smoke: single-mod `./alpha/...` → still **one** workspace suite plan.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Minimal CLI touch so leaf is a full doctest run; plan matrix is pure Assert.
	req.WorkDir = createSingleModTwoTrees(t)
	req.Args = []string{"test", "-v", "./alpha/..."}
	return nil
}
```
