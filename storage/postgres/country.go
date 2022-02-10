package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const insertcountry = `
INSERT INTO country(
	name, 
	status,
	position,
	created_by,
	updated_by
) VALUES (
	:name, 
	:status,
	:position,
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) CreateCountry(con context.Context, cntry storage.Country) (string, error) {
	logger.Info("create country db")
	stmt, err := s.db.PrepareNamed(insertcountry)
	if err != nil {
		logger.Error(ewpq + err.Error())
		return "", err
	}

	var id string
	if err := stmt.Get(&id, cntry); err != nil {
		logger.Error(ewmq + err.Error())
		return "", err
	}
	return id, nil
}

func (s *Storage) GetAllCountry(ctx context.Context, sts bool) ([]storage.Country, error) {
	logger.Info("get all Country")
	cntryl := `SELECT * from country WHERE deleted_at IS NULL`
	if sts {
		cntryl = cntryl + " AND status=1"
	}
	cntry := make([]storage.Country, 0)
	if err := s.db.Select(&cntry, cntryl); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return cntry, nil
}

const delcntry = `
UPDATE
	country
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteCountry(ctx context.Context, userid string, duid string) error {
	logger.Info("delete country")
	_, err := s.db.Exec(delcntry, duid, userid)
	if err != nil {
		logger.Error("delete country")
		return err
	}
	return nil
}

const getcntry = `
SELECT
	id,
	name, 
	status,
	position,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM country
WHERE (id = $1 OR name = $1) AND  deleted_at IS NULL
`

func (s *Storage) GetCountryBy(ctx context.Context, idname string) (*storage.Country, error) {
	var cntry storage.Country
	if err := s.db.Get(&cntry, getcntry, idname); err != nil {
		logger.Error("error while get country data. " + err.Error())
		return nil, fmt.Errorf("executing country details: %w", err)
	}
	return &cntry, nil
}

const gdbp = `
SELECT
	id,
	name, 
	status,
	position,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM country
WHERE position = $1 AND deleted_at IS NULL
`

func (s *Storage) GetCountryByPosition(ctx context.Context, pos int32) (*storage.Country, error) {
	var cntry storage.Country
	if err := s.db.Get(&cntry, gdbp, pos); err != nil {
		logger.Error("error while get country data. " + err.Error())
		return nil, fmt.Errorf("executing country details: %w", err)
	}
	return &cntry, nil
}

const updateCntry = `
UPDATE country SET
	name = :name,
	status = :status,
	position = :position,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING *
`

func (s *Storage) UpdateCountry(ctx context.Context, cnt storage.Country) (*storage.Country, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateCntry)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&cnt, cnt); err != nil {
		return nil, fmt.Errorf("executing country update: %w", err)
	}
	return &cnt, nil
}

const updateCntStatus = `
	UPDATE 
		country 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateCountryStatus(ctx context.Context, p storage.Country) (*storage.Country, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateCntStatus)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing country status: %w", err)
	}
	return &p, nil
}
