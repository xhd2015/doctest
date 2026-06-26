# Scenario

**Feature**: Realistic CLI output integration templates

```
# realistic doctest CLI output templates
Author -> Matcher: multi-construct template
Matcher <- simulated build/help output
```

## Steps
1. Templates mirror doctest build/help output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Operation = "match"
	return nil
}
```
