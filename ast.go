package expr

import (
	"fmt"
)

type valueType int
type cmp int
type logic int

var (
	stringValue = valueType(0)
	numberValue = valueType(1)
	regexpValue = valueType(2)
)

var (
	equal    = cmp(0) // ==
	notEqual = cmp(1) // !=
	gt       = cmp(2) // >
	ge       = cmp(3) // >=
	lt       = cmp(4) // <
	le       = cmp(5) // <=
	in       = cmp(6) // in
	notIn    = cmp(7) // !in
	match    = cmp(8) // ~
	notMatch = cmp(9) // !~
)

var (
	and = logic(0)
	or  = logic(1)
)

// Expr expr
type Expr interface {
	Exec(interface{}) bool
}

// VarExpr var expr
type VarExpr interface {
	AsString(interface{}) string
	AsNumber(interface{}) int64
}

// VarExprFactory factory method
type VarExprFactory func([]byte) (VarExpr, error)

type calcNode interface {
	Expr
	AddCmp(cmp, []byte) error
}

type node struct {
	exprs []Expr
	logic logic
}

func (n *node) lastExpr() Expr {
	return n.exprs[len(n.exprs)-1]
}

func (n *node) add(expr Expr) {
	n.exprs = append(n.exprs, expr)
}

func (n *node) append(logic logic, expr Expr) error {
	if len(n.exprs) > 1 && n.logic != logic {
		return fmt.Errorf("and/or can't mixin")
	}

	n.exprs = append(n.exprs, expr)
	n.logic = logic
	return nil
}

func (n *node) Exec(ctx interface{}) bool {
	if len(n.exprs) == 1 {
		return n.exprs[0].Exec(ctx)
	}

	var result bool
	for _, expr := range n.exprs {
		result = expr.Exec(ctx)
		switch n.logic {
		case and:
			if !result {
				return false
			}
		case or:
			if result {
				return true
			}
		}
	}
	return result
}
