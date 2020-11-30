package auth

import (
	"errors"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var (
	ErrAuthNotFound = errors.New("no auth found with host")
)

type Provider interface {
	Get(host string) (*http.BasicAuth, error)
}

type noOpProvider struct {
}

func (n *noOpProvider) Get(host string) (*http.BasicAuth, error) {
	return nil, nil
}
func NewNoOPProvider() Provider {
	return &noOpProvider{}
}
