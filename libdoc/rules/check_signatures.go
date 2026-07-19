package rules

import (
	"fmt"
	"strings"
)

func CheckSetupSignature(params, results string, path string) *Violation {
	if !matchOptionalDoctestParams(params, "t*testing.T,", ",req*Request") || normalize(results) != "error" {
		return &Violation{Path: path, Msg: "Setup must be func Setup(t *testing.T, req *Request) error or Setup(t *testing.T, d *session.Doctest, req *Request) error"}
	}
	return nil
}

func CheckRunSignature(params, results string, path string) *Violation {
	if !matchOptionalDoctestParams(params, "t*testing.T,", ",req*Request") || normalize(results) != "(*Response,error)" {
		return &Violation{Path: path, Msg: "Run must be func Run(t *testing.T, req *Request) (*Response, error) or Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error)"}
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
	if !matchOptionalDoctestParams(params, "t*testing.T,", ",req*Request,resp*Response,errerror") || normalize(results) != "" {
		return &Violation{Path: path, Msg: "Assert must be func Assert(t *testing.T, req *Request, resp *Response, err error) or Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error)"}
	}
	return nil
}

// matchOptionalDoctestParams accepts either without-d (prefix without trailing
// comma + suffix) or with-d (prefix + <name>*session.Doctest + suffix).
// The inject param name is free; type must be *session.Doctest.
// prefix ends with "," and suffix starts with "," (both normalized, no spaces).
func matchOptionalDoctestParams(params, prefix, suffix string) bool {
	n := normalize(params)
	withoutD := strings.TrimSuffix(prefix, ",") + suffix
	if n == withoutD {
		return true
	}
	if !strings.HasPrefix(n, prefix) || !strings.HasSuffix(n, suffix) {
		return false
	}
	mid := n[len(prefix) : len(n)-len(suffix)]
	if mid == "" || strings.Contains(mid, ",") {
		return false
	}
	return strings.HasSuffix(mid, "*session.Doctest") && len(mid) > len("*session.Doctest")
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
