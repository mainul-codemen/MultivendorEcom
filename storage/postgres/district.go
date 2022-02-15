package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const insertdis = `
INSERT INTO district(
	name, 
	status,
	country_id,
	position,
	created_by,
	updated_by
) VALUES (
	:name, 
	:status,
	:country_id,
	:position,
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) CreateDistrict(con context.Context, des storage.District) (string, error) {
	logger.Info("create district db")
	stmt, err := s.db.PrepareNamed(insertdis)
	if err != nil {
		logger.Error(ewpq + err.Error())
		return "", err
	}

	var id string
	if err := stmt.Get(&id, des); err != nil {
		logger.Error(ewmq + err.Error())
		return "", err
	}
	return id, nil
}

var dislist = `
SELECT 
  district.id,
  district.name,
  district.position,
  country.name AS country_name,
  district.status,
  district.created_at,
  district.created_by,
  district.updated_at,
  district.updated_by
FROM district
LEFT JOIN country ON country.id = country_id
WHERE district.deleted_at IS NULL
`

func (s *Storage) GetDistrictList(ctx context.Context, sts bool) ([]storage.District, error) {
	logger.Info("get all district")
	fmt.Println(sts)
	if sts {
		dislist = dislist + " AND district.status=1"
	}

	dis := make([]storage.District, 0)
	if err := s.db.Select(&dis, dislist); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return dis, nil
}

const gdisquery = `
SELECT
	district.id,
	district.name,
	district.position,
	country_id,
	country.name AS country_name,
	district.status,
	district.created_at,
	district.created_by,
	district.updated_at,
	district.updated_by
FROM district
LEFT JOIN country ON country.id = country_id
WHERE (district.id = $1 OR district.name = $1) AND  district.deleted_at IS NULL
`

func (s *Storage) GetDistrictBy(ctx context.Context, idname string) (*storage.District, error) {
	var des storage.District
	if err := s.db.Get(&des, gdisquery, idname); err != nil {
		logger.Error("error while get district data. " + err.Error())
		return nil, fmt.Errorf("executing district details: %w", err)
	}
	return &des, nil
}

const gdispq = `
SELECT
	id,
	name, 
	country_id, 
	status,
	position,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM district
WHERE position = $1 AND deleted_at IS NULL
`

func (s *Storage) GetDistrictByPosition(ctx context.Context, pos int32) (*storage.District, error) {
	var des storage.District
	if err := s.db.Get(&des, gdispq, pos); err != nil {
		logger.Error("error while get district data. " + err.Error())
		return nil, fmt.Errorf("executing district details: %w", err)
	}
	return &des, nil
}

const updateDist = `
UPDATE district SET
	name = :name,
	status = :status,
	country_id = :country_id,
	position = :position,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING *
`

func (s *Storage) UpdateDistrict(ctx context.Context, dis storage.District) (*storage.District, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateDist)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&dis, dis); err != nil {
		return nil, fmt.Errorf("executing district update: %w", err)
	}

	return &dis, nil
}

const updateDisStatus = `
	UPDATE 
		district 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateDistrictStatus(ctx context.Context, dis storage.District) (*storage.District, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateDisStatus)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&dis, dis); err != nil {
		return nil, fmt.Errorf("executing district status: %w", err)
	}

	return &dis, nil
}

const deldis = `
UPDATE
	district
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteDistrict(ctx context.Context, userid string, duid string) error {
	logger.Info("delete district")
	_, err := s.db.Exec(deldis, duid, userid)
	if err != nil {
		logger.Error("delete district")
		return err
	}
	return nil
}
