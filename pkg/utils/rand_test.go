package utils

import (
	"github.com/bmizerany/assert"
	"testing"
	"time"
)

func TestRandomString(t *testing.T) {
	s := RandomString(32)
	t.Log(s)
	assert.Equal(t, 32, len(s))
}

func TestRandomInt(t *testing.T) {
	t.Log(RandomInt64n(time.Now().Unix()))
}
