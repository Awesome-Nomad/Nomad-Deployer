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
	"fmt"
	"github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog/log"
	root "github.com/Awesome-Nomad/Nomad-Deployer/cmd"
	"github.com/Awesome-Nomad/Nomad-Deployer/cmd/local"
	"github.com/Awesome-Nomad/Nomad-Deployer/internal/cli_ui"
	pdeployer "github.com/Awesome-Nomad/Nomad-Deployer/internal/deployer"

	"github.com/spf13/cobra"
)

// projectDeployCmd represents the projectDeploy command
var projectDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy project",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
		env := root.GetInputEnv()
		cfg := root.GetDeployerCfg()
		templateFile, err := cfg.GetTemplateFile()
		root.CheckIfErr(err)
		// 1. Checkout project
		cEnv := root.GetEnvironment(env)
		defaultRef := cEnv.Git.DefaultRef
		defaultProvider := cEnv.Git.DefaultProvider
		log.Printf("env: %s, default_ref: %s", env, defaultRef)
		ref, project := checkoutProject(defaultRef, string(defaultProvider))
		deployer, err := root.CreateDeployer(env)
		root.CheckIfErr(err)
		// 2. Generate job file.
		workingDir := project.GetProjectDir()
		varFiles, err := getGitVarFiles(ref, getDockerImage(cEnv.Docker, project), workingDir)
		root.CheckIfErr(err)
		log.Printf("varFiles: %s", varFiles)
		jobSpec, err := deployer.GenerateJobSpec(templateFile, varFiles, workingDir)
		root.CheckIfErr(err)
		_, meta, err := deployer.Diff(jobSpec)
		root.CheckIfErr(err)
		deploySpec, err := deployer.Deploy(jobSpec, meta)
		root.CheckIfErr(err)
		diffResponse := meta[pdeployer.RawResponse].(*api.JobPlanResponse)
		diffContent, err := cli_ui.NomadFormatDryRun(diffResponse, jobSpec)
		root.CheckIfErr(err)
		local.BasicUIOutput(diffContent)
		fmt.Println(deploySpec.Pretty())
	},
}

func init() {
	projectCmd.AddCommand(projectDeployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectDeployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectDeployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
