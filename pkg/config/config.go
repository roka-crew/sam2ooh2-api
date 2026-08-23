package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env    Env       `yaml:"env"`
	Name   string    `yaml:"name"`
	Listen string    `yaml:"listen"`
	Sqlite RDBConfig `yaml:"sqlite"`
}

func New(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open: %v", err)
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, fmt.Errorf("failed new decode: %v", err)
	}

	return config, nil
}

type Env uint8

const (
	EnvDev Env = iota
	EnvQa
	EnvProd
)

func (e *Env) UnmarshalYAML(value *yaml.Node) error {
	var decodedValue string
	if err := value.Decode(&decodedValue); err != nil {
		return err
	}

	switch decodedValue {
	case "dev", "development":
		*e = EnvDev
	case "qa":
		*e = EnvQa
	case "prod", "production":
		*e = EnvProd
	default:
		return fmt.Errorf("unknown env value: %s", decodedValue)
	}

	return nil
}

type RDBConfig struct {
	Host     string     `yaml:"host"`
	User     string     `yaml:"user"`
	Port     int        `yaml:"port"`
	Password string     `yaml:"password"`
	DBname   string     `yaml:"dbName"`
	Options  RDBOptions `yaml:"options"`
}

type RDBOptions struct {
	MaxIdleConns    int           `yaml:"maxIdleConns"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"`
}

func (c *RDBOptions) UnmarshalYAML(value *yaml.Node) error {
	var raw struct {
		MaxIdleConns    int    `yaml:"maxIdleConns"`
		MaxOpenConns    int    `yaml:"maxOpenConns"`
		ConnMaxLifetime string `yaml:"connMaxLifetime"`
	}

	if err := value.Decode(&raw); err != nil {
		return err
	}

	c.MaxIdleConns = raw.MaxIdleConns
	c.MaxOpenConns = raw.MaxOpenConns

	parsedDuration, err := time.ParseDuration(raw.ConnMaxLifetime)
	if err != nil {
		return err
	}
	c.ConnMaxLifetime = parsedDuration

	return nil
}
