package config

type DockerImageResolver interface {
	Resolve() string
}
