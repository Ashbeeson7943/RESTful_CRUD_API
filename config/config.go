package internalConfig

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Config struct {
	DB_URI     string    `json:"dbURI"`
	API_CONFIG ApiConfig `json:"api_config"`
}

type ApiConfig struct {
	HOST string `json:"host"`
	PORT string `json:"port"`
	BASE_PATH string `json:"base_path"`
}

func (ac *ApiConfig) FullURL() string {
	return fmt.Sprintf("%v:%v", ac.HOST, ac.PORT)
}

func LoadConfig(configFilePath string) Config {
	configJson, err := os.Open(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	jsonByteValue, _ := io.ReadAll(configJson)

	var conf Config

	json.Unmarshal(jsonByteValue, &conf)

	defer configJson.Close()

	return conf
}