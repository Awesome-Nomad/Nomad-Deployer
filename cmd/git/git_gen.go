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
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	root "github.com/Awesome-Nomad/Nomad-Deployer/cmd"
	"github.com/Awesome-Nomad/Nomad-Deployer/pkg/utils"
	mutils "github.com/Awesome-Nomad/Nomad-Deployer/utils"
	"io/ioutil"
)

// projectGenCmd represents the projectGen command
var projectGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate job file for project",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
		env := viper.GetString("environment")
		deployerCfg := root.GetDeployerCfg()
		templateFile, err := deployerCfg.GetTemplateFile()
		root.CheckIfErr(err)
		// 1. Checkout project
		cEnv := root.GetEnvironment(env)
		defaultRef := cEnv.Git.DefaultRef
		defaultProvider := cEnv.Git.DefaultProvider
		log.Printf("env: %s, default_ref: %s", env, defaultRef)
		ref, project := checkoutProject(defaultRef, string(defaultProvider))
		// 2. Generate job file.
		deployer, err := root.CreateDeployer(env)
		root.CheckIfErr(err)
		workingDir, err := utils.ExpandDir(project.GetProjectDir())
		root.CheckIfErr(err)
		varFiles, err := getGitVarFiles(ref, getDockerImage(cEnv.Docker, project), workingDir)
		log.Printf("varFiles: %s", varFiles)
		root.CheckIfErr(err)
		jobSpec, err := deployer.GenerateJobSpec(templateFile, varFiles, workingDir)
		root.CheckIfErr(err)
		jobFile, err := mutils.NewTempFile("job.*.nomad")
		root.CheckIfErr(err)
		err = ioutil.WriteFile(jobFile.Name(), []byte(jobSpec.Pretty()), 0644)
		root.CheckIfErr(err)
		log.Info().Msgf("Job file generated and placed at %s", jobFile.Name())
	},
}

func init() {
	projectCmd.AddCommand(projectGenCmd)
}
