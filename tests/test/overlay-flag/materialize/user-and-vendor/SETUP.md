# Scenario

**Feature**: user seed ∪ vendor-gomod overlay without pre_test (later vendor wins same key)

```
MaterializeUserVendorOverlay(user, vendor-gomod-overlay.json, dest)
  -> single -overlay= with union Replace (vendor wins conflicts)
```

## Preconditions

- No pre_test hooks (`PreTest` empty).
- `UseMaterializeHelper` true.
- Vendor overlay is the gen `vendor-gomod-overlay.json` style Replace map.

## Steps

1. Parent sets materialize mode + materialize helper.
2. Leaves set same-key vs disjoint Replace pairs.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.UseMaterializeHelper = true
	req.PreTest = nil
	return nil
}
```
