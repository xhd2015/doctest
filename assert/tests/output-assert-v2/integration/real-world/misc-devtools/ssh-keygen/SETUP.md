# Scenario

**Feature**: ssh-keygen

```
# ssh-keygen
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__PATH__: 'type=string, example=id_ed25519'\n",
		"Your identification has been saved in id_ed25519",
	)
	req.Actual = "Your identification has been saved in id_ed25519"
	return nil
}
```
