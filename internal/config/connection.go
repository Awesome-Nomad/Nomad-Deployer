package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/rs/zerolog/log"
	"sync"
)

type Connection interface {
	Init(address string) error
	GetAddress() (string, error)
	Destroy() error
}
type SSHConnectionConfig struct {
	Username string `hcl:"username"`
	KeyFile  string `hcl:"key_file"`
	Address  string `hcl:"address,optional"`
}

type sshConnectionWrapper struct {
	Config *SSHConnectionConfig `hcl:"ssh,block"`
}

type conn struct {
	Type ConnectionType `hcl:"type"`
	HCL  hcl.Body       `hcl:",remain" json:"-"`
	Connection
}
type HashiConnConfig struct {
	Address          string `hcl:"address"`
	Token            string `hcl:"acl_token"`
	ConnectionConfig *conn  `hcl:"connection,block"`
	mu               sync.RWMutex
	initialized      bool
}

func (c *HashiConnConfig) Connect() (string, error) {
	if c.ConnectionConfig == nil {
		return "", ErrNoConnection
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.initialized {
		err := c.ConnectionConfig.Init(c.Address)
		if err != nil {
			log.Printf("Create connection problem. %+v\n", err)
			return "", err
		}
	}
	return c.ConnectionConfig.GetAddress()
}

func (c *HashiConnConfig) Destroy() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.initialized = false
	return c.ConnectionConfig.Destroy()
}
