package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const inserthub = `
INSERT INTO hub(
	hub_name, 
	status,
	hub_phone_1,
	hub_phone_2,
	hub_email,
	hub_address,
	country_id,
	district_id,
	station_id,
	position,
	created_by,
	updated_by
) VALUES (
	:hub_name, 
	:status,
	:hub_phone_1,
	:hub_phone_2,
	:hub_email,
	:hub_address,
	:country_id,
	:district_id,
	:station_id,
	:position,
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) CreateHub(con context.Context, des storage.Hub) (string, error) {
	logger.Info("create hub db")
	stmt, err := s.db.PrepareNamed(inserthub)
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

var hublist = `
SELECT 
  hub.id,
  hub.hub_name,
  hub.position,
  hub.hub_phone_1,
  hub.hub_phone_2,
  hub.hub_email,
  hub.hub_address,
  country.id as country_id,
  country.name AS country_name,
  district.id as district_id,
  district.name AS district_name,
  station.id as station_id,
  station.name AS station_name,
  hub.status,
  hub.created_at,
  hub.created_by,
  hub.updated_at,
  hub.updated_by
FROM hub
LEFT JOIN country ON country.id = country_id
LEFT JOIN district ON district.id = district_id
LEFT JOIN station ON station.id = station_id
WHERE hub.deleted_at IS NULL
`

func (s *Storage) GetHubList(ctx context.Context, sts bool) ([]storage.Hub, error) {
	logger.Info("get all hub")
	if sts {
		hublist = hublist + " AND status=1"
	}
	hub := make([]storage.Hub, 0)
	if err := s.db.Select(&hub, hublist); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return hub, nil
}

const ghub = `
SELECT
	hub.id,
	hub.hub_name,
	hub.position,
	hub.country_id,
	hub.district_id,
	hub.station_id,
	hub.hub_phone_1,
	hub.hub_phone_2,
	hub.hub_email,
	hub.hub_address,
	country.name AS country_name,
	district.name AS district_name,
	station.name AS station_name,
	hub.status,
	hub.created_at,
	hub.created_by,
	hub.updated_at,
	hub.updated_by
FROM hub
LEFT JOIN country ON country.id = country_id
LEFT JOIN district ON district.id = district_id
LEFT JOIN station ON station.id = station_id
WHERE (hub.id = $1 OR hub.hub_name = $1 OR hub.hub_phone_1=$1 OR hub.hub_phone_2=$1 OR hub.hub_email=$1) AND  hub.deleted_at IS NULL
`

func (s *Storage) GetHubBy(ctx context.Context, idname string) (*storage.Hub, error) {
	var hub storage.Hub
	if err := s.db.Get(&hub, ghub, idname); err != nil {
		logger.Error("error while get hub data. " + err.Error())
		return nil, fmt.Errorf("executing hub details: %w", err)
	}
	return &hub, nil
}

const gthubBypn = `
SELECT
	id,
	hub_name, 
	hub_phone_1,
	hub_phone_2,
	hub_email,
	hub_address,
	country_id, 
	district_id, 
	station_id, 
	status,
	position,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM hub
WHERE position = $1 AND deleted_at IS NULL
`

func (s *Storage) GetHubByPosition(ctx context.Context, pos int32) (*storage.Hub, error) {
	var hub storage.Hub
	if err := s.db.Get(&hub, gthubBypn, pos); err != nil {
		logger.Error("error while get district data. " + err.Error())
		return nil, fmt.Errorf("executing district details: %w", err)
	}
	return &hub, nil
}

const updateHub = `
UPDATE hub SET
	hub_name = :hub_name,
	status = :status,
	hub_phone_1 = :hub_phone_1,
	hub_phone_2 = :hub_phone_2,
	hub_email = :hub_email,
	hub_address = :hub_address,
	country_id = :country_id,
	district_id = :district_id,
	station_id = :station_id,
	position = :position,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING *
`

func (s *Storage) UpdateHub(ctx context.Context, hub storage.Hub) (*storage.Hub, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateHub)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&hub, hub); err != nil {
		return nil, fmt.Errorf("executing hub update: %w", err)
	}

	return &hub, nil
}

const updateHubStatus = `
	UPDATE 
		hub 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateHubStatus(ctx context.Context, hub storage.Hub) (*storage.Hub, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateHubStatus)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&hub, hub); err != nil {
		return nil, fmt.Errorf("executing hub status: %w", err)
	}

	return &hub, nil
}

const delHub = `
UPDATE
	hub
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteHub(ctx context.Context, userid string, duid string) error {
	logger.Info("delete hub")
	_, err := s.db.Exec(delHub, duid, userid)
	if err != nil {
		logger.Error("delete hub")
		return err
	}
	return nil
}
