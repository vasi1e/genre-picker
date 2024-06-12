package presenter

import (
	"fmt"
	"genre-pick-up/internal/model"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Repository interface {
	GetPeople() ([]string, error)
	GetPersonGenre(name string) (model.GenrePartners, error)
	GetGenres() ([]model.GenrePartners, error)
	SavePickedGenre(name string, genre string) error
	SaveGenre(model.GenrePartners) error
}

type presenter struct {
	repo   Repository
	people []string
}

func NewPresenter(repo Repository) (presenter, error) {
	people, err := repo.GetPeople()
	if err != nil {
		return presenter{}, fmt.Errorf("failed to get geople: %w", err)
	}

	return presenter{
		repo:   repo,
		people: people,
	}, nil
}

func (p presenter) getFreeGenres() ([]model.GenrePartners, error) {
	var free []model.GenrePartners
	genres, err := p.repo.GetGenres()
	if err != nil {
		return nil, fmt.Errorf("failed to get genre: %w", err)
	}

	for _, genrePartners := range genres {
		if len(genrePartners.Partners) < 2 {
			free = append(free, genrePartners)
		}
	}

	return free, nil
}

func (p presenter) pickGenre(name string) (model.GenrePartners, error) {
	freeGenres, err := p.getFreeGenres()
	if err != nil {
		return model.GenrePartners{}, fmt.Errorf("failed to get free genre: %w", err)
	}

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(len(freeGenres))

	genrePartners := freeGenres[randNum]
	genrePartners.Partners = append(genrePartners.Partners, name)

	if err := p.repo.SavePickedGenre(name, genrePartners.Genre); err != nil {
		return model.GenrePartners{}, fmt.Errorf("failed to save picked genre: %w", err)
	}

	if err := p.repo.SaveGenre(genrePartners); err != nil {
		return model.GenrePartners{}, fmt.Errorf("failed to save genre: %w", err)
	}

	return genrePartners, nil
}

func (p presenter) Pick(ctx *fiber.Ctx) error {
	name := ctx.FormValue("name")
	partner := "..."
	genre := "..."

	if name != "" {
		genrePartners, err := p.repo.GetPersonGenre(name)
		if err != nil {
			panic(fmt.Errorf("failed to get person genre: %v", err))
		}

		genre = genrePartners.Genre

		if genre == "" {
			genrePartners, err = p.pickGenre(name)
			if err != nil {
				panic(fmt.Errorf("failed to pick genre: %v", err))
			}

			genre = genrePartners.Genre
		}

		partner = getPartner(name, genrePartners)
	}

	return ctx.Render("index", fiber.Map{
		"people":   p.people,
		"selected": name,
		"with":     partner,
		"genre":    genre,
	})
}

func getPartner(name string, genrePartners model.GenrePartners) string {
	if len(genrePartners.Partners) == 1 {
		return "..."
	}

	if genrePartners.Partners[0] == name {
		return genrePartners.Partners[1]
	}

	return genrePartners.Partners[0]
}
