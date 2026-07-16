# Scenario

**Feature**: Invalid v3 templates fail at parse time

```
# invalid v3 syntax rejected by parser
Author -> Facade: malformed v3 template
Facade -> parse error
```

## Steps
1. Expect parse failure.

```go
func Setup(t *testing.T, req *Request) error {
	req.ExpectParseError = true
	return nil
}
```
