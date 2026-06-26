# Scenario

**Feature**: P15 — unknown ansi color name

```
# invalid tag syntax rejected at parse time
Author -> Parser: malformed template
Parser -> parse error (position + message)
```

## Steps
1. Set template/actual fields for P15 — unknown ansi color name.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<ansi-color orange>x</ansi-color>"
	return nil
}
```
