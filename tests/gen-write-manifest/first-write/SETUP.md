# Scenario

**Feature**: cold gen root — first WriteGoMod establishes unified bookkeeping

```
# empty gen root
WriteGoMod(genDir, modRoot, ...)
  -> go.mod written
  -> doctest.gen-manifest created (lists go.mod)
  -> doctest.gomod-fp not created
```

## Preconditions

- Gen dir does not exist or is empty before Run.
- Parent module has a simple `go.mod`.

## Steps

1. Create fresh ModRoot + GenDir.
2. Run Mode `write-gomod` once (measured first write).

## Context

- Leaves assert different artifacts of the same first-write outcome (MECE on
  artifact under check: manifest presence vs legacy fp absence).

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "write-gomod"
	req.ModPath = "example.com/app"
	req.HasMod = true
	prepareFreshGen(t, req, "module example.com/app\n\ngo 1.21\n")
	return nil
}
```
