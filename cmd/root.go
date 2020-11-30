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
package cmd

import (
	"errors"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	nomadapi "github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/Awesome-Nomad/Nomad-Deployer/internal/config"
	"github.com/Awesome-Nomad/Nomad-Deployer/internal/deployer"
	"github.com/Awesome-Nomad/Nomad-Deployer/internal/deployer/nomad"
	. "github.com/Awesome-Nomad/Nomad-Deployer/pkg/utils"
	"io/ioutil"
	"os"
)

var (
	deployerCfgFile string
	environment     string
	deployerCfg     *config.Config
	hclContext      = &hcl.EvalContext{
		Variables: make(map[string]cty.Value),
		Functions: make(map[string]function.Function),
	}
	hclCtxOptions []HCLCtxOption
	Version       = "undefined"
)

var (
	ErrEnvironmentNotFound = errors.New("environment not found")
)

type HCLCtxOption interface {
	Apply(ctx *hcl.EvalContext)
}

func AddHCLCtxOption(option HCLCtxOption) {
	hclCtxOptions = append(hclCtxOptions, option)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "deployer",
	Short: "",
	Long:  ``,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		log.Debug().Msg("Root")
		if viper.GetBool("version") {
			return
		}
		// Setup config
		hclBytes, err := ioutil.ReadFile(deployerCfgFile)
		if err != nil {
			panic(err)
		}
		// Parse config
		for _, opt := range hclCtxOptions {
			opt.Apply(hclContext)
		}
		deployerCfg, err = config.LoadConfig(hclBytes, hclContext)
		if err != nil {
			panic(err)
		}
		// Printout current config
		log.Debug().Msgf("deployerCfgFile: %+v", deployerCfgFile)
		log.Printf("environment: %+v", environment)
	},
	Run: func(_ *cobra.Command, _ []string) {
		if viper.GetBool("version") {
			fmt.Println(Version)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&deployerCfgFile, "config", "c", "deployer.hcl", "config file")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Show verbose log")
	rootCmd.Flags().Bool("version", false, "Show verbose log")
	rootCmd.PersistentFlags().StringVarP(&environment, "environment", "e", "dev",
		"Environment. Eg: dev, prod, stag, etc...")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
	err := viper.BindPFlags(rootCmd.Flags())
	log.Err(err)
	verboseLog := viper.GetBool("verbose")
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if verboseLog {
		log.Logger = log.With().Caller().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}

func GetVarFiles(workingDir string) (vars []string, err error) {
	workingDir, err = ExpandDir(workingDir)
	if err != nil {
		return
	}
	baseVarFiles := deployerCfg.GetVarFilesForEnv(environment)
	// 3.2 Init project var files
	projectVarFiles := deployerCfg.GetVarFilesForEnvWithDir(environment, workingDir)
	// 3.3 Append all var files into one
	vars = append(baseVarFiles, projectVarFiles...)
	return
}

func GetEnvironment(env string) *config.Environment {
	if "" == env {
		env = GetInputEnv()
	}
	selectEnvironment := deployerCfg.GetEnvironment(env)
	if selectEnvironment == nil {
		panic(ErrEnvironmentNotFound)
	}
	return selectEnvironment
}

func GetInputEnv() string {
	return environment
}

func CreateDeployer(env string) (deployer.Deployer, error) {
	selectEnvironment := GetEnvironment(env)
	var (
		consulAddr string
		nomadAddr  string
		err        error
	)
	// Init Nomad deployer
	nomadConfig := nomadapi.DefaultConfig()

	nomadAddr, err = selectEnvironment.Nomad.Connect()
	if err != nil {
		log.Error().Msgf("Failed to init Nomad connection. %+v", err)
		return nil, err
	}

	nomadConfig.Address = nomadAddr
	nomadConfig.SecretID = selectEnvironment.Nomad.Token
	if consul := selectEnvironment.Consul; consul != nil {
		consulAddr, err = selectEnvironment.Consul.Connect()
		if err != nil {
			log.Warn().Msgf("Failed to init consul connection. %+v", err)
		}
	}
	nomadDeployer, err := nomad.NewNomadDeployer(nomadConfig, consulAddr)
	return nomadDeployer, err
}

func CheckIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
func GetDeployerCfg() *config.Config {
	return deployerCfg
}

func RegisterCommand(c *cobra.Command) {
	rootCmd.AddCommand(c)
}
