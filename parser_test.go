package expr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/fagongzi/util/format"
	"github.com/stretchr/testify/assert"
)

func testIn(left interface{}, right Expr, ctx interface{}) (interface{}, error) {
	if _, ok := left.(int64); !ok {
		return nil, fmt.Errorf("%+v is not int64", left)
	}

	v2, err := right.Exec(ctx)
	if err != nil {
		return nil, err
	}

	if _, ok := v2.([]string); !ok {
		return nil, fmt.Errorf("%+v is not []string", v2)
	}

	expect := left.(int64)
	for _, v := range v2.([]string) {
		vn, err := format.ParseStrInt64(v)
		if err != nil {
			return false, err
		}

		if vn == expect {
			return true, nil
		}
	}

	return false, nil
}

func testStrIn(left interface{}, right Expr, ctx interface{}) (interface{}, error) {
	if _, ok := left.(string); !ok {
		return nil, fmt.Errorf("%+v is not string", left)
	}

	v2, err := right.Exec(ctx)
	if err != nil {
		return nil, err
	}

	if _, ok := v2.([]string); !ok {
		return nil, fmt.Errorf("%+v is not []string", v2)
	}

	expect := left.(string)
	for _, v := range v2.([]string) {
		if v == expect {
			return true, nil
		}
	}

	return false, nil
}

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

func testMatch(left interface{}, right Expr, ctx interface{}) (interface{}, error) {
	if _, ok := left.(string); !ok {
		return nil, fmt.Errorf("expect string left value but %T", left)
	}

	rightValue, err := right.Exec(ctx)
	if err != nil {
		return nil, err
	}

	if _, ok := rightValue.(*regexp.Regexp); !ok {
		return nil, fmt.Errorf("expect regexp right value but %T", rightValue)
	}

	return rightValue.(*regexp.Regexp).MatchString(left.(string)), nil
}

func TestParser(t *testing.T) {
	p := NewParser(nil,
		WithOp("+", testAdd),
		WithOp("==", testEqual),
		WithOp("===", testStrEqual),
		WithVarType("num:", Num),
		WithVarType("str:", Str))

	expr, err := p.Parse([]byte("((4+(1+2)+3)+5)==15"), nil)
	assert.NoError(t, err, "TestParser failed")

	value, err := expr.Exec(nil)
	assert.NoError(t, err, "TestParser failed")
	assert.Equal(t, true, value, "TestParser failed")

	expr, err = p.Parse([]byte("abcd===abcd"), nil)
	assert.NoError(t, err, "TestParser failed")
	assert.Equal(t, true, value, "TestParser failed")
}

func TestParserWithVarAndLiteral(t *testing.T) {
	p := NewParser(testVarFactory,
		WithOp("==", testStrEqual),
		WithVarType("num:", Num),
		WithVarType("str:", Str))

	ctx := make(map[string]string)
	ctx["1"] = `{\"abc}`

	expr, err := p.Parse([]byte(`{str:1}=="{\\\"abc}"`), nil)
	assert.NoError(t, err, "TestParser failed")

	value, err := expr.Exec(ctx)
	assert.NoError(t, err, "TestParser failed")
	assert.Equal(t, true, value, "TestParser failed")

	expr, err = p.Parse([]byte(`"{\\\"abc}"=={str:1}`), nil)
	assert.NoError(t, err, "TestParser failed")

	value, err = expr.Exec(ctx)
	assert.NoError(t, err, "TestParser failed")
	assert.Equal(t, true, value, "TestParser failed")
}

func TestParserWithVar(t *testing.T) {
	p := NewParser(testVarFactory,
		WithOp("+", testAdd),
		WithOp("==", testEqual),
		WithOp("===", testStrEqual),
		WithOp("&&", testAndLogic),
		WithOp("||", testOrLogic),
		WithVarType("num:", Num),
		WithVarType("str:", Str),
		WithDefaultVarType(Num))

	ctx := make(map[string]string)
	ctx["1"] = "1"
	ctx["2"] = "2"
	ctx["3"] = "3"
	ctx["4"] = "4"
	ctx["5"] = "5"

	expr, err := p.Parse([]byte("{1}+{2}"), nil)
	assert.NoError(t, err, "TestParser failed")
	value, err := expr.Exec(ctx)
	assert.NoError(t, err, "TestParser failed")
	assert.Equal(t, int64(3), value, "TestParser failed")

	expr, err = p.Parse([]byte("(({4}+({1}+{2})+{3})+{5})==15"), nil)
	assert.NoError(t, err, "TestParser failed")
	value, err = expr.Exec(ctx)
	assert.NoError(t, err, "TestParser failed")
	assert.Equal(t, true, value, "TestParser failed")

	expr, err = p.Parse([]byte("((({4}+({1}+{2})+{3})+{5})==15)&&(({1}+{2})==3)"), nil)
	assert.NoError(t, err, "TestParser failed")
	value, err = expr.Exec(ctx)
	assert.NoError(t, err, "TestParser failed")
	assert.Equal(t, true, value, "TestParser failed")

	expr, err = p.Parse([]byte("((({4}+({1}+{2})+{3})+{5})==12)||(({1}+{2})==4)"), nil)
	assert.NoError(t, err, "TestParser failed")
	value, err = expr.Exec(ctx)
	assert.NoError(t, err, "TestParser failed")
	assert.Equal(t, false, value, "TestParser failed")
}

