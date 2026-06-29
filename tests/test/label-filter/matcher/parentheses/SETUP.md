# Scenario

**Feature**: parentheses override default precedence

```
(slow || heavy) && ui  matches only when ui plus (slow or heavy)
```

## Steps

1. Assert grouped OR then AND.

```go
func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}
```