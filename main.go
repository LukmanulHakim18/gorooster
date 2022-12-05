package main

import (
	"github.com/LukmanulHakim18/gorooster/database"
	"github.com/LukmanulHakim18/gorooster/helpers"
	"github.com/LukmanulHakim18/gorooster/router"
	"github.com/LukmanulHakim18/gorooster/services"
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
