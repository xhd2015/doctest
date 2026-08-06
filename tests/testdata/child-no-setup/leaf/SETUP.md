# Scenario

**Feature**: this leaf SETUP.md defines only a type declaration — no Setup, no Run

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

## Steps
1. This leaf SETUP.md defines only a type declaration — no Setup, no Run.
2. Non-root SETUP without func Setup is allowed (organization / prose-only).

```go
type ExtraType struct {
    Label string
}
```
