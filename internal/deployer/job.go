package deployer

import (
	"github.com/imdario/mergo"
	"strings"
)

type Job interface {
	GetVersion() string
	GetSpec() Spec
}

type Config struct {
	ServiceName string           `yaml:"service_name"`
	Image       string           `yaml:"image"`
	Count       uint             `yaml:"count"`
	Canary      uint             `yaml:"count"`
	AutoRevert  bool             `yaml:"auto_revert"`
	AutoPromote bool             `yaml:"auto_promote"`
	MaxParallel uint             `yaml:"max_parallel"`
	AppFiles    []AppFile        `yaml:"app_files"`
	Entrypoint  []string         `yaml:"entrypoint"`
	Directory   []Directory      `yaml:"directories"`
	Services    []ConsulService  `yaml:"services"`
	Resources   AppResource      `yaml:"resources"`
	Constraints []MetaConstraint `yaml:"constraints"`
	ExtraHosts  []string         `yaml:"extra_hosts"`
}

func (s *Config) GetSpec() Spec {
	return Spec{}
}

func (s *Config) GetVersion() string {
	return "0"
}

func (s *Config) Merge(other Config) error {
	return mergo.Merge(s, other, mergo.WithAppendSlice, mergo.WithOverride)
}

type AppResource struct {
	MBits  uint `yaml:"mbits"`
	CPU    uint `yaml:"cpu"`
	Memory uint `yaml:"memory"`
}

type ConsulService struct {
	Name           string   `yaml:"name"`
	Port           string   `yaml:"port"`
	Prometheus     bool     `yaml:"prometheus"`
	PrometheusPort string   `yaml:"prometheus_port"`
	Tags           []string `yaml:"tags"`
}

type AppFile struct {
	SourcePath  string `yaml:"src"`
	Destination string `yaml:"destination"`
	Environment bool   `yaml:"env"`
}

type MetaConstraint string

func (m MetaConstraint) GetMetaKey() string {
	ms := string(m)
	l := strings.Split(ms, "=")
	if len(l) > 1 {
		return l[0]
	}
	return ms
}

func (m MetaConstraint) GetValue() string {
	ms := string(m)
	l := strings.Split(ms, "=")
	if len(l) > 1 {
		return l[1]
	}
	return "True"
}

type Directory string

func (d Directory) GetBindPath() string {
	da := strings.Split(string(d), ":")
	if len(da) > 1 {
		return da[1]
	}
	return string(d)
}
func (d Directory) GetHostPath() string {
	da := strings.Split(string(d), ":")
	return da[0]
}

func NewSpec() *Config {
	return &Config{
		AutoRevert:  false,
		AutoPromote: false,
		MaxParallel: 1,
		AppFiles:    nil,
		Entrypoint:  nil,
		Services:    nil,
		Resources: AppResource{
			MBits:  1,
			CPU:    100,
			Memory: 300,
		},
		Constraints: nil,
		ExtraHosts:  nil,
	}
}
