package deployer

type Storer interface {
	Store(projectId string, jobSpec, deployResultSpec *Spec, meta map[string]interface{}) error
	Close() error
}
