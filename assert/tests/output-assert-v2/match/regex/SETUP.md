# Scenario

**Feature**: Regex-intent template lines

```
# strong-signal lines compiled as Go regexp full-line match
Matcher <- actual line must fully match regex
```

## Steps
1. Body lines use regex-intent patterns (e.g. `.*`, alternation).

```go
func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}
```