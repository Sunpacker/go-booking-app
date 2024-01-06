package dbrepo

import (
	"database/sql"
	"github.com/Sunpacker/go-booking-app/internal/config"
	"github.com/Sunpacker/go-booking-app/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

type testDBRepo struct {
	App *config.AppConfig
}

func NewTestRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}