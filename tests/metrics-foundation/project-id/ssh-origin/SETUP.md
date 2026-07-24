# Scenario

**Feature**: SSH scp-like origin produces the same slug as HTTPS

```
# ssh origin
git@github.com:xhd2015/doctest.git -> github.com_xhd2015_doctest
```

## Preconditions

- Origin uses `git@host:path` form with trailing `.git`.

## Steps

1. Set origin to `git@github.com:xhd2015/doctest.git`.
2. Resolve project id from origin.

## Context

- SSH and HTTPS remotes for the same repo must share one metrics directory.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "project_id_from_origin"
	req.Origin = "git@github.com:xhd2015/doctest.git"
	return nil
}
```
