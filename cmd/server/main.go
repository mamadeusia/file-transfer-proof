package main

import (
	"log"

	config "github.com/mamadeusia/file-transfer-proof/config/server"
	"github.com/mamadeusia/file-transfer-proof/internal/server/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	log.Println("config: ", cfg)
	// Run
	app.Run(cfg)
}
