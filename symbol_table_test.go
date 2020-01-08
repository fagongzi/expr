package expr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddSynmbol(t *testing.T) {
	st := &symbolTable{}
	st.addSymbol([]byte("123"), 1)
	assert.Equal(t, 1, len(st.items), "TestAddSynmbol failed")
	st.addSymbol([]byte("124"), 2)
	assert.Equal(t, 1, len(st.items), "TestAddSynmbol failed")
	st.addSymbol([]byte("abc"), 3)
	assert.Equal(t, 2, len(st.items), "TestAddSynmbol failed")
	st.addSymbol([]byte("abd"), 4)
	assert.Equal(t, 2, len(st.items), "TestAddSynmbol failed")
}

func TestFindSynmbol(t *testing.T) {
	st := &symbolTable{}
	st.addSymbol([]byte("123"), 1)
	st.addSymbol([]byte("124"), 2)
	st.addSymbol([]byte("abc"), 3)
	st.addSymbol([]byte("abd"), 4)

	token, maybe := st.findToken([]byte("12"))
	assert.Equal(t, 0, token, "TestFindSynmbol failed")
	assert.True(t, maybe, "TestFindSynmbol failed")

	token, maybe = st.findToken([]byte("122"))
	assert.Equal(t, 0, token, "TestFindSynmbol failed")
	assert.False(t, maybe, "TestFindSynmbol failed")

	token, maybe = st.findToken([]byte("123"))
	assert.Equal(t, 1, token, "TestFindSynmbol failed")
	assert.True(t, maybe, "TestFindSynmbol failed")

	token, maybe = st.findToken([]byte("124"))
	assert.Equal(t, 2, token, "TestFindSynmbol failed")
	assert.True(t, maybe, "TestFindSynmbol failed")

	token, maybe = st.findToken([]byte("1241"))
	assert.Equal(t, 0, token, "TestFindSynmbol failed")
	assert.False(t, maybe, "TestFindSynmbol failed")

	token, maybe = st.findToken([]byte("ab"))
	assert.Equal(t, 0, token, "TestFindSynmbol failed")
	assert.True(t, maybe, "TestFindSynmbol failed")

	token, maybe = st.findToken([]byte("aba"))
	assert.Equal(t, 0, token, "TestFindSynmbol failed")
	assert.False(t, maybe, "TestFindSynmbol failed")

	token, maybe = st.findToken([]byte("abc"))
	assert.Equal(t, 3, token, "TestFindSynmbol failed")
	assert.True(t, maybe, "TestFindSynmbol failed")

	token, maybe = st.findToken([]byte("abd"))
	assert.Equal(t, 4, token, "TestFindSynmbol failed")
	assert.True(t, maybe, "TestFindSynmbol failed")

	token, maybe = st.findToken([]byte("abc1"))
	assert.Equal(t, 0, token, "TestFindSynmbol failed")
	assert.False(t, maybe, "TestFindSynmbol failed")
}
