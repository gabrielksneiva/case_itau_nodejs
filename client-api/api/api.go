package api

import (
	"case-itau/api/handler"
	"case-itau/config"
	"case-itau/repositories"
	"case-itau/repositories/connection"
	"case-itau/services/customer"
	l "case-itau/utils/logger"

	"github.com/gofiber/fiber/v2"
)

func Start() {
	cfg := config.Load()

	app := fiber.New(fiber.Config{
		AppName: "Customer API",
	})

	// init db
	db, err := connection.NewSqliteConnection(cfg.DBPath)
	if err != nil {
		l.Logger.Sugar().Fatalf("failed to connect database: %v", err)
	}
	err = db.AutoMigrate(&repositories.Clientes{})
	if err != nil {
		l.Logger.Sugar().Fatalf("failed to migrate database: %v", err)
	}

	// init repo
	repo := repositories.NewGormRepository[repositories.Clientes](db)

	// init services and handlers
	svc := customer.NewService(repo)
	h := handler.NewCustomerHandler(svc)

	Register(app, db, cfg, h)

	l.Logger.Sugar().Fatal(app.Listen(":" + cfg.APIPort))
}
