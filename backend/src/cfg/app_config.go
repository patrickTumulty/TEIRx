package cfg 

import (
	"encoding/json"

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

type AppConfig struct {
	Database DbInfo `json:"database"`
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
	}
}

func LoadAppConfig(filepath string) (*AppConfig, error) {

	jsonFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	appConfig := NewAppConfig()

	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &appConfig)

	return &appConfig, nil
}
