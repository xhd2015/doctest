# Scenario

**Feature**: `doctest test` parses `-overlay` / `--overlay` and abs-resolves paths

```
ParseTestOptions([-overlay|--overlay FILE, dir])
  -> Options.Overlay absolute | parse error
```

## Preconditions

- Parse-only L2; no materialize helpers required for these leaves.
- Directory remain is a dummy path (parse does not require a real tree).

## Steps

1. Set `Mode=parse`.
2. Leaf sets `ParseArgs` for the form under test.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = modeParse
	return nil
}
```
