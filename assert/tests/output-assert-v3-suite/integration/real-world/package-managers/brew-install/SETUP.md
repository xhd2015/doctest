# Scenario

**Feature**: brew install

```
# brew install
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__PKG__: 'type=string, example=jq'\n",
		"==> Pouring jq--1\\.6\\.arm64_sonoma\\.bottle\\.tar\\.gz",
	)
	req.Actual = "==> Pouring jq--1.6.arm64_sonoma.bottle.tar.gz"
	return nil
}
```
