# Scenario

**Feature**: V3S-M15 — gray color span pass

```
# <ansi-color gray> asserts ANSI gray wrap around inner text
Matcher <- gray-wrapped 1 Cached
```

## Steps
1. Set color span template and gray-wrapped actual.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("", "<ansi-color gray>1 Cached</ansi-color>")
	req.Actual = grayWrap("1 Cached")
	return nil
}
```