package v1

import (
	"github.com/Masterminds/semver"
	"github.com/Awesome-Nomad/Nomad-Deployer/internal/deployer"
)

var maxSupportedVersion, _ = semver.NewVersion("1")

type Spec struct {
	Version string           `yaml:"version"`
	JobSpec *deployer.Config `yaml:",inline"`
	Steps   []Step           `yaml:"steps"`
}

func (s *Spec) GetVersion() string {
	return s.Version
}

type Step struct {
	Name    string          `yaml:"name"`
	JobSpec deployer.Config `yaml:",inline"`
}

func (s *Spec) IsValid() bool {
	v, err := semver.NewVersion(s.Version)
	if err != nil {
		panic("Invalid version. " + s.Version)
	}
	return maxSupportedVersion.Compare(v) >= 0
}

func NewSpec() *Spec {
	return &Spec{
		Version: maxSupportedVersion.String(),
		JobSpec: deployer.NewSpec(),
		Steps:   nil,
	}
}
