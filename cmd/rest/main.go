package main

import (
	"api/cmd/rest/application"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// env
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	
	// app
	config := application.NewConfigDefault()
	router := chi.NewRouter()

	app := application.NewApp(config, router)
	if err := app.Dependencies(); err != nil {
		panic(err)
	}

	// run
	if err := app.Run(); err != nil {
		panic(err)
	}
}