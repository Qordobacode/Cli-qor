package file

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_Example(t *testing.T) {
	matched, err := filepath.Match("*.strings", "C:\\data\\go\\src\\github.com\\qordobacode\\cli-v2\\test\\general\\rest.strings")
	assert.Nil(t, err)
	assert.True(t, matched)
}