func TestParserRegexpWithVar(t *testing.T) {
	p := NewParser(testVarFactory,
		WithOp("~", testMatch),
		WithVarType("num:", Num),
		WithVarType("str:", Str),
		WithVarType("regexp:", Regexp))

	ctx := make(map[string]string)
	ctx["1"] = "||||"

	expr, err := p.Parse([]byte("{str:1}~|^[\\|]+$|"), nil)
	assert.NoError(t, err, "TestParserRegexpWithVar failed")
	value, err := expr.Exec(ctx)
	assert.NoError(t, err, "TestParserRegexpWithVar failed")
	assert.Equal(t, true, value, "TestParserRegexpWithVar failed")
}

func TestParserArrayWithVar(t *testing.T) {
	p := NewParser(testVarFactory,
		WithOp("in", testStrIn),
		WithVarType("num:", Num),
		WithVarType("str:", Str))

	ctx := make(map[string]string)
	ctx["1"] = "|"

	expr, err := p.Parse([]byte(`{str:1} in [\\1,\|,3]`), nil)
	assert.NoError(t, err, "TestParserArrayWithVar failed")
	value, err := expr.Exec(ctx)
	assert.NoError(t, err, "TestParserArrayWithVar failed")
	assert.Equal(t, true, value, "TestParserArrayWithVar failed")

	expr, err = p.Parse([]byte("{str:1} in [4,2,3]"), nil)
	assert.NoError(t, err, "TestParserArrayWithVar failed")
	value, err = expr.Exec(ctx)
	assert.NoError(t, err, "TestParserArrayWithVar failed")
	assert.Equal(t, false, value, "TestParserArrayWithVar failed")
}

func TestConversionAndRevert(t *testing.T) {
	value := conversion([]byte(`"`))
	assert.Equal(t, []byte(`"`), value, "TestConversion failed")
	assert.Equal(t, []byte(`"`), revertConversion(value), "TestConversion failed")

	value = conversion([]byte(`""`))
	assert.Equal(t, []byte(`""`), value, "TestConversion failed")
	assert.Equal(t, []byte(`""`), revertConversion(value), "TestConversion failed")

	value = conversion([]byte(`\"`))
	assert.Equal(t, []byte{quotationConversion}, value, "TestConversion failed")
	assert.Equal(t, []byte(`"`), revertConversion(value), "TestConversion failed")

	value = conversion([]byte(`b\"a`))
	assert.Equal(t, []byte{'b', quotationConversion, 'a'}, value, "TestConversion failed")
	assert.Equal(t, []byte(`b"a`), revertConversion(value), "TestConversion failed")

	value = conversion([]byte(`\\`))
	assert.Equal(t, []byte{slashConversion}, value, "TestConversion failed")
	assert.Equal(t, []byte(`\`), revertConversion(value), "TestConversion failed")

	value = conversion([]byte(`\\\"`))
	assert.Equal(t, []byte{slashConversion, quotationConversion}, value, "TestConversion failed")
	assert.Equal(t, []byte(`\"`), revertConversion(value), "TestConversion failed")

	value = conversion([]byte(`\\\"\`))
	assert.Equal(t, []byte{slashConversion, quotationConversion, '\\'}, value, "TestConversion failed")
	assert.Equal(t, []byte(`\"\`), revertConversion(value), "TestConversion failed")
}

type testMapBasedVarExpr struct {
	valueType VarType
	attr      string
}

func (expr *testMapBasedVarExpr) Exec(ctx interface{}) (interface{}, error) {
	m, ok := ctx.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("error ctx %T", ctx)
	}

	return ValueByType([]byte(m[expr.attr]), expr.valueType)
}

func testVarFactory(value []byte, valueType VarType) (Expr, error) {
	return &testMapBasedVarExpr{
		valueType: valueType,
		attr:      string(value),
	}, nil
}
