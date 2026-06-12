package rules

import "fmt"

func CheckSetupSignature(params, results string, path string) *Violation {
	if normalize(params) != "t*testing.T,req*Request" || normalize(results) != "error" {
		return &Violation{Path: path, Msg: "Setup must be func Setup(t *testing.T, req *Request) error"}
	}
	return nil
}

func CheckRunSignature(params, results string, path string) *Violation {
	if normalize(params) != "t*testing.T,req*Request" || normalize(results) != "(*Response,error)" {
		return &Violation{Path: path, Msg: "Run must be func Run(t *testing.T, req *Request) (*Response, error)"}
	}
	return nil
}

func CheckAssertExists(assertSet bool, path string) *Violation {
	if !assertSet {
		return &Violation{Path: path, Msg: fmt.Sprintf("missing func Assert(t *testing.T, req *Request, resp *Response, err error)")}
	}
	return nil
}

func CheckAssertSignature(params, results string, path string) *Violation {
	if normalize(params) != "t*testing.T,req*Request,resp*Response,errerror" || normalize(results) != "" {
		return &Violation{Path: path, Msg: "Assert must be func Assert(t *testing.T, req *Request, resp *Response, err error)"}
	}
	return nil
}

func normalize(s string) string {
	n := len(s)
	b := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		c := s[i]
		if c == ' ' || c == '\t' || c == '\n' {
			continue
		}
		b = append(b, c)
	}
	return string(b)
}
