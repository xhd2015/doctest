# Scenario

**Feature**: V3-M3 — type=number loose numeric (1, -2, 3.14)

```
# type=number → -?\d+(?:\.\d+)? matches int, negative, float
Author -> v3 Matcher: three number placeholders
Matcher <- a=1 b=-2 c=3.14
Matcher -> pass
```

## Steps
1. Set three type=number placeholders and actuals 1, -2, 3.14.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__A__: type=number\n__B__: type=number\n__C__: type=number\n",
		"a=__A__\nb=__B__\nc=__C__",
	)
	req.Actual = "a=1\nb=-2\nc=3.14"
	return nil
}
```
