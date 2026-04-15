package activation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeCode(t *testing.T) {
	var codeNeeded bool
	var err error

	_, _, err = DecodeCode("")
	assert.ErrorIs(t, err, ErrInvalidCodeFormat)

	_, _, err = DecodeCode("LPA:")
	assert.ErrorIs(t, err, ErrInvalidCodeFormat)

	_, _, err = DecodeCode("LPA:1")
	assert.ErrorIs(t, err, ErrInvalidCodeFormat)

	info, _, err := DecodeCode("LPA:1$example.com")
	assert.NoError(t, err)
	assert.Equal(t, "example.com", info.SMDP)

	info, _, err = DecodeCode("LPA:1$example.com$matching-id")
	assert.Equal(t, "example.com", info.SMDP)
	assert.Equal(t, "matching-id", info.MatchID)

	info, codeNeeded, err = DecodeCode("LPA:1$example.com$matching-id$$1")
	assert.Equal(t, "example.com", info.SMDP)
	assert.Equal(t, "matching-id", info.MatchID)
	assert.True(t, codeNeeded, "Confirm Code Required Flag")
}

func TestCompleteCode(t *testing.T) {
	const lpaString = "LPA:1$example.com$matching-id"
	assert.Equal(t, lpaString, CompleteCode("LPA:1$example.com$matching-id"))
	assert.Equal(t, lpaString, CompleteCode("1$example.com$matching-id"))
	assert.Equal(t, lpaString, CompleteCode("$example.com$matching-id"))
}
