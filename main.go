package main

import (
	"net/http"

	"github.com/LukmanulHakim18/gorooster/v2/database"
	"github.com/LukmanulHakim18/gorooster/v2/helpers"
	"github.com/LukmanulHakim18/gorooster/v2/router"
	"github.com/LukmanulHakim18/gorooster/v2/services"

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
