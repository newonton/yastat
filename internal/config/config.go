package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	dirName  = ".config/yastat"
	fileName = "config.json"
)

type Config struct {
	APIKey string `json:"api_key"`
	AppID  int    `json:"app_id"`
}

func MustLoad() (*Config, error) {
	return Load()
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := createEmptyConfig(path); err != nil {
			return nil, err
		}
		return nil, firstRunMessage(path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config format")
	}

	if cfg.APIKey == "" || cfg.AppID == 0 {
		return nil, fmt.Errorf(
			"config is not filled\nedit %s\n\n"+
				"Get api_key:\n"+
				"https://oauth.yandex.ru/authorize?response_type=token&client_id=da092e6d50b443308da7a28e638070b9",
			path,
		)
	}

	return &cfg, nil
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, dirName, fileName), nil
}

func createEmptyConfig(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	empty := `{
  "api_key": "",
  "app_id": 0
}

`
	return os.WriteFile(path, []byte(empty), 0o644)
}

func firstRunMessage(path string) error {
	return fmt.Errorf(
		`config file created: %s

Fill it with your data:

{
  "api_key": "YOUR_TOKEN",
  "app_id": 123456
}

Get api_key:
https://oauth.yandex.ru/authorize?response_type=token&client_id=da092e6d50b443308da7a28e638070b9`,
		path)
}
