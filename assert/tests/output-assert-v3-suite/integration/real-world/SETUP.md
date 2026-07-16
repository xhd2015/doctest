# Scenario

**Feature**: Real-world CLI output cookbook (simulated transcripts)

```
# grouped by toolchain; each leaf asserts v3 template vs simulated bytes
Author -> Facade: version 3 templates
Matcher <- familiar CLI stdout/stderr shapes
```

## Steps
1. Category and leaf Setup functions build Template and Actual fields.

```go
func Setup(t *testing.T, req *Request) error {
	req.Operation = "match"
	return nil
}
```
