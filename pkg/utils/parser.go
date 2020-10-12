package utils

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Generate(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("parser config fail, %v", err)
	}
	cfg := &Config{}
	return cfg, yaml.Unmarshal(data, &cfg)
}
