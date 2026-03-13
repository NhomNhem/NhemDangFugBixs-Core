package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMath_Add(t *testing.T) {
	result := 1 + 1
	assert.Equal(t, 2, result, "1 + 1 should equal 2")
}
