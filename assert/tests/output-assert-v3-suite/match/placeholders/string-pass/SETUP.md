# Scenario

**Feature**: V3S-M3 — string placeholder pass

```
# Hello __USER__ matches arbitrary same-line text
Matcher <- actual Hello alice
```

## Steps
1. Set USER string placeholder and matching actual.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("__USER__: type=string\n", "Hello __USER__")
	req.Actual = "Hello alice"
	return nil
}
```