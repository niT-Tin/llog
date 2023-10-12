package llog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithColor(t *testing.T) {
	assert.Equal(t, "[31;1mhello[0m", WithColor("hello", FgRed))
	assert.Equal(t, "[32;1mhello[0m", WithColor("hello", FgGreen))
	assert.Equal(t, "[34;1mhello[0m", WithColor("hello", FgBlue))
	assert.Equal(t, "[44;97;1mhello[0m", WithColor("hello", BgBlue))
}
