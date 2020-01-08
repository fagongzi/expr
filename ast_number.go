package expr

import (
	"github.com/fagongzi/util/format"
)

type numberNode struct {
	value   int64
	cmp     cmp
	varExpr VarExpr
}

func (n *numberNode) AddCmp(cmp cmp, value []byte) error {
	v, err := format.ParseStrInt64(string(value))
	if err != nil {
		return err
	}

	n.value = v
	n.cmp = cmp

	return nil
}

func (n *numberNode) Exec(ctx interface{}) bool {
	value := n.varExpr.AsNumber(ctx)

	switch n.cmp {
	case equal:
		return value == n.value
	case notEqual:
		return value != n.value
	case gt:
		return value > n.value
	case ge:
		return value >= n.value
	case lt:
		return value < n.value
	case le:
		return value <= n.value
	default:
		return false
	}
}
