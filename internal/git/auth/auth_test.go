package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoOpProvider_Get(t *testing.T) {
	p := NewNoOPProvider()
	auth, err := p.Get("any_host")
	assert.Nil(t, auth)
	assert.Nil(t, err)
}
