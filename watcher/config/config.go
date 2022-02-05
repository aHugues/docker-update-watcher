package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type DockerLink int64

const (

	// DockerLive means that the daemon is directly linked to the Docker daemon and gets information from current images
	DockerLive DockerLink = iota

	// DockerStandalone means that the daemon is not linked to any Docker daemon and only has access to its initial data
	DockerStandalone
)

type LogConfig struct {
	Level logrus.Level `json:"level"`
}

type DockerStandaloneConfig struct {
	CacheDirectory string
	IntialListPath string `json:"initial-list"`
}

type DockerConfig struct {
	Username  string `json:"username"`
	AuthToken string `json:"auth-token"`
}

type Config struct {
	DockerLink DockerLink

	Logging LogConfig `json:"logging"`

	DockerConfig     DockerConfig           `json:"docker-hub"`
	StandaloneConfig DockerStandaloneConfig `json:"docker-standalone"`
}

var DefaultConfig = Config{
	DockerLink: DockerLive,
	Logging: LogConfig{
		Level: logrus.InfoLevel,
	},
}

func ReadFromFile(filepath string) (*Config, error) {
	confFile, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer confFile.Close()

	content, err := io.ReadAll(confFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	NewConf := Config(DefaultConfig)
	if err := json.Unmarshal(content, &NewConf); err != nil {
		return nil, fmt.Errorf("failed to parse JSON config file: %w", err)
	}
	return &NewConf, nil
}
