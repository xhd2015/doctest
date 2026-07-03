package legacy_v1

import "fmt"

type parseError struct {
	pos int
	msg string
}

func (e *parseError) Error() string {
	if e.pos >= 0 {
		return fmt.Sprintf("parse error at position %d: %s", e.pos, e.msg)
	}
	return fmt.Sprintf("parse error: %s", e.msg)
}

func parseErr(pos int, format string, args ...any) error {
	return &parseError{pos: pos, msg: fmt.Sprintf(format, args...)}
}

type matchError struct {
	msg string
}

func (e *matchError) Error() string {
	return e.msg
}

func matchErr(format string, args ...any) error {
	return &matchError{msg: fmt.Sprintf(format, args...)}
}