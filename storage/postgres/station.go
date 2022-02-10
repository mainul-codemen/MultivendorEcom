package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const insertstn = `
INSERT INTO station(
	name, 
	status,
	district_id,
	position,
	created_by,
	updated_by
) VALUES (
	:name, 
	:status,
	:district_id,
	:position,
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) CreateStation(con context.Context, des storage.Station) (string, error) {
	logger.Info("create station db")
	stmt, err := s.db.PrepareNamed(insertstn)
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

var stnlist = `
SELECT 
  station.id,
  station.name,
  station.position,
  district.name AS district_name,
  station.status,
  station.created_at,
  station.created_by,
  station.updated_at,
  station.updated_by
FROM station
LEFT JOIN district ON district.id = district_id
WHERE station.deleted_at IS NULL
`

func (s *Storage) GetStationList(ctx context.Context, sts bool) ([]storage.Station, error) {
	logger.Info("get all station")
	if sts {
		stnlist = stnlist + " AND station.status=1"
	}
	stn := make([]storage.Station, 0)
	if err := s.db.Select(&stn, stnlist); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return stn, nil
}

const gstn = `
SELECT
	station.id,
	station.name,
	station.position,
	district_id,
	district.name AS district_name,
	station.status,
	station.created_at,
	station.created_by,
	station.updated_at,
	station.updated_by
FROM station
LEFT JOIN district ON district.id = district_id
WHERE (station.id = $1 OR station.name = $1) AND  station.deleted_at IS NULL
`

func (s *Storage) GetStationBy(ctx context.Context, idname string) (*storage.Station, error) {
	var stn storage.Station
	if err := s.db.Get(&stn, gstn, idname); err != nil {
		logger.Error("error while get station data. " + err.Error())
		return nil, fmt.Errorf("executing station details: %w", err)
	}
	return &stn, nil
}

const gtstnBypn = `
SELECT
	id,
	name, 
	district_id, 
	status,
	position,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM station
WHERE position = $1 AND deleted_at IS NULL
`

func (s *Storage) GetStationByPosition(ctx context.Context, pos int32) (*storage.Station, error) {
	var stn storage.Station
	if err := s.db.Get(&stn, gtstnBypn, pos); err != nil {
		logger.Error("error while get district data. " + err.Error())
		return nil, fmt.Errorf("executing district details: %w", err)
	}
	return &stn, nil
}

const updateStn = `
UPDATE station SET
	name = :name,
	status = :status,
	district_id = :district_id,
	position = :position,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING *
`

func (s *Storage) UpdateStation(ctx context.Context, stn storage.Station) (*storage.Station, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateStn)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&stn, stn); err != nil {
		return nil, fmt.Errorf("executing station update: %w", err)
	}

	return &stn, nil
}

const updateStnStatus = `
	UPDATE 
		station 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateStationStatus(ctx context.Context, stn storage.Station) (*storage.Station, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateStnStatus)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&stn, stn); err != nil {
		return nil, fmt.Errorf("executing station status: %w", err)
	}

	return &stn, nil
}

const delStn = `
UPDATE
	station
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteStation(ctx context.Context, userid string, duid string) error {
	logger.Info("delete station")
	_, err := s.db.Exec(delStn, duid, userid)
	if err != nil {
		logger.Error("delete station")
		return err
	}
	return nil
}
