package cfg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type DbInfo struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	DbName   string `json:"dbname"`
}

type LoggingInfo struct {
	Level    string `json:"level"`
	Filepath string `json:"filepath"`
}

type AppConfig struct {
	Database DbInfo      `json:"database"`
	Logging  LoggingInfo `json:"logging"`
}

func NewAppConfig() AppConfig {
	return AppConfig{
		DbInfo{
			"",
			"",
			"",
			"",
			"",
		},
		LoggingInfo{
			"",
			"",
		},
	}
}

var appCfg *AppConfig

func GetAppConfig() *AppConfig {
	if appCfg == nil {
		return nil
	}

	return appCfg
}

func LoadAppConfig(filepath string) error {

	jsonFile, err := os.Open(filepath)

	if err != nil {
		return errors.New(fmt.Sprintf("Unable to load app config: path=%s: %s", filepath, err))
	}

	defer jsonFile.Close()

	cfg := NewAppConfig()

	byteValue, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal JSON to app config: %s", err))
	}

	appCfg = &cfg

	return nil
}
