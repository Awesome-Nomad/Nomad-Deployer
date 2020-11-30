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
	"github.com/spf13/cobra"
	root "github.com/Awesome-Nomad/Nomad-Deployer/cmd"
	"io"
	"os"
)

var outFile string

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate nomad job from input",
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
		fmt.Printf("templateFile: %s, %+v, %s", templateFile, varFiles, workingDir)
		spec, err := nomadDeployer.GenerateJobSpec(templateFile, varFiles, workingDir)
		if err != nil {
			panic(err)
		}
		var writer io.Writer
		switch outFile {
		case "/dev/stdout":
			writer = os.Stdout
		case "/dev/stderr":
			writer = os.Stderr
		default:
			writer, err = os.Create(outFile)
			if err != nil {
				panic(err)
			}
		}
		err = spec.Write(writer)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	localCmd.AddCommand(genCmd)
	genCmd.Flags().StringVarP(&outFile, "out", "o", "/dev/stdout", "Output stream.")
}
