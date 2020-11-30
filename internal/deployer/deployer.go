package deployer

import (
	"github.com/Awesome-Nomad/Nomad-Deployer/internal/beauter"
	"io"
)

type Spec struct {
	Content     string
	ContentType beauter.ContentType
}

func (s *Spec) Write(writer io.Writer) error {
	_, err := writer.Write([]byte(s.Content))
	return err
}

func (s *Spec) Pretty() string {
	v, _ := s.ContentType.Format(s.Content)
	return v
}

type Deployer interface {
	// GenerateJobSpec will merge all job into one and return
	// Nomad: HCL/JSON
	// K8S	: YAML/JSON
	GenerateJobSpec(jobTemplateFile string, jobConfigFiles []string, workingDir string) (jobSpec *Spec, err error)
	// Diff show what changes will be made
	// The diff between current state and new desired state
	// Nomad: nomad plan
	// K8S	: kubectl diff
	Diff(jobSpec *Spec) (diff *Spec, meta map[string]interface{}, err error)
	// Deploy Job. Extends var may be good for some system. Eg: Nomad will use update-index to safeguard deployment
	Deploy(jobSpec *Spec, extendVars map[string]interface{}) (deployResult *Spec, err error)

	GenerateJobSpecV1(levantJobFile string, jobCfg Config) (jobSpec *Spec, err error)
}
