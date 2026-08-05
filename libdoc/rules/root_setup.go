package rules

import "fmt"

// CheckRootHasGoBlock requires a Go fence on DOCTEST.md (not SETUP.md).
// SETUP.md may be prose-only; do not call this for intermediate SETUP paths.
func CheckRootHasGoBlock(goBlockSet bool, path string) *Violation {
	if !goBlockSet {
		return &Violation{Path: path, Msg: "must have a Go code block"}
	}
	return nil
}

func CheckRootHasRequestResponse(types map[string]bool, path string) *Violation {
	if !types["Request"] || !types["Response"] {
		return &Violation{Path: path, Msg: "must define type Request and type Response"}
	}
	return nil
}

func CheckRootHasRun(runSet bool, path string) *Violation {
	if !runSet {
		return &Violation{Path: path, Msg: fmt.Sprintf("must have func Run")}
	}
	return nil
}

func CheckRootSetupNoRequestResponseRun(types map[string]bool, runSet bool, path string) *Violation {
	if types["Request"] || types["Response"] {
		return &Violation{Path: path, Msg: "Request and Response must be defined in DOCTEST.md, not SETUP.md"}
	}
	if runSet {
		return &Violation{Path: path, Msg: "func Run must be defined in DOCTEST.md, not SETUP.md"}
	}
	return nil
}
