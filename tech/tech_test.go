package tech

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Tech(t *testing.T) {
	tm, err := String2Time("2022-04-22 0:50:35")
	assert.Nil(t, err)
	assert.Equal(t, int64(1650588635), tm.Unix())

	tm, err = String2Time("2022-04-22 0:50:")
	assert.NotNil(t, err)

	tm, err = IntString2Time("1650588635")
	assert.Nil(t, err)
	assert.Equal(t, int64(1650588635), tm.Unix())

	ts := IntString2yymmdd_hhmmss("1650588635")
	assert.Nil(t, err)
	assert.Equal(t, "2022-04-22 00:50:35", ts)

	LogInfo("")
	LogWarn("")

}
