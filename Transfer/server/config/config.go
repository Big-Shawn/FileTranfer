package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	Port        string   `yaml:"port"`
	Path        []string `yaml:"path"`
	PackageSize int      `yaml:"package_size"`
	Logger      Log      `yaml:"logger"`
}

type Log struct {
	Path        []string `yaml:"path"`
	Level       int8     `yaml:"level"`
	Development bool     `yaml:"development"`
	Encoding    string   `yaml:"encoding"`
}

var C *Config

const path = "config.yaml"

func Load() error {

	open, err := os.Open(path)

	if err != nil {
		return err
	}
	all, err := io.ReadAll(open)

	if err != nil {
		return err
	}
	return yaml.Unmarshal(all, &C)
}
