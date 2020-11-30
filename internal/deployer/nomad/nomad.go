package nomad

import (
	"encoding/json"
	"fmt"
	nomadapi "github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec"
	levant "github.com/jrasell/levant/template"
	"github.com/rs/zerolog/log"
	"github.com/Awesome-Nomad/Nomad-Deployer/internal/beauter"
	"github.com/Awesome-Nomad/Nomad-Deployer/internal/deployer"
	mutils "github.com/Awesome-Nomad/Nomad-Deployer/utils"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var chDirLock = &sync.Mutex{}

type nomadDeployer struct {
	nomadClient *nomadapi.Client
	consulAddr  string
}

func (d *nomadDeployer) GenerateJobSpecV1(levantJobFile string, jobCfg deployer.Config) (jobSpec *deployer.Spec, err error) {
	var jobFile *os.File
	jobFile, err = mutils.NewTempFile(fmt.Sprintf("%s.*.nomad", jobCfg.ServiceName))
	if err != nil {
		return
	}

	_jobSpec := jobCfg.GetSpec()
	err = ioutil.WriteFile(jobFile.Name(), []byte(_jobSpec.Content), 0644)

	return
}

func NewNomadDeployer(nomadCfg *nomadapi.Config, consulAddr string) (deployer deployer.Deployer, err error) {
	nomadCli, err := nomadapi.NewClient(nomadCfg)
	if err != nil {
		return
	}
	return &nomadDeployer{
		nomadClient: nomadCli,
		consulAddr:  consulAddr,
	}, nil
}

// GenerateJobSpec generate nomad job as HCL file.
// jobTemplate 	must be a valid levant template
// jobConfigs	a list of yaml files to input as levant var-file
// workingDir	directory to run template file. Because template may use local file of project
func (d *nomadDeployer) GenerateJobSpec(jobTemplateFile string, jobConfigFiles []string, workingDir string) (jobSpec *deployer.Spec, err error) {
	chDirLock.Lock()
	// Only unlock when other defers finished running
	defer chDirLock.Unlock()
	// 1. Change working directory to target directory
	curWD, _ := os.Getwd()
	err = os.Chdir(workingDir)
	// Back to right working dir
	defer func() {
		e := os.Chdir(curWD)
		if e != nil {
			panic(e)
		}
	}()

	if err != nil {
		return nil, err
	}
	levantFlag := make(map[string]string)
	tpl, err := levant.RenderTemplate(jobTemplateFile, jobConfigFiles, d.consulAddr, &levantFlag)
	log.Err(err)
	if err != nil {
		return nil, err
	}
	jobContent := tpl.String()
	// 2. Verify Job is valid
	_, err = jobspec.Parse(tpl)
	if err != nil {
		return nil, err
	}
	jobSpec = &deployer.Spec{
		Content:     jobContent,
		ContentType: beauter.HCL,
	}
	return
}

func (d *nomadDeployer) Diff(jobSpec *deployer.Spec) (diff *deployer.Spec, meta map[string]interface{}, err error) {
	job, err := jobspec.Parse(strings.NewReader(jobSpec.Content))
	if err != nil {
		return nil, nil, err
	}
	resp, _, err := d.nomadClient.Jobs().Plan(job, true, nil)
	if err != nil {
		return nil, nil, err
	}
	meta = make(map[string]interface{})
	meta[deployer.NomadPlanIndexMetaKey] = resp.JobModifyIndex
	meta[deployer.RawResponse] = resp
	jsonb, err := json.Marshal(resp)
	if err != nil {
		return nil, nil, err
	}

	return &deployer.Spec{
		Content:     string(jsonb),
		ContentType: beauter.JSON,
	}, meta, nil
}

func (d *nomadDeployer) Deploy(jobSpec *deployer.Spec, extendVars map[string]interface{}) (deployResult *deployer.Spec, err error) {
	job, err := jobspec.Parse(strings.NewReader(jobSpec.Content))
	if err != nil {
		return nil, err
	}
	index, enforce := extendVars[deployer.NomadPlanIndexMetaKey]
	resp, _, err := d.nomadClient.Jobs().RegisterOpts(job, &nomadapi.RegisterOptions{
		EnforceIndex:   enforce,
		ModifyIndex:    index.(uint64),
		PolicyOverride: false,
	}, nil)
	if err != nil {
		return nil, err
	}
	jsonb, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	deployResult = &deployer.Spec{
		Content:     string(jsonb),
		ContentType: beauter.JSON,
	}

	return deployResult, nil
}
