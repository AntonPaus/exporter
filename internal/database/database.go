package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

type Migration struct {
	Name    string
	SQL     string
	Applied bool
}

// loadMigrations reads all SQL files from the migrations directory
func LoadMigrations(dir string) ([]Migration, error) {
	var migrations []Migration
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			sqlContent, err := os.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}
			migrations = append(migrations, Migration{
				Name: file.Name(),
				SQL:  string(sqlContent),
			})
		}
	}
	return migrations, nil
}

// applyMigrations executes the SQL migrations
func ApplyMigrations(db *sql.DB, migrations []Migration) error {
	fmt.Println("Applying migrations:")
	for _, migration := range migrations {
		// Check if migration is already applied
		if !migration.Applied {
			fmt.Printf("\t%s\n", migration.Name)
			// Execute the SQL
			_, err := db.Exec(migration.SQL)
			if err != nil {
				return err
			}
			migration.Applied = true
		}
	}
	return nil
}
