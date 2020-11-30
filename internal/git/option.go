package git

import (
	"net/url"
	"strings"
)

const gitRepoSuffix = ".git"

type Option interface {
	apply(r *Project) error
}

type URL string
type IsAuthOptional bool

func (u IsAuthOptional) apply(r *Project) error {
	r.authOptional = bool(u)
	return nil
}

func (u URL) apply(r *Project) (err error) {
	r.gitURL, err = url.Parse(ensureGitRepoURL(string(u)))
	if "http" == r.gitURL.Scheme {
		// Always secure
		r.gitURL.Scheme = "https"
	}
	return
}

func ensureGitRepoURL(url string) string {
	if !strings.HasSuffix(url, gitRepoSuffix) {
		url = url + gitRepoSuffix
	}
	return url
}
