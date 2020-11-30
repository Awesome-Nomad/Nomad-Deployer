package auth

import (
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNetRCProvider(t *testing.T) {
	provider, _ := NewNetRCProvider("assets/netrc")
	a, err := provider.Get("github.com")
	assert.Nil(t, err)
	assert.Equal(t, &http.BasicAuth{
		Username: "liemdeptrai",
		Password: "holymollydolly",
	}, a)
	_, err = provider.Get("gitlab.com")
	assert.Equal(t, ErrAuthNotFound, err)
}
