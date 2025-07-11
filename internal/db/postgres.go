package db

import (
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/model"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var Database *sql.DB

func InitPostgres() error {
	connectionStr := os.Getenv("POSTGRES_ENV")
	db, err := sql.Open("postgres", connectionStr)

	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	Database = db
	return createTableIfNotExists()
}
func createTableIfNotExists() error {
	query := `
	CREATE TABLE IF NOT EXISTS characters (
		id INT PRIMARY KEY,
		name TEXT NOT NULL,
		status TEXT,
		species TEXT,
		origin TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := Database.Exec(query)
	return err
}

func InsertCharacters(chars []model.Character) error {
	tx, err := Database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //safer rollback option, it is completely optional

	stmt, err := tx.Prepare(`INSERT INTO characters (id, name, status, species, origin) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (id) DO NOTHING;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range chars {
		if _, err := stmt.Exec(c.ID, c.Name, c.Status, c.Species, c.Origin.Name); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func GetCharacters(limit, offset int, sortBy string) ([]model.Character, error) {
	if sortBy != "name" && sortBy != "id" {
		sortBy = "id" // Fallback to default
	}
	query := fmt.Sprintf(`
		SELECT id, name, status, species, origin
		FROM characters
		ORDER BY %s
		LIMIT $1 OFFSET $2
	`, sortBy)

	rows, err := Database.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.Character
	for rows.Next() {
		var c model.Character
		err := rows.Scan(&c.ID, &c.Name, &c.Status, &c.Species, &c.Origin.Name)
		if err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	return results, nil

}
