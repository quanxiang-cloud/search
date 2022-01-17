package config

import (
	"context"
	"io/ioutil"

	"github.com/quanxiang-cloud/cabin/tailormade/db/elastic"
	"github.com/quanxiang-cloud/search/pkg/util"
	"gopkg.in/yaml.v2"
)

// Config configuration item
type Config struct {
	Port          string         `yaml:"port"`
	Elasticsearch elastic.Config `yaml:"elasticsearch"`
}

// New reuturn config from file path
func New(ctx context.Context, path string) (*Config, error) {
	log := util.LoggerFromContext(ctx).WithName("config")

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf *Config
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		log.Error(err, "yaml unmarshal")
		return nil, err
	}

	return conf, nil
}
