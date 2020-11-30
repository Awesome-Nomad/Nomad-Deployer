package config

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	var merror error
	direct, err := os.Open("assets/basic_connection.hcl")
	merror = multierr.Append(merror, err)
	nomadOnly, err := os.Open("assets/nomad_only.hcl")
	merror = multierr.Append(merror, err)
	consulOnly, err := os.Open("assets/consul_only.hcl")
	merror = multierr.Append(merror, err)
	emptyGitRef, err := os.Open("assets/empty_git_ref.hcl")
	merror = multierr.Append(merror, err)
	gitRef, err := os.Open("assets/git_ref.hcl")
	merror = multierr.Append(merror, err)
	sshWithAddress, err := os.Open("assets/ssh_with_address.hcl")
	merror = multierr.Append(merror, err)
	docker, err := os.Open("assets/docker.hcl")
	merror = multierr.Append(merror, err)
	if merror != nil {
		t.Fatal(err)
	}

	type ConfigChecker func(config *Config)

	type Testcase struct {
		name    string
		input   io.Reader
		isError bool
		checker ConfigChecker
	}
	testcases := []Testcase{
		{
			name:    direct.Name(),
			input:   direct,
			isError: false,
			checker: func(result *Config) {
				expected := NewDefaultConfig(&Config{
					TemplateDir:   "./example",
					GitProjectDir: "./projects",
					Environments: []*Environment{
						{
							Name: "basic",
							Nomad: &HashiConnConfig{
								Address: "localhost:4646",
								Token:   "",
								ConnectionConfig: &conn{
									Type:       ConnDirect,
									Connection: &DirectConnection{},
								},
							},
							Consul: &HashiConnConfig{
								Address: "localhost:8500",
								Token:   "",
								ConnectionConfig: &conn{
									Type: ConnSSH,
									Connection: &SSHConnection{
										Config: &SSHConnectionConfig{
											Username: "root",
											KeyFile:  "/root/.ssh/id_rsa",
										},
									},
								},
							},
						},
					},
				})
				result.equals(t, expected)
			},
		},
		{
			name:    nomadOnly.Name(),
			input:   nomadOnly,
			isError: false,
			checker: func(result *Config) {
				expected := NewDefaultConfig(&Config{
					Environments: []*Environment{
						{
							Name: "nomad",
							Nomad: &HashiConnConfig{
								Address: "localhost:4646",
								Token:   "",
								ConnectionConfig: &conn{
									Type:       ConnDirect,
									Connection: &DirectConnection{},
								},
							},
						},
					},
				})
				result.equals(t, expected)
			},
		}, {
			name:    docker.Name(),
			input:   docker,
			isError: false,
			checker: func(result *Config) {
				expected := NewDefaultConfig(&Config{
					Environments: []*Environment{
						{
							Name: "nomad",
							Nomad: &HashiConnConfig{
								Address: "localhost:4646",
								Token:   "",
								ConnectionConfig: &conn{
									Type:       ConnDirect,
									Connection: &DirectConnection{},
								},
							},
							Docker: &DockerConfig{Registry: "registry-1.docker.io"},
						},
					},
				})
				result.equals(t, expected)
			},
		},
		{
			name:    consulOnly.Name(),
			input:   consulOnly,
			isError: false,
			checker: func(result *Config) {
				expected := NewDefaultConfig(&Config{
					Environments: []*Environment{
						{
							Name: "consul",
							Nomad: &HashiConnConfig{
								Address: "localhost:4646",
								Token:   "",
								ConnectionConfig: &conn{
									Type:       ConnDirect,
									Connection: &DirectConnection{},
								},
							},
							Consul: &HashiConnConfig{
								Address: "localhost:8500",
								Token:   "",
								ConnectionConfig: &conn{
									Type: ConnSSH,
									Connection: &SSHConnection{
										Config: &SSHConnectionConfig{
											Username: "root",
											KeyFile:  "/root/.ssh/id_rsa",
										},
									},
								},
							},
						},
					},
				})
				result.equals(t, expected)
			},
		},
		{
			name:    emptyGitRef.Name(),
			input:   emptyGitRef,
			isError: false,
			checker: func(result *Config) {
				expected := NewDefaultConfig(&Config{
					Environments: []*Environment{
						{
							Name: "empty",
							Nomad: &HashiConnConfig{
								Address: "localhost:4646",
								Token:   "",
								ConnectionConfig: &conn{
									Type:       ConnDirect,
									Connection: &DirectConnection{},
								},
							},
						},
					},
				})
				result.equals(t, expected)
			},
		},
		{
			name:    gitRef.Name(),
			input:   gitRef,
			isError: false,
			checker: func(result *Config) {
				expected := NewDefaultConfig(&Config{
					Environments: []*Environment{
						{
							Name: "git_ref",
							Nomad: &HashiConnConfig{
								Address: "localhost:4646",
								Token:   "",
								ConnectionConfig: &conn{
									Type:       ConnDirect,
									Connection: &DirectConnection{},
								},
							},
							Git: &GitConfig{DefaultRef: "refs/remotes/origin/develop"},
						},
					},
				})
				result.equals(t, expected)
			},
		},
		{
			name:    sshWithAddress.Name(),
			input:   sshWithAddress,
			isError: false,
			checker: func(result *Config) {
				expected := NewDefaultConfig(&Config{
					Environments: []*Environment{
						{
							Name: "ssh_with_address",
							Nomad: &HashiConnConfig{
								Address: "localhost:4646",
								Token:   "",
								ConnectionConfig: &conn{
									Type:       ConnDirect,
									Connection: &DirectConnection{},
								},
							},
							Consul: &HashiConnConfig{
								Address: "localhost:8500",
								Token:   "",
								ConnectionConfig: &conn{
									Type: ConnSSH,
									Connection: &SSHConnection{
										Config: &SSHConnectionConfig{
											Username: "root",
											KeyFile:  "/root/.ssh/id_rsa",
											Address:  "localhost:22",
										},
									},
								},
							},
						},
					},
				})
				result.equals(t, expected)
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			hcl, err := ioutil.ReadAll(tc.input)
			if err != nil {
				t.Error(err)
			}
			config, err := LoadConfig(hcl, nil)
			if err != nil {
				if !tc.isError {
					t.Error(err)
				}
			}
			tc.checker(config)
		})
	}

}

