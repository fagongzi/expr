package expr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type mapVarExpr struct {
	attr string
}

func (v *mapVarExpr) AsString(ctx interface{}) string {
	m := ctx.(map[string]interface{})

	if v, ok := m[v.attr]; ok {
		return v.(string)
	}

	return ""
}

func (v *mapVarExpr) AsNumber(ctx interface{}) int64 {
	m := ctx.(map[string]interface{})

	if v, ok := m[v.attr]; ok {
		return v.(int64)
	}

	return 0
}

func TestParseSingle(t *testing.T) {
	input := []byte("{ num: key1 } == 1")
	expr, err := Parse(input, func(expr []byte) (VarExpr, error) {
		return &mapVarExpr{attr: string(expr)}, nil
	})
	assert.NoError(t, err, "TestParseSingle failed")
	assert.NotNil(t, expr, "TestParseSingle failed")

	m := make(map[string]interface{})
	assert.False(t, expr.Exec(m), "TestParseSingle failed")

	m["key1"] = int64(0)
	assert.False(t, expr.Exec(m), "TestParseSingle failed")

	m["key1"] = int64(1)
	assert.True(t, expr.Exec(m), "TestParseSingle failed")
}

func TestParseWithParen(t *testing.T) {
	input := []byte("( { num: key1 } == 1 ) || { str: key2 } == abc")
	expr, err := Parse(input, func(expr []byte) (VarExpr, error) {
		return &mapVarExpr{attr: string(expr)}, nil
	})
	assert.NoError(t, err, "TestParseWithParen failed")
	assert.NotNil(t, expr, "TestParseWithParen failed")

	m := make(map[string]interface{})
	assert.False(t, expr.Exec(m), "TestParseWithParen failed")

	m["key1"] = int64(0)
	m["key2"] = "abd"
	assert.False(t, expr.Exec(m), "TestParseWithParen failed")

	m["key2"] = "abc"
	assert.True(t, expr.Exec(m), "TestParseWithParen failed")
}

func TestParseWithParenNestNested(t *testing.T) {
	input := []byte("( ( { num: key1 } == 1 ) || { str: key2 } == abc ) && { key3 } == value3")
	expr, err := Parse(input, func(expr []byte) (VarExpr, error) {
		return &mapVarExpr{attr: string(expr)}, nil
	})
	assert.NoError(t, err, "TestParseWithParenNestNested failed")
	assert.NotNil(t, expr, "TestParseWithParenNestNested failed")

	m := make(map[string]interface{})
	assert.False(t, expr.Exec(m), "TestParseWithParenNestNested failed")

	m["key1"] = int64(0)
	m["key2"] = "abd"
	assert.False(t, expr.Exec(m), "TestParseWithParenNestNested failed")

	m["key2"] = "abc"
	assert.False(t, expr.Exec(m), "TestParseWithParenNestNested failed")

	m["key3"] = "value2"
	assert.False(t, expr.Exec(m), "TestParseWithParenNestNested failed")

	m["key3"] = "value3"
	assert.True(t, expr.Exec(m), "TestParseWithParenNestNested failed")

}
