/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package git

import (
	"bufio"
	"fmt"
	root "github.com/Awesome-Nomad/Nomad-Deployer/cmd"
	"github.com/Awesome-Nomad/Nomad-Deployer/internal/config"
	"github.com/Awesome-Nomad/Nomad-Deployer/internal/git"
	gitauth "github.com/Awesome-Nomad/Nomad-Deployer/internal/git/auth"
	mutils "github.com/Awesome-Nomad/Nomad-Deployer/utils"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/hcl/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"os"
	"path/filepath"
	"strings"
)

var (
	gitUrlStr   string
	netrcFile   string
	alwaysAuth  bool
	gitRef      string
	gitProvider string
	projectsDir string
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "git",
	Short: "Clone Git Repo as Project.",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
		// 1. Checkout project
		environment := viper.GetString("environment")
		cEnv := root.GetEnvironment(environment)
		defaultRef := cEnv.Git.DefaultRef
		log.Printf("env: %s, default_ref: %s", environment, defaultRef)
		checkoutProject(defaultRef, gitProvider)
	},
}

func init() {
	root.RegisterCommand(projectCmd)
	projectCmd.PersistentFlags().StringVarP(&gitUrlStr, "url", "u", "", "Git Url.")
	_ = projectCmd.MarkPersistentFlagRequired("url")
	projectCmd.PersistentFlags().BoolVarP(&alwaysAuth, "always-auth", "", false, "Force always auth to use git.")
	defaultNetRC, _ := homedir.Expand(filepath.Join("~", ".netrc"))
	projectCmd.PersistentFlags().StringVarP(&netrcFile, "netrc", "", defaultNetRC, "netrc file location.")
	projectCmd.PersistentFlags().StringVarP(&gitRef, "ref", "r", "",
		"Git Reference Name. Eg: refs/remotes/origin/master, refs/tags/v1.0.0")
	projectCmd.PersistentFlags().StringVar(&gitProvider, "provider", "gitlab", "Git DefaultProvider. Eg: gitlab, github, etc...")
	defaultProjectDirs, _ := filepath.Abs(config.DefaultProjectsDir)
	projectCmd.PersistentFlags().StringVarP(&projectsDir, "projects-dir", "d", defaultProjectDirs, "Projects download location.")
	// Add function resolver
	root.AddHCLCtxOption(&hclFunction{
		name: "latestTag",
		fn:   createGetLatestTagFunction(),
	})
}

type hclFunction struct {
	name string
	fn   function.Function
}

func (f *hclFunction) Apply(ctx *hcl.EvalContext) {
	ctx.Functions[f.name] = f.fn
}

func createGetLatestTagFunction() function.Function {
	return function.New(&function.Spec{
		Params:   []function.Parameter{},
		VarParam: &function.Parameter{},
		Type:     function.StaticReturnType(cty.String),
		Impl: func(_ []cty.Value, _ cty.Type) (cty.Value, error) {
			prj := createProject(gitProvider)
			latestTag, err := prj.GetLatestSemver()
			return cty.StringVal(plumbing.NewTagReferenceName(latestTag).String()), err
		},
	})
}

func createProject(gitProvider string) *git.Project {
	var authProvider gitauth.Provider
	if _, err := os.Open(netrcFile); os.IsNotExist(err) {
		if alwaysAuth {
			panic("No netrc file found. Always auth must have at least one auth provider")
		}
		authProvider = gitauth.NewNoOPProvider()
	} else {
		authProvider, err = gitauth.NewNetRCProvider(netrcFile)
		if err != nil {
			panic(err)
		}
	}
	project := git.NewProject(projectsDir, authProvider, gitProvider, git.URL(gitUrlStr))
	return project
}

func checkoutProject(defaultRef string, defaultProvider string) (*plumbing.Reference, *git.Project) {
	targetRef := gitRef
	if targetRef == "" {
		targetRef = defaultRef
	}
	gProvider := defaultProvider
	if gProvider == "" {
		gProvider = gitProvider
	}
	log.Printf("Checkout with target ref: %s", targetRef)
	project := createProject(gProvider)
	ref, err := project.Checkout(plumbing.ReferenceName(targetRef))
	if err != nil {
		panic(err)
	}
	log.Printf("Git checkout completed. Project is located at %s", project.GetProjectDir())
	return ref, project
}

func getGitVarFiles(reference *plumbing.Reference, dockerImage string, workingDir string) (vars []string, err error) {
	vars, err = root.GetVarFiles(workingDir)
	var (
		imageFile *os.File
		dockerTag string
	)
	refName := reference.Name()
	if refName.IsTag() {
		dockerTag = refName.Short()
	} else {
		dockerTag = reference.Hash().String()
	}
	imageFile, err = mutils.NewTempFile("image.*.yaml")
	root.CheckIfErr(err)
	defer imageFile.Close()
	w := bufio.NewWriter(imageFile)
	_, err = w.WriteString(fmt.Sprintf("image: %s:%s", dockerImage, dockerTag))
	defer w.Flush()
	vars = append(vars, imageFile.Name())
	return

}

func getDockerImage(dockerConfig *config.DockerConfig, prj *git.Project) string {
	dockerRegistry := strings.TrimSuffix(dockerConfig.Registry, "/")
	return fmt.Sprintf("%s/%s", dockerRegistry, prj.Slug())
}
