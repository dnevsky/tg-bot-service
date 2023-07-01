package repos

import (
	"github.com/dnevsky/tg-bot-service/repos/postgres"
	"github.com/jmoiron/sqlx"
)

type Services interface {
	Save(from int64, service, login, password string) error
	Read(from int64, service string) (string, string, error)
	Delete(from int64, service string) error
	GetAll(from int64) ([]string, error)
}

type Repos struct {
	Services
}

func NewRepos(db *sqlx.DB) *Repos {
	return &Repos{
		Services: postgres.NewServices(db),
	}
}

// func NewRepos() *Repos {
// 	return &Repos{
// 		Services: file.NewServices(),
// 	}
// }
