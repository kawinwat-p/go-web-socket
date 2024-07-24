package main

import (
	"os"
	"websocketjingjing/configuration"
	repo "websocketjingjing/domain/repositories"
	gw "websocketjingjing/src/gateways"
	"websocketjingjing/src/middlewares"
	sv "websocketjingjing/src/services"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	// // // remove this before deploy ###################
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	// /// ############################################

	app := fiber.New(configuration.NewFiberConfiguration())
	middlewares.Logger(app)
	app.Use(recover.New())
	app.Use(cors.New())

	// app.Use("/ws", func(c *fiber.Ctx) error {
	// 	if websocket.IsWebSocketUpgrade(c) {
	// 		c.Locals("allowed", true)
	// 		return c.Next()
	// 	}
	// 	return fiber.ErrUpgradeRequired
	// })

	config := websocket.Config{
		HandshakeTimeout:  0,             // No timeout for handshake
		Origins:           []string{"*"}, // Allow all origins
		EnableCompression: false,         // Disable compression
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
	}

	hub := repo.NewHub()
	sv1 := sv.NewHubService(hub)

	gw.NewHTTPGateway(app, sv1,config)

	PORT := os.Getenv("DB_PORT_LOGIN")

	if PORT == "" {
		PORT = "8080"
	}

	app.Listen(":" + PORT)
}
