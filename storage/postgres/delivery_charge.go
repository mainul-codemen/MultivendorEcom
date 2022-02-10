package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const insertdc = `
INSERT INTO delivery_charge(
	country_id,
	district_id,
	station_id,
	weight_min,
	weight_max,
	delivery_charge,
	created_by,
	updated_by
) VALUES (
	:country_id,
	:district_id,
	:station_id,
	:weight_min,
	:weight_max,
	:delivery_charge,
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) CreateDeliveryCharge(con context.Context, des storage.DeliveryCharge) (string, error) {
	logger.Info("create delivery_charge db")
	stmt, err := s.db.PrepareNamed(insertdc)
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

const deliveryChargelist = `
SELECT 
  delivery_charge.id,
  country.name AS country_id,
  district.name AS district_id,
  station.name AS station_id,
  delivery_charge.dc_status,
  delivery_charge.created_at,
  delivery_charge.created_by,
  delivery_charge.updated_at,
  delivery_charge.updated_by
FROM delivery_charge
LEFT JOIN country ON country.id = country_id
LEFT JOIN district ON district.id = district_id
LEFT JOIN station ON station.id = station_id
WHERE delivery_charge.deleted_at IS NULL
`

func (s *Storage) GetDeliveryChargeList(ctx context.Context) ([]storage.DeliveryCharge, error) {
	logger.Info("get all delivery_charge")
	dc := make([]storage.DeliveryCharge, 0)
	if err := s.db.Select(&dc, deliveryChargelist); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return dc, nil
}

const gdc = `
SELECT
	delivery_charge.id,
	delivery_charge.country_id,
	delivery_charge.district_id,
	delivery_charge.station_id,
	country.name AS country_name,
	district.name AS district_name,
	station.name AS station_name,
	delivery_charge.dc_status,
	delivery_charge.created_at,
	delivery_charge.created_by,
	delivery_charge.updated_at,
	delivery_charge.updated_by
FROM delivery_charge
LEFT JOIN country ON country.id = country_id
LEFT JOIN district ON district.id = district_id
LEFT JOIN station ON station.id = station_id
WHERE delivery_charge.id = $1 AND  delivery_charge.deleted_at IS NULL
`

func (s *Storage) GetDeliveryChargeBy(ctx context.Context, idname string) (*storage.DeliveryCharge, error) {
	var deliveryCharge storage.DeliveryCharge
	if err := s.db.Get(&deliveryCharge, gdc, idname); err != nil {
		logger.Error("error while get delivery_charge data. " + err.Error())
		return nil, fmt.Errorf("executing delivery_charge details: %w", err)
	}
	return &deliveryCharge, nil
}

const updateDeliveryCharge = `
UPDATE delivery_charge SET
	country_id = :country_id,
	district_id = :district_id,
	station_id = :station_id,
	dc_status = :dc_status,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING *
`

func (s *Storage) UpdateDeliveryCharge(ctx context.Context, dc storage.DeliveryCharge) (*storage.DeliveryCharge, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateDeliveryCharge)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&dc, dc); err != nil {
		return nil, fmt.Errorf("executing delivry charge update: %w", err)
	}

	return &dc, nil
}

const updateDeliveryChargeStatus = `
	UPDATE 
		delivery_charge 
	SET
		dc_status = :dc_status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateDeliveryChargeStatus(ctx context.Context, dc storage.DeliveryCharge) (*storage.DeliveryCharge, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateDeliveryChargeStatus)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&dc, dc); err != nil {
		return nil, fmt.Errorf("executing delivery_charge status: %w", err)
	}

	return &dc, nil
}

const delDeliveryCharge = `
UPDATE
	delivery_charge
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteDeliveryCharge(ctx context.Context, userid string, duid string) error {
	logger.Info("delete delivery_charge")
	_, err := s.db.Exec(delDeliveryCharge, duid, userid)
	if err != nil {
		logger.Error("delete delivery_charge")
		return err
	}
	return nil
}
