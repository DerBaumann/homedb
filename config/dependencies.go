package config

import (
	"database/sql"
	"homedb/repository"
	"os"
)

type dependencies struct {
	Repo *repository.Queries
}

func SetupDependencies() (*dependencies, error) {
	db, err := sql.Open("postgres", os.Getenv("DB_STRING"))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	repo := repository.New(db)

	deps := dependencies{
		Repo: repo,
	}

	return &deps, nil
}
