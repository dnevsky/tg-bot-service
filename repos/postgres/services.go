package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/dnevsky/tg-bot-service/models"
	"github.com/jmoiron/sqlx"
)

type Services struct {
	db *sqlx.DB
}

func NewServices(db *sqlx.DB) *Services {
	return &Services{db}
}

func (r *Services) Save(from int64, service, login, password string) error {
	var id int

	query := fmt.Sprintf(`INSERT INTO %s (userid, nameofservice, login, password)
							SELECT * FROM (
								SELECT $1::bigint, $2, $3, $4
							) AS new_service
							WHERE NOT EXISTS (
								SELECT 1 FROM services WHERE userid = $5::bigint AND nameofservice = $6
							) LIMIT 1 RETURNING id;`, servicesTable)
	row := r.db.QueryRow(query, from, service, login, password, from, service)
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("Вы пытаетесь создать ещё один сервис с таким же названием.")
		}
		return err
	}

	return nil
}

func (r *Services) Read(from int64, service string) (string, string, error) {
	var serviceObj models.Service

	query := fmt.Sprintf(`SELECT * FROM %s WHERE userid = $1 AND nameofservice = $2`, servicesTable)

	if err := r.db.Get(&serviceObj, query, int(from), service); err != nil {
		return "", "", err
	}

	return serviceObj.Login, serviceObj.Password, nil
}

func (r *Services) Delete(from int64, service string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE userid = $1 AND nameofservice = $2", servicesTable)

	res, err := r.db.Exec(query, from, service)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Services) GetAll(from int64) ([]string, error) {
	var res_temp []models.Service
	var res []string

	query := fmt.Sprintf("SELECT * FROM %s WHERE userid = $1", servicesTable)

	if err := r.db.Select(&res_temp, query, from); err != nil {
		return nil, err
	}

	for _, v := range res_temp {
		res = append(res, v.NameOfService)
	}

	return res, nil
}
