package apperror

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"wentee/blog/app/schema/apperror/errcode"
)

func TestNewAndMethods(t *testing.T) {
	raw := errors.New("raw")
	ae := New(http.StatusBadRequest, errcode.BAD_REQUEST, raw)

	assert.Equal(t, http.StatusBadRequest, ae.Status)
	assert.Equal(t, errcode.BAD_REQUEST, ae.Code)
	assert.Equal(t, raw, ae.RawErr())
	// Message field should be empty by default
	assert.Empty(t, ae.Message)

	// GetMessage should return the mapping from errcode
	assert.Equal(t, errcode.Message(errcode.BAD_REQUEST), ae.GetMessage())
}
