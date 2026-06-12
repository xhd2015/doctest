package rules

import "fmt"

func CheckNoHelperRedefinition(ancestorHelpers map[string]bool, childHelpers []string, path string, ancestorDepth int) *Violation {
	if ancestorDepth == 0 {
		return nil
	}
	for _, name := range childHelpers {
		if ancestorHelpers[name] {
			return &Violation{
				Path: path,
				Msg:  fmt.Sprintf("helper %q already defined by an ancestor SETUP.md — use a unique name (e.g., with a suffix)", name),
			}
		}
	}
	return nil
}
