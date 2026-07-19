# Scenario

**Feature**: ref-instead-of-inline assemble (`AssembleRef*`) uses the same inject contract

```
# ref path
TreeCase -> AssembleRefRootSource + AssembleRefLeafTestSource
  -> thin leaf test constructs d, no Chdir, no free DOCTEST_* vars
  -> root package also drops free DOCTEST_* package vars
```

## Preconditions

- `req.Op = "ref"`.

## Steps

1. Set Op to ref for all descendants.

## Context

- Root package source is returned as `resp.RootSrc`; leaf as `resp.Source`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "ref"
	return nil
}
```
