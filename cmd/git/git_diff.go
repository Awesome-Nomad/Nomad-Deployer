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
	"github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	root "github.com/Awesome-Nomad/Nomad-Deployer/cmd"
	"github.com/Awesome-Nomad/Nomad-Deployer/cmd/local"
	"github.com/Awesome-Nomad/Nomad-Deployer/internal/cli_ui"
	pdeployer "github.com/Awesome-Nomad/Nomad-Deployer/internal/deployer"
)

// projectDiffCmd represents the projectDiff command
var projectDiffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Display diff between current and expected",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
		env := root.GetInputEnv()
		deployerCfg := root.GetDeployerCfg()
		templateFile, err := deployerCfg.GetTemplateFile()
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

		diffResponse := meta[pdeployer.RawResponse].(*api.JobPlanResponse)
		diffContent, err := cli_ui.NomadFormatDryRun(diffResponse, jobSpec)
		if err != nil {
			panic(err)
		}
		local.BasicUIOutput(diffContent)
	},
}

func init() {
	projectCmd.AddCommand(projectDiffCmd)
}
