package main

import (
	"gorooster/database"
	"gorooster/helpers"
	"gorooster/router"
	"gorooster/services"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	RedisClient := database.GetRedisClient()

	go services.StartEventListener(*RedisClient)

	if err := http.ListenAndServe(helpers.EnvGetString("RUNING_PORT", ":1407"), router.GetRouter()); err != nil {
		panic(err)
	}

}
