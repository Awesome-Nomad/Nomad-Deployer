// +build local

package git

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/Awesome-Nomad/Nomad-Deployer/internal/git/auth"
	"os"
	"path/filepath"
	"testing"
)

func TestNewProject(t *testing.T) {
	userhomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	netrcPath := filepath.Join(userhomeDir, ".netrc")
	authProvider, err := auth.NewNetRCProvider(netrcPath)
	if err != nil {
		panic(err)
	}
	prj := NewProject("./projects",
		authProvider,
		"",
		URL("https://github.com/hashicorp/nomad.git"),
	)

	fmt.Printf("Path: %s\n", prj.Slug())

	checkout(prj, plumbing.NewTagReferenceName("v1.6.0"))
	checkout(prj, "refs/remotes/origin/develop")
	checkout(prj, plumbing.NewRemoteReferenceName("origin", "master"))
	checkout(prj, plumbing.NewRemoteReferenceName("origin", "develop"))

}
func checkout(repo *Project, ref plumbing.ReferenceName) {
	commitHash, err := repo.Checkout(ref)
	fmt.Printf("ref: %s, hash: %s, name: %+v\n", ref, commitHash, commitHash.Name())
	if err != nil {
		panic(err)
	}
}

func TestAbc(t *testing.T) {

	p := NewProject("./projects",
		auth.NewNoOPProvider(),
		"",
		URL("https://github.com/hashicorp/nomad.git"))

	latestVersion, err := p.GetLatestSemver()

	t.Logf("latestVersion: %s. err: %+v", latestVersion, err)
}
