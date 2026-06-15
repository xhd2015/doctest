package rules

import "fmt"

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
