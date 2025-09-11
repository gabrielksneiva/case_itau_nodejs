package api

import (
	"case-itau/api/handler"
	"case-itau/api/middleware"
	"case-itau/config"

	_ "case-itau/docs"
	repository "case-itau/repository/interface"
	"case-itau/services/customer"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"gorm.io/gorm"
)

func Register(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	// apply middlewares
	middleware.RegisterMiddlewares(app, int64(cfg.RateLimitMax))

	// auto migrate
	db.AutoMigrate(&repository.Clientes{})

	// init services and handlers
	svc := customer.NewService(db)
	h := handler.NewCustomerHandler(svc)

	// swagger
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html", fiber.StatusFound)
	})
	app.Get("/docs/", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html", fiber.StatusFound)
	})
	app.Get("/docs/*", fiberSwagger.WrapHandler)

	// routes
	v1 := app.Group("/clientes")
	v1.Get("/", h.List)
	v1.Get("/:id", h.Get)
	v1.Post("/", h.Create)
	v1.Put("/:id", h.Update)
	v1.Delete("/:id", h.Delete)
	v1.Post("/:id/depositar", h.Deposit)
	v1.Post("/:id/sacar", h.Withdraw)
}
