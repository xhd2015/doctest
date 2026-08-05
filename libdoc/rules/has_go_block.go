package rules

// CheckHasGoBlock formerly required every SETUP.md to contain a Go fence.
// SETUP.md Go blocks are optional (prose-only / organization nodes). Callers
// must not use this for SETUP.md; DOCTEST.md still requires a Go block via
// CheckRootHasGoBlock / explicit discover+validate checks on DOCTEST only.
//
// Kept as a no-op so old call sites cannot reintroduce the SETUP requirement.
func CheckHasGoBlock(goBlockSet bool, path string) *Violation {
	_ = goBlockSet
	_ = path
	return nil
}

func CheckGoBlockIsFinal(isFinal bool, path string) *Violation {
	if !isFinal {
		return &Violation{Path: path, Msg: "go block must be final content"}
	}
	return nil
}
