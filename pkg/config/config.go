package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/fr0stylo/imployer/pkg/build"
	"github.com/fr0stylo/imployer/pkg/deploy"
)

type Config struct {
	Build  *build.BuildOpts   `yaml:"build"`
	Deploy *deploy.DeployOpts `yaml:"deploy"`
}

func Load(path string) (*Config, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer fd.Close()

	var cfg Config
	if err := yaml.NewDecoder(fd).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %v", err)
	}

	return &cfg, nil
}
