# Scenario

**Feature**: project_id from git origin URL or hashed workspace root

```
# origin present
origin URL -> strip scheme/user/.git -> host/path -> slugify / to _ -> project_id

# origin missing
abs root -> sha256 -> nogit_<12 hex>
```

## Preconditions

- Leaves set `req.Op` to `project_id_from_origin` or `project_id_fallback`.
- Origin strings are injectable (no live `git` required).

## Steps

1. Provide origin or abs root.
2. Call ProjectID helpers via root Run.
3. Assert exact slug or prefix+hex shape.

## Context

- Significance: identity is the top-level cache key under `metrics/<project_id>/`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Group default: project-id leaves exercise origin/fallback helpers only.
	t.Helper()
	if req.ProjectID != "" {
		// Clear path-oriented fields so identity leaves stay pure.
		req.ProjectID = ""
	}
	return nil
}
```
