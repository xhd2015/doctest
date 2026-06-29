package core

import (
	"fmt"
	"strings"
	"unicode"
)

// EvalLabelExpr parses expr and returns whether labels satisfy the boolean expression.
// Unlabeled leaves (len(labels)==0) never match a non-empty expression.
func EvalLabelExpr(expr string, labels []string) (bool, error) {
	ast, err := parseLabelExpr(expr)
	if err != nil {
		return false, err
	}
	set := labelSet(labels)
	return ast.eval(set), nil
}

// ParseLabelExpr validates a label expression without evaluating it.
func ParseLabelExpr(expr string) error {
	_, err := parseLabelExpr(expr)
	return err
}

// MatchLabelExprs returns true if any expression in exprs matches labels (OR across flags).
func MatchLabelExprs(exprs []string, labels []string) (bool, error) {
	if len(exprs) == 0 {
		return false, nil
	}
	set := labelSet(labels)
	for _, expr := range exprs {
		ast, err := parseLabelExpr(expr)
		if err != nil {
			return false, err
		}
		if ast.eval(set) {
			return true, nil
		}
	}
	return false, nil
}

func labelSet(labels []string) map[string]bool {
	set := make(map[string]bool, len(labels))
	for _, l := range labels {
		set[l] = true
	}
	return set
}

type labelExpr interface {
	eval(set map[string]bool) bool
}

type labelAtom struct {
	name string
}

func (a labelAtom) eval(set map[string]bool) bool {
	return set[a.name]
}

type labelUnary struct {
	op    rune // '&' or '|' for && and ||
	left  labelExpr
	right labelExpr
}

func (u labelUnary) eval(set map[string]bool) bool {
	if u.op == '&' {
		return u.left.eval(set) && u.right.eval(set)
	}
	return u.left.eval(set) || u.right.eval(set)
}

type labelParser struct {
	input string
	pos   int
}

func parseLabelExpr(expr string) (labelExpr, error) {
	p := &labelParser{input: strings.TrimSpace(expr)}
	if p.pos >= len(p.input) {
		return nil, fmt.Errorf("label expression parse error: empty expression")
	}
	ast, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	p.skipSpace()
	if p.pos < len(p.input) {
		return nil, fmt.Errorf("label expression syntax error: unexpected input %q", p.input[p.pos:])
	}
	return ast, nil
}

func (p *labelParser) parseOr() (labelExpr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for {
		p.skipSpace()
		if !p.consume("||") {
			return left, nil
		}
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = labelUnary{op: '|', left: left, right: right}
	}
}

func (p *labelParser) parseAnd() (labelExpr, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for {
		p.skipSpace()
		if !p.consume("&&") {
			return left, nil
		}
		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		left = labelUnary{op: '&', left: left, right: right}
	}
}

func (p *labelParser) parsePrimary() (labelExpr, error) {
	p.skipSpace()
	if p.pos >= len(p.input) {
		return nil, fmt.Errorf("label expression parse error: expected label or '('")
	}
	if p.peek() == '(' {
		p.pos++
		inner, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		p.skipSpace()
		if !p.consume(")") {
			return nil, fmt.Errorf("label expression syntax error: missing ')'")
		}
		return inner, nil
	}
	start := p.pos
	for p.pos < len(p.input) {
		r := rune(p.input[p.pos])
		if r == '(' || r == ')' {
			break
		}
		if p.matchAt("&&") || p.matchAt("||") {
			break
		}
		if unicode.IsSpace(r) {
			break
		}
		if !isLabelRune(r) {
			return nil, fmt.Errorf("label expression syntax error: invalid character %q in label", r)
		}
		p.pos++
	}
	if start == p.pos {
		return nil, fmt.Errorf("label expression parse error: expected label")
	}
	name := p.input[start:p.pos]
	return labelAtom{name: name}, nil
}

func isLabelRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_'
}

func (p *labelParser) skipSpace() {
	for p.pos < len(p.input) && unicode.IsSpace(rune(p.input[p.pos])) {
		p.pos++
	}
}

func (p *labelParser) peek() byte {
	if p.pos >= len(p.input) {
		return 0
	}
	return p.input[p.pos]
}

func (p *labelParser) consume(s string) bool {
	if !p.matchAt(s) {
		return false
	}
	p.pos += len(s)
	return true
}

func (p *labelParser) matchAt(s string) bool {
	return p.pos+len(s) <= len(p.input) && p.input[p.pos:p.pos+len(s)] == s
}