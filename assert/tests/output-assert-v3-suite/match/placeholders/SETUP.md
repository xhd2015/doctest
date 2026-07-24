# Scenario

**Feature**: Typed placeholder expansion in pattern lines

```
# __NAME__ replaced with type-specific subpatterns before match
Matcher <- actual with placeholder values
Matcher -> pass when values match type rules
```

## Steps
1. Templates declare placeholders in v3 YAML header.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = req
	return nil
}
```