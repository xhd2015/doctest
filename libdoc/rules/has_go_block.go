package rules

func CheckHasGoBlock(goBlockSet bool, path string) *Violation {
	if !goBlockSet {
		return &Violation{Path: path, Msg: "must have a Go code block"}
	}
	return nil
}

func CheckGoBlockIsFinal(isFinal bool, path string) *Violation {
	if !isFinal {
		return &Violation{Path: path, Msg: "go block must be final content"}
	}
	return nil
}
