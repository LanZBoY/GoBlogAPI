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

func TestGenerateSaltUnique(t *testing.T) {
	pu := &PasswordUtils{}
	s1, err1 := pu.GenerateSalt(4)
	s2, err2 := pu.GenerateSalt(4)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Len(t, s1, 8)
	assert.Len(t, s2, 8)
	assert.NotEqual(t, s1, s2)
}

func TestHashAndVerifyPassword(t *testing.T) {
	pu := &PasswordUtils{}
	salt, _ := pu.GenerateSalt(4)
	pwd := "secret"
	hashed, err := pu.HashPassword(pwd, salt)
	assert.NoError(t, err)
	assert.NotEqual(t, pwd, hashed)
	assert.True(t, pu.VerifyPassword(hashed, pwd, salt))
	assert.False(t, pu.VerifyPassword(hashed, "wrong", salt))
	assert.False(t, pu.VerifyPassword(hashed, pwd, salt+"1"))
}
