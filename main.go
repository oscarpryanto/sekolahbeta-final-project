package main

import (
	"final-project/bckp-database/config"
	"final-project/bckp-database/controllers"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func InitEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Warn("Cannot load env file, using system env")
	}
}

func main() {

	InitEnv()
	config.OpenDB()

	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024, //50mb
	})

	// controllers.RouteCars(app)
	controllers.RouteBckpDatabase(app)

	err := app.Listen(":3000")
	if err != nil {
		logrus.Fatal(
			"Error on running fiber, ",
			err.Error())
	}

	_ = os.Mkdir("upload", 0777)
}
