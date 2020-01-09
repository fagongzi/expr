package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkipWhitespaces(t *testing.T) {
	scan := NewScanner([]byte("  \r\n && 1&2   ||"))
	scan.(*scanner).skipWhitespaces()
	assert.Equal(t, 5, scan.(*scanner).bp, "TestSkipWhitespaces failed")
}

func TestNextToken(t *testing.T) {
	scan := NewScanner([]byte("  \r\n && 1&2   || 1 "))
	scan.AddSymbol([]byte("&&"), 2)
	scan.AddSymbol([]byte("||"), 3)

	assert.Equal(t, "", string(scan.ScanString()), "TestNextToken failed")

	scan.NextToken()
	assert.Equal(t, 2, scan.Token(), "TestNextToken failed")
	assert.Equal(t, "", string(scan.ScanString()), "TestNextToken failed")

	scan.NextToken()
	assert.Equal(t, 3, scan.Token(), "TestNextToken failed")
	assert.Equal(t, "1&2", string(scan.ScanString()), "TestNextToken failed")

	scan.NextToken()
	assert.Equal(t, TokenEOI, scan.Token(), "TestNextToken failed")
	assert.Equal(t, "1", string(scan.ScanString()), "TestNextToken failed")
}

func TestNextTokenWithGap(t *testing.T) {
	scan := NewScanner([]byte("1 >==2 "))
	scan.AddSymbol([]byte(">"), 2)
	scan.AddSymbol([]byte(">="), 3)
	scan.AddSymbol([]byte("="), 4)

	scan.NextToken()
	assert.Equal(t, 3, scan.Token(), "TestNextToken failed")
	assert.Equal(t, "1", string(scan.ScanString()), "TestNextToken failed")

	scan.NextToken()
	assert.Equal(t, 4, scan.Token(), "TestNextToken failed")
	assert.Equal(t, "", string(scan.ScanString()), "TestNextToken failed")
}
