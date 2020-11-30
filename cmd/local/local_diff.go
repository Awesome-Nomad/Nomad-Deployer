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
	"github.com/hashicorp/nomad/api"
	"github.com/mattn/go-colorable"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/colorstring"
	"github.com/spf13/cobra"
	root "github.com/Awesome-Nomad/Nomad-Deployer/cmd"
	"github.com/Awesome-Nomad/Nomad-Deployer/internal/cli_ui"
	"github.com/Awesome-Nomad/Nomad-Deployer/internal/deployer"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Generate difference between current state and new state",
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
		diffResp := meta[deployer.RawResponse].(*api.JobPlanResponse)
		diffContent, _ := cli_ui.NomadFormatDryRun(diffResp, spec)
		BasicUIOutput(diffContent)
	},
}

func BasicUIOutput(v string) {
	outputUI := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      colorable.NewColorableStdout(),
		ErrorWriter: colorable.NewColorableStderr(),
	}
	outputUI.Output(Colorize().Color(v))
}

func Colorize() *colorstring.Colorize {
	return &colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Disable: !terminal.IsTerminal(int(os.Stdout.Fd())),
		Reset:   true,
	}
}

func init() {
	localCmd.AddCommand(diffCmd)
}
