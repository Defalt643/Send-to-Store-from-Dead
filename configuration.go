package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"

	"veteran.socialenable.co/se4/middlelibrary/graylog2"
)

// Configuration object
type Configuration struct {
	Rabbit  Rabbit            `json:"rabbit"`
	Graylog graylog2.Graylog2 `json:"graylog"`
}

// GetConfiguration is get config follow environment
func GetConfiguration() Configuration {
	_, filename, _, _ := runtime.Caller(1)
	graylog2.Info(filename)
	switch Env {
	case BetaEnvironment, ProductionEnvironment:
		configFile, _ := os.Open(path.Join(path.Dir(filename), fmt.Sprintf("config/%s.json", Env)))
		return decodeConfig(configFile)
	default:
		Env = DevelopmentEnvironment
		configFile, _ := os.Open(path.Join(path.Dir(filename), fmt.Sprintf("config/%s.json", DevelopmentEnvironment)))
		return decodeConfig(configFile)
	}
}
func decodeConfig(file *os.File) Configuration {
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		message := fmt.Sprintf("Decode config error because : %s", err)
		graylog2.Fatal(13400, message)
		panic(message)
	}
	return configuration
}
