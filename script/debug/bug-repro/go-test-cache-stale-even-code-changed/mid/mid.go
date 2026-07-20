package mid

// Version is the only behavioral knob. run.sh rewrites the return value
// between warm runs. Keep the body tiny so inlining is likely (mirrors
// small intermediate Setup packages in doctest gen).
func Version() int {
	return 1
}
