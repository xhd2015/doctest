package droot

import "testing"

type Request struct{ WorkDir string }
type Response struct{}

func RootSetup(t *testing.T, req *Request) error { return nil }
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
