package main

import (
	"os"
	"websocketjingjing/configuration"
	ds "websocketjingjing/domain/datasources"
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

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	mongodb := ds.NewMongoDB(10)
	redisdb := ds.NewRedisConnection()
	

	userMongo := repo.NewUsersRepository(mongodb)
	alertRepo := repo.NewAlertMessageRepositories(mongodb)
	redisRepo := repo.NewRedisRepository(redisdb)
	// socketRepo := repo.NewSocketRepo()
	hub := repo.NewHub()

	sv0 := sv.NewSurveyPointService(userMongo, alertRepo, redisRepo)
	// sv1 := sv.NewSocketService(socketRepo)
	sv2 := sv.NewHubService(hub)

	gw.NewHTTPGateway(app, sv0, sv2)

	PORT := os.Getenv("DB_PORT_LOGIN")

	if PORT == "" {
		PORT = "8080"
	}

	app.Listen(":" + PORT)
}
