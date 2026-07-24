# Scenario

**Feature**: HTTPS GitHub origin normalizes to host_owner_repo slug

```
# https origin
https://github.com/xhd2015/doctest.git -> github.com_xhd2015_doctest
```

## Preconditions

- Origin is a typical HTTPS remote including trailing `.git`.

## Steps

1. Set origin to `https://github.com/xhd2015/doctest.git`.
2. Resolve project id from origin.

## Context

- Scheme, optional credentials, and `.git` suffix must not appear in the slug.
- Path separators become underscores.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "project_id_from_origin"
	req.Origin = "https://github.com/xhd2015/doctest.git"
	return nil
}
```
