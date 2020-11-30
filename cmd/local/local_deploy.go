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
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	root "github.com/Awesome-Nomad/Nomad-Deployer/cmd"
)

// applyCmd represents the Apply command
var applyCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy application",
	Long:  ``,
	PreRun: func(_ *cobra.Command, _ []string) {
		commandInitialized()
	},
	Run: func(_ *cobra.Command, _ []string) {
		// Create deployer
		nomadDeployer, err := root.CreateDeployer(root.GetInputEnv())
		if err != nil {
			panic(err)
		}
		// Generate doc
		varFiles, err := root.GetVarFiles(workingDir)
		if err != nil {
			panic(err)
		}
		spec, err := nomadDeployer.GenerateJobSpec(templateFile, varFiles, workingDir)
		if err != nil {
			panic(err)
		}
		_, meta, err := nomadDeployer.Diff(spec)
		if err != nil {
			panic(err)
		}
		spec, err = nomadDeployer.Deploy(spec, meta)
		if err != nil {
			panic(err)
		}
		fmt.Println(spec.Content)
	},
}

func init() {
	localCmd.AddCommand(applyCmd)
	err := localCmd.MarkPersistentFlagRequired("project-directory")
	log.Err(err)
}
