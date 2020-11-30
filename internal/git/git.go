package git

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/rs/zerolog/log"
	config2 "github.com/Awesome-Nomad/Nomad-Deployer/internal/config"
	gitauth "github.com/Awesome-Nomad/Nomad-Deployer/internal/git/auth"
	"github.com/Awesome-Nomad/Nomad-Deployer/pkg/semver"
	"net/url"
	"path/filepath"
	"strings"
)

var (
	ErrNoTagFound = errors.New("no tag found in project")
)

type Project struct {
	gitURL       *url.URL
	projectsDir  string
	authProvider gitauth.Provider
	authOptional bool
	provider     config2.GitProvider
}

func (p *Project) GetURL() url.URL {
	return *p.gitURL
}

// Slug
func (p *Project) Slug() string {
	gitPath := strings.ReplaceAll(p.gitURL.RequestURI(), filepath.Ext(p.gitURL.RequestURI()), "")
	slug := strings.ReplaceAll(strings.TrimPrefix(gitPath, "/"), "/", "-")
	return p.provider.NormalizeSlug(slug)
}

func (p *Project) GetLatestSemver() (string, error) {
	r := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		URLs: []string{p.gitURL.String()},
	})

	auth, err := p.authProvider.Get(p.gitURL.Hostname())
	if err != nil && !p.authOptional {
		return "", err
	}

	refs, err := r.List(&git.ListOptions{Auth: auth})
	if err != nil {
		return "", nil
	}

	tagsSemver := semver.NewTagsSemver()
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tagRef := ref.Name()
			tagAnnotation := getTagAnnotated(tagRef.String())
			tagsSemver.AddTag(tagAnnotation)
		}
	}
	latestVersion := tagsSemver.GetLatestTag()
	return latestVersion, nil
}

func (p *Project) GetProjectDir() string {
	return GetGitProjectDir(p.projectsDir, p.gitURL)
}

func (p *Project) String() string {
	return fmt.Sprintf("projectDir: %s", p.GetProjectDir())
}

func (p *Project) Checkout(ref plumbing.ReferenceName) (reference *plumbing.Reference, err error) {
	log.Info().Msgf("Checkout project %s with ref: %s", p.gitURL.String(), ref)
	var (
		w *git.Worktree
	)
	auth, err := p.authProvider.Get(p.gitURL.Hostname())
	if err != nil && !p.authOptional {
		return
	}
	gitProjectPath := p.GetProjectDir()
	repo, err := git.PlainOpen(gitProjectPath)
	if err != nil {
		switch err {
		case git.ErrRepositoryNotExists:
			// Do clone
			repo, err = git.PlainClone(gitProjectPath, false, &git.CloneOptions{
				URL:               p.gitURL.String(),
				Auth:              auth,
				Depth:             MaxDepth,
				RecurseSubmodules: 0,
			})
			if err != nil {
				return
			}
		default:
			return
		}
	} else {
		// Work with already exist project
		w, err = repo.Worktree()
		if err != nil {
			return
		}
		// 1. Clean project
		err = w.Clean(&git.CleanOptions{Dir: true})
		if err != nil {
			return
		}
		// 2. Pull source with ref
		err = w.Pull(&git.PullOptions{
			RemoteName: "origin",
			Depth:      MaxDepth,
			Auth:       auth,
			Progress:   nil,
			Force:      true,
		})
		if err != nil {
			switch err {
			case git.NoErrAlreadyUpToDate:
			case git.ErrNonFastForwardUpdate:
			default:
				return
			}
		}
	}
	// Work with already exist project
	w, err = repo.Worktree()
	if err != nil {
		return
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: ref,
		Create: false,
		Force:  true,
		Keep:   false,
	})
	if err != nil {
		return nil, err
	}
	reference, err = repo.Reference(ref, true)
	return
}

func NewProject(projectsDir string, authProvider gitauth.Provider, gitProvider string, options ...Option) (r *Project) {
	r = &Project{
		projectsDir:  projectsDir,
		authProvider: authProvider,
		provider:     config2.GitProvider(gitProvider),
	}
	for _, opt := range options {
		err := opt.apply(r)
		if err != nil {
			panic(err)
		}
	}
	return
}

func getTagAnnotated(tagRef string) string {
	return strings.ReplaceAll(tagRef, "refs/tags/", "")
}
