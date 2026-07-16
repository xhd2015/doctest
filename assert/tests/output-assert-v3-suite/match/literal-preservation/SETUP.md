# Scenario

**Feature**: Escaped literals under v3 raw-RE content lines

```
# authors escape metachars so CLI-looking text matches literally
Matcher <- version 1\.0 / cost: \$5\.00 / \(1 Cached\)
```

## Steps
1. Leaf templates use escaped RE metacharacters; actual bytes stay unescaped.

```go
func Setup(t *testing.T, req *Request) error {
	req.Operation = "match"
	return nil
}
```
