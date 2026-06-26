# Scenario

**Feature**: Invalid templates fail at parse time

```
# invalid tag syntax rejected at parse time
Author -> Parser: malformed template
Parser -> parse error (position + message)
```

## Steps
1. Expect parse failure.

```go
func Setup(t *testing.T, req *Request) error {
	req.ExpectParseError = true
	return nil
}
```
