package expr

import (
	"regexp"
	"strings"
)


type stringNode struct {
	value   string
	cmp     cmp
	varExpr VarExpr
	reg     *regexp.Regexp
}

func newStringNode() {

}

func (n *stringNode) AddCmp(cmp cmp, value []byte) error {
	n.value = string(value)
	n.cmp = cmp

	if cmp == match || cmp == notMatch {
		reg, err := regexp.Compile(n.value)
		if err != nil {
			return err
		}
		n.reg = reg
	}

	return nil
}

func (n *stringNode) Exec(ctx interface{}) bool {
	value := n.varExpr.AsString(ctx)

	switch n.cmp {
	case equal:
		return n.value == value
	case notEqual:
		return n.value != value
	case gt:
		return value > n.value
	case ge:
		return value >= n.value
	case lt:
		return value < n.value
	case le:
		return value <= n.value
	case in:
		return strings.Index(value, n.value) > 0
	case notIn:
		return strings.Index(value, n.value) == -1
	case match:
		return n.match(value)
	case notMatch:
		return n.notMatch(value)
	default:
		return false
	}
}

func (n *stringNode) match(value string) bool {
	return n.reg.MatchString(value)
}

func (n *stringNode) notMatch(value string) bool {
	return !n.reg.MatchString(value)
}
