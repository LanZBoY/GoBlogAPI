package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordUtils(t *testing.T) {
	pu := &PasswordUtils{}
	salt, err := pu.GenerateSalt(4)
	assert.NoError(t, err)
	assert.Len(t, salt, 8)

	pwd := "secret"
	hashed, err := pu.HashPassword(pwd, salt)
	assert.NoError(t, err)
	assert.True(t, pu.VerifyPassword(hashed, pwd, salt))
	assert.False(t, pu.VerifyPassword(hashed, "bad", salt))
}
