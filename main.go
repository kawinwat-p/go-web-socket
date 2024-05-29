package main

import (
	"bn-survey-point/configuration"
	ds "bn-survey-point/domain/datasources"
	repo "bn-survey-point/domain/repositories"
	gw "bn-survey-point/src/gateways"
	"bn-survey-point/src/middlewares"
	sv "bn-survey-point/src/services"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {

	// // remove this before deploy ###################
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	/// ############################################

	app := fiber.New(configuration.NewFiberConfiguration())
	middlewares.Logger(app)
	app.Use(recover.New())
	app.Use(cors.New())

	mongodb := ds.NewMongoDB(10)
	redisdb := ds.NewRedisConnection()

	userMongo := repo.NewUsersRepository(mongodb)
	alertRepo := repo.NewAlertMessageRepositories(mongodb)
	redisRepo := repo.NewRedisRepository(redisdb)


	sv0 := sv.NewSurveyPointService(userMongo,alertRepo,redisRepo)

	gw.NewHTTPGateway(app, sv0)

	PORT := os.Getenv("DB_PORT_LOGIN")

	if PORT == "" {
		PORT = "8080"
	}

	app.Listen(":" + PORT)
}
