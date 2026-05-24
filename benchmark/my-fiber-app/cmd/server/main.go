package main

import (
	"my-fiber-app/pkg/config"
	"my-fiber-app/pkg/mongo"
	"my-fiber-app/pkg/redis"
	"my-fiber-app/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()

	// telemetry.InitTracing(context.Background()) // <--- adiciona o tracer
	mongo.Connect()
	redis.Connect()

	app := fiber.New()
	// app.Use(otelfiber.Middleware())
	// telemetry.InitMetrics(app)

	router.SetupRoutes(app)

	app.Listen(":3001")
}
