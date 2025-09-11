package api

import (
	"case-itau/config"
	"case-itau/repository/connection"
	repository "case-itau/repository/interface"
	"case-itau/utils/logger"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Start() {
	cfg := config.Load()

	logger.NewLogger()

	db, err := connection.NewSqliteConnection(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	err = db.AutoMigrate(&repository.Clientes{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: "Customer API",
	})

	Register(app, db, cfg)

	log.Printf("listening on %s", cfg.APIPort)
	log.Fatal(app.Listen(":" + cfg.APIPort))
}
