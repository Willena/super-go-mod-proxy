package config

import (
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

var logger, _ = zap.NewDevelopment()

type Config struct {
	General AppConfig           `json:"general"`
	Phases  PhasesConfiguration `json:"phases"`
}

type AppConfig struct {
	Secure            bool   `json:"secure"`
	DefaultRelayProxy string `json:"defaultRelayProxy"`
}

type PluginsDefinitions map[string]PluginDefinition
type PluginConfiguration map[string]interface{}

type PluginDefinition struct {
	Kind   string              `json:"kind"`
	Config PluginConfiguration `json:"config"`
}

type PhasesConfiguration struct {
	Receive  PluginsDefinitions `json:"receive,omitempty"`
	PreFetch PluginsDefinitions `json:"prefetch,omitempty"`
	Fetch    PluginsDefinitions `json:"fetch,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	logger.Info("Loading Main configuration")

	jsonFile, err := os.Open(path)
	if err != nil {
		logger.Error("Could not open configuration ", zap.Error(err), zap.String("FilePath", path))
		return nil, err
	}
	defer jsonFile.Close()

	logger.Debug("Reading configuration...")
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logger.Error("Could not read configuration file", zap.Error(err))
		return nil, err
	}

	var conf Config

	err = json.Unmarshal(byteValue, &conf)
	if err != nil {
		logger.Error("Could not decode Json configuration", zap.Error(err))
		return nil, err
	}

	logger.Info("Configuration Loaded")
	logger.Debug("Configuration summary is ", zap.Any("config", conf))

	return &conf, nil
}
