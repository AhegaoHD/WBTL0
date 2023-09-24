package main

import (
	"log"
	"os"

	"github.com/AhegaoHD/WBTL0/config"
	"github.com/AhegaoHD/WBTL0/internal/wb0"
)

func main() {
	configPath := findConfigPath()

	cfg, err := config.Parse(configPath)
	if err != nil {
		log.Fatal(err)
	}

	wb0.Run(cfg)
}
func findConfigPath() string {
	const (
		devConfig  = "config/dev.config.toml"
		prodConfig = "config/config.toml"
	)

	if os.Getenv("CFG") == "DEV" {
		return devConfig
	}

	return prodConfig
}
