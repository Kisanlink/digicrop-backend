package config

import (
	"log"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	OpenAIKey string `yaml:"openai_api_key"`
}

var AppConfig Config

func LoadConfig() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Load YAML file
	yamlFile, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("Error reading config.yaml: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &AppConfig)
	if err != nil {
		log.Fatalf("Error unmarshalling config.yaml: %v", err)
	}
}

func GetOpenAIKey() string {
	return AppConfig.OpenAIKey
}
