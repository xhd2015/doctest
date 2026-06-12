package rules

func CheckChainHasRun(runSource string, assertPath string) *Violation {
	if runSource == "none" {
		return &Violation{Path: assertPath, Msg: "no Run(t *testing.T, req *Request) (*Response, error) in setup chain"}
	}
	return nil
}
