# Scenario

**Feature**: V3-P5 — custom regex= placeholder subpattern

```
# __ID__: regex=[a-z]+ custom fragment
Author -> v3 Parser: regex= placeholder
Parser -> Pattern with Placeholder{ID, regex subpattern}
```

## Steps
1. Set template with __ID__: regex=[a-z]+ and body id=__ID__.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("__ID__: regex=[a-z]+\n", "id=__ID__")
	return nil
}
```
