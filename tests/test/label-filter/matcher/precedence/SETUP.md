# Scenario

**Feature**: AND binds tighter than OR

```
a || b && c  ≡  a || (b && c)
```

## Steps

1. Compare match results for equivalent label sets.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = req
	return nil
}
```