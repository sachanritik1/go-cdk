package app

import (
	"lambda-func/api"
	"lambda-func/database"
)

type App struct {
	ApiHandler api.ApiHandler
}

func NewApp() *App {
	db, err := database.NewDynamoDBClient()
	if err != nil {
		panic(err)
	}

	apiHandler := api.NewApiHandler(*db)

	return &App{
		ApiHandler: *apiHandler,
	}
}
