package main

import (
	"git.bluebird.id/mybb/gorooster/database"
	"git.bluebird.id/mybb/gorooster/helpers"
	"git.bluebird.id/mybb/gorooster/router"
	"git.bluebird.id/mybb/gorooster/services"
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
