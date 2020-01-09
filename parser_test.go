package expr

import (
	"fmt"
	"github.com/fagongzi/util/format"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testAdd(left interface{}, right Expr, ctx interface{}) (interface{}, error) {
	if _, ok := left.(int64); !ok {
		return nil, fmt.Errorf("%+v is not int64", left)
	}

	v2, err := right.Exec(ctx)
	if err != nil {
		return nil, err
	}

	if _, ok := v2.(int64); !ok {
		return nil, fmt.Errorf("%+v is not int64", v2)
	}

	return left.(int64) + v2.(int64), nil
}

func testEqual(left interface{}, right Expr, ctx interface{}) (interface{}, error) {
	if _, ok := left.(int64); !ok {
		return nil, fmt.Errorf("%+v is not int64", left)
	}

	v2, err := right.Exec(ctx)
	if err != nil {
		return nil, err
	}

	if _, ok := v2.(int64); !ok {
		return nil, fmt.Errorf("%+v is not int64", v2)
	}

	return left.(int64) == v2.(int64), nil
}

func testStrEqual(left interface{}, right Expr, ctx interface{}) (interface{}, error) {
	if _, ok := left.(string); !ok {
		return nil, fmt.Errorf("%+v is not string", left)
	}

	v2, err := right.Exec(ctx)
	if err != nil {
		return nil, err
	}

	if _, ok := v2.(string); !ok {
		return nil, fmt.Errorf("%+v is not string", v2)
	}

	return left.(string) == v2.(string), nil
}

func testAndLogic(left interface{}, right Expr, ctx interface{}) (interface{}, error) {
	if _, ok := left.(bool); !ok {
		return nil, fmt.Errorf("%+v is not bool", left)
	}

	if !left.(bool) {
		return false, nil
	}

	v2, err := right.Exec(ctx)
	if err != nil {
		return nil, err
	}

	if _, ok := v2.(bool); !ok {
		return nil, fmt.Errorf("%+v is not bool", v2)
	}

	return v2.(bool), nil
}

func testOrLogic(left interface{}, right Expr, ctx interface{}) (interface{}, error) {
	if _, ok := left.(bool); !ok {
		return nil, fmt.Errorf("%+v is not bool", left)
	}

	if left.(bool) {
		return true, nil
	}

	v2, err := right.Exec(ctx)
	if err != nil {
		return nil, err
	}

	if _, ok := v2.(bool); !ok {
		return nil, fmt.Errorf("%+v is not bool", v2)
	}

	return v2.(bool), nil
}

func TestParser(t *testing.T) {
	p := NewParser(nil)
	p.AddOP("+", testAdd)
	p.AddOP("==", testEqual)
	p.AddOP("===", testStrEqual)

	expr, err := p.Parse([]byte("((4+(1+2)+3)+5)==15"))
	assert.NoError(t, err, "TestParser failed")

	value, err := expr.Exec(nil)
	assert.NoError(t, err, "TestParser failed")
	assert.Equal(t, true, value, "TestParser failed")

	expr, err = p.Parse([]byte("abcd===abcd"))
	assert.NoError(t, err, "TestParser failed")
	assert.Equal(t, true, value, "TestParser failed")
}

func TestParserWithVar(t *testing.T) {
	p := NewParser(testVarFactory)
	p.AddOP("+", testAdd)
	p.AddOP("==", testEqual)
	p.AddOP("===", testStrEqual)
	p.AddOP("&&", testAndLogic)
	p.AddOP("||", testOrLogic)
	p.ValueType("num:", "str:")

	ctx := make(map[string]string)
	ctx["1"] = "1"
	ctx["2"] = "2"
	ctx["3"] = "3"
	ctx["4"] = "4"
	ctx["5"] = "5"

	expr, err := p.Parse([]byte("{1}+{2}"))
	assert.NoError(t, err, "TestParser failed")
	value, err := expr.Exec(ctx)
	assert.NoError(t, err, "TestParser failed")
	assert.Equal(t, int64(3), value, "TestParser failed")

	expr, err = p.Parse([]byte("(({4}+({1}+{2})+{3})+{5})==15"))
	assert.NoError(t, err, "TestParser failed")
	value, err = expr.Exec(ctx)
	assert.NoError(t, err, "TestParser failed")
	assert.Equal(t, true, value, "TestParser failed")

	expr, err = p.Parse([]byte("((({4}+({1}+{2})+{3})+{5})==15)&&(({1}+{2})==3)"))
	assert.NoError(t, err, "TestParser failed")
	value, err = expr.Exec(ctx)
	assert.NoError(t, err, "TestParser failed")
	assert.Equal(t, true, value, "TestParser failed")

	expr, err = p.Parse([]byte("((({4}+({1}+{2})+{3})+{5})==12)||(({1}+{2})==4)"))
	assert.NoError(t, err, "TestParser failed")
	value, err = expr.Exec(ctx)
	assert.NoError(t, err, "TestParser failed")
	assert.Equal(t, false, value, "TestParser failed")
}

type testMapBasedVarExpr struct {
	valueType string
	attr      string
}

func (expr *testMapBasedVarExpr) Exec(ctx interface{}) (interface{}, error) {
	m, ok := ctx.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("error ctx %T", ctx)
	}

	switch expr.valueType {
	case "str:":
		return m[expr.attr], nil
	case "num:":
		if v, ok := m[expr.attr]; ok {
			return format.ParseStrInt64(v)
		}
		return 0, nil
	}

	return nil, fmt.Errorf("not support value type")
}

func testVarFactory(value []byte, valueType string) (Expr, error) {
	return &testMapBasedVarExpr{
		valueType: valueType,
		attr:      string(value),
	}, nil
}
