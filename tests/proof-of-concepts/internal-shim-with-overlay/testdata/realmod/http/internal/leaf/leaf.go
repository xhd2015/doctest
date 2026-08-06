package leaf

// Hello is only reachable from packages under .../http/...
func Hello() string { return "from-internal-leaf" }
