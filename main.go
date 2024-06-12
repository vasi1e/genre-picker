package main

import (
	"genre-pick-up/env"
	"genre-pick-up/internal/presenter"
	"genre-pick-up/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		Immutable: true,
		Views:     html.New("./views", ".html"),
	})

	dbConfig, err := env.LoadDBConfig()
	if err != nil {
		panic(err)
	}

	repo, err := repository.NewInstance(dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Name)
	if err != nil {
		panic(err)
	}

	pre, err := presenter.NewPresenter(repo)
	if err != nil {
		panic(err)
	}

	app.Get("/pick", pre.Pick)
	app.Post("/pick", pre.Pick)

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}
