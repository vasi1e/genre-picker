package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"genre-pick-up/internal/model"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type sqlDB struct {
	connection *sql.DB
}

func NewInstance(user, pass, host, dbName string) (*sqlDB, error) {
	dbConnection, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, dbName))
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection with DB: %w", err)
	}

	if err := dbConnection.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return &sqlDB{
		connection: dbConnection,
	}, nil
}

func (sdb *sqlDB) GetPeople() ([]string, error) {
	results, err := sdb.connection.Query(selectPeople)
	if err != nil {
		return nil, fmt.Errorf("failed to execute selecting all people: %w", err)
	}
	defer results.Close()

	var people []string
	for results.Next() {
		var person []byte
		if err := results.Scan(&person); err != nil {
			return nil, fmt.Errorf("failed to scan person: %w", err)
		}

		people = append(people, string(person))
	}

	return people, nil
}

func (sdb *sqlDB) GetPersonGenre(name string) (model.GenrePartners, error) {
	stmt, err := sdb.connection.Prepare(selectPersonGenre)
	if err != nil {
		return model.GenrePartners{}, fmt.Errorf("failed to prepare selecting person genre and partners: %w", err)
	}

	var genre, partners []byte
	if err := stmt.QueryRow(name).Scan(&genre, &partners); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return model.GenrePartners{}, fmt.Errorf("failed to scan person genre and partners: %w", err)
	}

	return model.GenrePartners{
		Genre:    string(genre),
		Partners: splitPartners(string(partners)),
	}, nil
}

func (sdb *sqlDB) GetGenres() ([]model.GenrePartners, error) {
	results, err := sdb.connection.Query(selectGenrePartners)
	if err != nil {
		return nil, fmt.Errorf("failed to execute selecting all genres: %w", err)
	}
	defer results.Close()

	var genrePartners []model.GenrePartners
	for results.Next() {
		var genre, partners []byte
		if err := results.Scan(&genre, &partners); err != nil {
			return nil, fmt.Errorf("failed to scan genre and partners: %w", err)
		}

		genrePartners = append(genrePartners, model.GenrePartners{
			Genre:    string(genre),
			Partners: splitPartners(string(partners)),
		})
	}

	return genrePartners, nil
}

func (sdb *sqlDB) SavePickedGenre(name string, genre string) error {
	stmt, err := sdb.connection.Prepare(updatePersonGenre)
	if err != nil {
		return fmt.Errorf("failed to prapare updating person genre: %w", err)
	}

	_, err = stmt.Exec(genre, name)
	if err != nil {
		return fmt.Errorf("failed to execute updating person genre: %w", err)
	}

	return nil
}

func (sdb *sqlDB) SaveGenre(genrePartners model.GenrePartners) error {
	stmt, err := sdb.connection.Prepare(updateGenrePartners)
	if err != nil {
		return fmt.Errorf("failed to prapare updating genre partners: %w", err)
	}

	_, err = stmt.Exec(joinPartners(genrePartners.Partners), genrePartners.Genre)
	if err != nil {
		return fmt.Errorf("failed to execute updating genre partners: %w", err)
	}

	return nil
}

func splitPartners(partners string) []string {
	if partners == "" {
		return []string{}
	}

	return strings.Split(partners, ", ")
}

func joinPartners(partners []string) string {
	return strings.Join(partners, ", ")
}
