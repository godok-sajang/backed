package config

import (
	"godok/clean/cmd/cli"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Listen     string   `yaml:"listen"`
	Env        string   `yaml:"env"`
	Postgresql DBConfig `yaml:"postgresql"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbName"`
}

func New(flags *cli.Flags) *Config {
	return ParseConfig(flags.ConfigPath)
}

func ParseConfig(path string) *Config {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error when reading file from config-path")
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		log.Fatalf("error when unmarshaling config-file: %v", err)
	}

	return c
}
