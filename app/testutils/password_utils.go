package testutils

import "github.com/stretchr/testify/mock"

type MockPasswordUtils struct {
	mock.Mock
}

// GenerateSalt implements utils.IPasswrodUtils.
func (m *MockPasswordUtils) GenerateSalt(length int) (string, error) {
	args := m.Called(length)

	return args.String(0), args.Error(1)
}

// HashPassword implements utils.IPasswrodUtils.
func (m *MockPasswordUtils) HashPassword(password string, salt string) (string, error) {
	args := m.Called(password, salt)
	return args.String(0), args.Error(1)
}

// VerifyPassword implements utils.IPasswrodUtils.
func (m *MockPasswordUtils) VerifyPassword(hashedPassword string, password string, salt string) bool {
	args := m.Called(hashedPassword, password, salt)
	return args.Bool(0)
}
