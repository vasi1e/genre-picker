package repository

const selectPeople = `SELECT name FROM people;`

const selectGenrePartners = `SELECT * FROM genres;`

const selectPersonGenre = `SELECT p.genre, g.partners
	FROM people p JOIN genres g ON p.genre = g.name
	WHERE p.name = ?;`

const updateGenrePartners = `UPDATE genres SET partners = ? WHERE name = ?;`

const updatePersonGenre = `UPDATE people SET genre = ? WHERE name = ?;`
