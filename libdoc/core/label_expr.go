package core

import (
	"fmt"
	"strings"
	"unicode"
)

// EvalLabelExpr parses expr and returns whether labels satisfy the boolean expression.
// Evaluation is pure set membership: an empty label set (unlabeled leaf) makes every
// positive atom false, so expressions like "!e2e" match unlabeled leaves.
// Grammar: atoms, !, &&, ||, and parentheses. Precedence: ! > && > ||.
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

// labelBinary is && or || (name historical; was labelUnary).
type labelBinary struct {
	op    rune // '&' or '|' for && and ||
	left  labelExpr
	right labelExpr
}

func (u labelBinary) eval(set map[string]bool) bool {
	if u.op == '&' {
		return u.left.eval(set) && u.right.eval(set)
	}
	return u.left.eval(set) || u.right.eval(set)
}

// labelNot is unary bang negation: !expr
type labelNot struct {
	inner labelExpr
}

func (n labelNot) eval(set map[string]bool) bool {
	return !n.inner.eval(set)
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
		left = labelBinary{op: '|', left: left, right: right}
	}
}

func (p *labelParser) parseAnd() (labelExpr, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	for {
		p.skipSpace()
		if !p.consume("&&") {
			return left, nil
		}
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = labelBinary{op: '&', left: left, right: right}
	}
}

// parseNot handles unary ! (right-associative: !!e2e).
func (p *labelParser) parseNot() (labelExpr, error) {
	p.skipSpace()
	if p.pos < len(p.input) && p.peek() == '!' {
		p.pos++
		inner, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return labelNot{inner: inner}, nil
	}
	return p.parsePrimary()
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
		if r == '(' || r == ')' || r == '!' {
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
