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
package local

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	root "github.com/Awesome-Nomad/Nomad-Deployer/cmd"
	"github.com/Awesome-Nomad/Nomad-Deployer/pkg/utils"
)

var (
	workingDir   string
	templateFile string
)

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Run Project as Local",
	Long:  ``,
	PreRun: func(_ *cobra.Command, _ []string) {
		commandInitialized()
	},
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.UseLine()
	},
}

func commandInitialized() {
	var err error
	templateFile, err = root.GetDeployerCfg().GetTemplateFile()
	root.CheckIfErr(err)
	workingDir, err = utils.ExpandDir(workingDir)
	root.CheckIfErr(err)
}

func init() {
	root.RegisterCommand(localCmd)
	localCmd.PersistentFlags().StringVarP(&workingDir, "project-directory", "w", "", "Project directory. Required")
	err := localCmd.MarkPersistentFlagRequired("project-directory")
	log.Err(err)
}
