package main

import (
	"net/http"

	"git.bluebird.id/mybb/gorooster/v2/database"
	"git.bluebird.id/mybb/gorooster/v2/helpers"
	"git.bluebird.id/mybb/gorooster/v2/router"
	"git.bluebird.id/mybb/gorooster/v2/services"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	RedisClient := database.GetRedisClient()

	go services.StartEventListeners(RedisClient)

	if err := http.ListenAndServe(helpers.EnvGetString("RUNING_PORT", ":1407"), router.GetRouter()); err != nil {
		panic(err)
	}

}
