# Scenario

**Feature**: `--gen-dir` equal to warm `$CacheHome/doctest/mapping-gen` is rejected

```
# equal warm home
seed warmHome/marker-before
doctest test --cold-cache --gen-dir $warmHome <tree>
  -> error; marker preserved
```

## Preconditions

- Parent created sandbox + fixture; warm home path is known.

## Steps

1. Seed marker under warm home.
2. Run with `--gen-dir` exactly equal to warm home.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	seedMarker(t, st.WarmHome, "marker-before")
	req.Args = []string{"test", "--cold-cache", "--gen-dir", st.WarmHome, st.TestDir}
	return nil
}
```