func (got *Environment) equals(t *testing.T, want *Environment) {
	// Check environment
	assert.Equal(t, got.Name, want.Name)
	assert.Equal(t, got.ID(), want.ID())
	// Check docker
	assert.Equal(t, got.Docker, want.Docker)
	// Check Nomad
	assert.Equal(t, got.Nomad.Address, want.Nomad.Address)
	assert.Equal(t, got.Nomad.Token, want.Nomad.Token)
	assert.Equal(t, got.Nomad.ConnectionConfig.Type, want.Nomad.ConnectionConfig.Type)
	assert.Equal(t, got.Nomad.ConnectionConfig.Connection, want.Nomad.ConnectionConfig.Connection)

	// Check Consul
	if got.Consul == nil {
		assert.Equal(t, got.Consul, want.Consul)
	} else {
		assert.Equal(t, got.Consul.Address, want.Consul.Address)
		assert.Equal(t, got.Consul.Token, want.Consul.Token)
		assert.Equal(t, got.Consul.ConnectionConfig.Type, want.Consul.ConnectionConfig.Type)
		assert.Equal(t, got.Consul.ConnectionConfig.Connection, want.Consul.ConnectionConfig.Connection)
	}
}

func (got *Config) equals(t *testing.T, want *Config) {
	assert.Equal(t, got.TemplateDir, want.TemplateDir)
	assert.Equal(t, got.GitProjectDir, want.GitProjectDir)
	assert.Equal(t, len(got.Environments), len(want.Environments))
	for idx, got := range got.Environments {
		want := want.Environments[idx]
		got.equals(t, want)
	}
}
