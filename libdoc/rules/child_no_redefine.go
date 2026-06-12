package rules

import "fmt"

func CheckChildNoRedefine(types map[string]bool, path string, ancestorDepth int) *Violation {
	if ancestorDepth == 0 {
		return nil
	}
	if types["Request"] {
		return &Violation{Path: path, Msg: fmt.Sprintf("child SETUP.md cannot redefine Request")}
	}
	if types["Response"] {
		return &Violation{Path: path, Msg: fmt.Sprintf("child SETUP.md cannot redefine Response")}
	}
	return nil
}
