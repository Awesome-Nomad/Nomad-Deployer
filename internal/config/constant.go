package config

import "errors"

var (
	ErrNoConnection = errors.New("no connection defined")
)

type ConnectionType string

const (
	ConnDirect ConnectionType = "direct"
	ConnSSH    ConnectionType = "ssh"
)

const (
	DeploymentDir = "deployment"
	SSHPort       = 22
)

const (
	DefaultProjectsDir = "./projects"
	DefaultTemplateDir = "./deploy-template/nomad/deployment"
)
