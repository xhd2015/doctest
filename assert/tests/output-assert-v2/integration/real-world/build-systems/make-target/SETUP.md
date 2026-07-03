# Scenario

**Feature**: make

```
# make
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"make: Entering directory '/tmp'\n...1 lines omitted...\nmake: Leaving directory '/tmp'",
	)
	req.Actual = "make: Entering directory '/tmp'\ngcc -o app main.c\nmake: Leaving directory '/tmp'"
	return nil
}
```
