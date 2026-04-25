package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	conf := Config{}
	if err := json.Unmarshal(data, &conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}

func (conf *Config) SetUser(currentUserName string) error {
	conf.CurrentUserName = currentUserName
	err := write(*conf)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, configFileName), nil
}

func write(conf Config) error {
	data, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	file, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(file, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
