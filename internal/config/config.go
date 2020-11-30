package config

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/rs/zerolog/log"
	"path/filepath"
)

type Destroyable interface {
	Destroy() error
}

type Connectable interface {
	Connect() (string, error)
}

type EnvironmentID string

type Environment struct {
	Name   string           `hcl:"name,label"`
	Nomad  *HashiConnConfig `hcl:"nomad,block"`
	Consul *HashiConnConfig `hcl:"consul,block"`
	Git    *GitConfig       `hcl:"git,block"`
	Docker *DockerConfig    `hcl:"docker,block"`
}

func (e *Environment) ID() string {
	return e.Name
}

type Config struct {
	TemplateDir   string         `hcl:"template_dir,optional"`
	GitProjectDir string         `hcl:"git_project_dir,optional"`
	Job           *JobConfig     `hcl:"job,block"`
	Environments  []*Environment `hcl:"env,block"`
	envMap        map[string]*Environment
}

func NewDefaultConfig(base *Config) (c *Config) {
	var (
		gitProjectDir = DefaultProjectsDir
		templateDir   = DefaultTemplateDir
		envMap        = make(map[string]*Environment)
	)

	c = &Config{}
	if base != nil {
		if gitProjectDir = base.GitProjectDir; gitProjectDir == "" {
			gitProjectDir = DefaultProjectsDir
		}
		if templateDir = base.TemplateDir; templateDir == "" {
			templateDir = DefaultTemplateDir
		}
		if base.envMap != nil {
			envMap = base.envMap
		}
		c.Job = base.Job
		c.Environments = base.Environments
	}
	c.TemplateDir = templateDir
	c.GitProjectDir = gitProjectDir
	c.envMap = envMap
	return
}

func (c *Config) GetVarFilesForEnv(envId string) []string {
	return c.GetVarFilesForEnvWithDir(envId, c.TemplateDir)
}

func (c *Config) GetVarFilesForEnvWithDir(envId string, baseDir string) []string {
	var varFiles []string
	env := c.GetEnvironment(envId)
	if env == nil {
		return varFiles
	}
	var baseDirAbs = baseDir
	if !filepath.IsAbs(baseDir) {
		baseDirAbs, _ = filepath.Abs(baseDir)
	}
	// By default, we will try to import base.yaml and ${env_name}.yaml if exists
	for _, f := range []string{"base.yaml", env.Name + ".yaml"} {
		fileAbsPath := filepath.Join(baseDirAbs, DeploymentDir, f)
		if FileExists(fileAbsPath) {
			varFiles = append(varFiles, fileAbsPath)
		} else {
			log.Debug().Msgf("file %s not found", fileAbsPath)
		}
	}
	return varFiles
}

type JobConfig struct {
	TemplateFile string `hcl:"template,optional"`
}

func (c *Config) GetEnvironment(envKey string) *Environment {
	return c.envMap[envKey]
}

func (c *Config) GetTemplateFile() (string, error) {
	var templateFile string
	if c.Job == nil || c.Job.TemplateFile == "" {
		templateFile = "job.nomadtpl"
	} else {
		templateFile = c.Job.TemplateFile
	}
	return filepath.Abs(filepath.Join(c.TemplateDir, DeploymentDir, templateFile))
}

func (c *Config) Destroy() error {
	for _, e := range c.Environments {
		err := e.Destroy()
		if err != nil {
			return err
		}
	}
	return nil
}
func (e *Environment) Destroy() error {
	err := e.Nomad.Destroy()
	if err != nil {
		return err
	}
	err = e.Consul.Destroy()
	if err != nil {
		return err
	}
	return nil
}

func LoadConfig(hclBytes []byte, evalCtx *hcl.EvalContext) (*Config, error) {
	parser := hclparse.NewParser()
	srcHCL, diag := parser.ParseHCL(hclBytes, "config.hcl")
	if diag != nil && diag.HasErrors() {
		return nil, fmt.Errorf("failed to parse config file. %w", diag)
	}
	config := NewDefaultConfig(nil)
	if diag := gohcl.DecodeBody(srcHCL.Body, evalCtx, config); diag.HasErrors() {
		return nil, fmt.Errorf("failed to parse config file. %w", diag)
	}
	var validateErr error

CONFIG:
	for _, env := range config.Environments {
		_, ok := config.envMap[env.ID()]
		if ok {
			validateErr = fmt.Errorf("duplication on environment %s", env.ID())
			break CONFIG
		}
		config.envMap[env.ID()] = env
		// Nomad
		if nomad := env.Nomad; nomad != nil {
			if conn := nomad.ConnectionConfig; conn != nil {
				nomad.ConnectionConfig.Connection, validateErr = createConnection(conn.Type, conn.HCL)
				if validateErr != nil {
					break CONFIG
				}
			}
		} else {
			cn, _ := createConnection(ConnDirect, nil)
			nomadConnectionCfg := &conn{
				Type:       ConnDirect,
				HCL:        nil,
				Connection: cn,
			}
			env.Nomad = &HashiConnConfig{
				Address:          "localhost:4646",
				Token:            "",
				ConnectionConfig: nomadConnectionCfg,
			}
		}
		// Consul
		if consul := env.Consul; consul != nil {
			if conn := consul.ConnectionConfig; conn != nil {
				consul.ConnectionConfig.Connection, validateErr = createConnection(conn.Type, conn.HCL)
				if validateErr != nil {
					break CONFIG
				}
			}
		}
	}
	if validateErr != nil {
		return nil, validateErr
	}
	return config, nil
}
func createConnection(connType ConnectionType, hclBody hcl.Body) (connection Connection, err error) {
	switch connType {
	case ConnDirect:
		connection = &DirectConnection{}
	case ConnSSH:
		sshConnWrapper := &sshConnectionWrapper{}
		if diag := gohcl.DecodeBody(hclBody, nil, sshConnWrapper); diag.HasErrors() {
			err = diag
		}
		connection = &SSHConnection{Config: sshConnWrapper.Config}
	default:
		err = fmt.Errorf("invalid connection type. %s", connType)
	}
	return connection, err
}
