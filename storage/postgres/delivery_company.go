package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const insertdelivery_company = `
INSERT INTO delivery_company(
	company_name, 
	company_status,
	phone,
	email,
	company_address,
	country_id,
	district_id,
	station_id,
	position,
	created_by,
	updated_by
) VALUES (
	:company_name, 
	:company_status,
	:phone,
	:email,
	:company_address,
	:country_id,
	:district_id,
	:station_id,
	:position,
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) CreateDeliveryCompany(con context.Context, des storage.DeliveryCompany) (string, error) {
	logger.Info("create delivery_company db")
	stmt, err := s.db.PrepareNamed(insertdelivery_company)
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

const delivery_companylist = `
SELECT 
  delivery_company.id,
  delivery_company.company_name,
  delivery_company.position,
  delivery_company.phone,
  delivery_company.email,
  delivery_company.company_address,
  country.name AS country_id,
  district.name AS district_id,
  station.name AS station_id,
  delivery_company.company_status,
  delivery_company.created_at,
  delivery_company.created_by,
  delivery_company.updated_at,
  delivery_company.updated_by
FROM delivery_company
LEFT JOIN country ON country.id = country_id
LEFT JOIN district ON district.id = district_id
LEFT JOIN station ON station.id = station_id
WHERE delivery_company.deleted_at IS NULL
`

func (s *Storage) GetDeliveryCompanyList(ctx context.Context) ([]storage.DeliveryCompany, error) {
	logger.Info("get all delivery_company")
	dc := make([]storage.DeliveryCompany, 0)
	if err := s.db.Select(&dc, delivery_companylist); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return dc, nil
}

const gdelivery_company = `
SELECT
	delivery_company.id,
	delivery_company.company_name,
	delivery_company.position,
	delivery_company.country_id,
	delivery_company.company_id,
	delivery_company.station_id,
	delivery_company.phone,
	delivery_company.email,
	delivery_company.company_address,
	country.name AS country_name,
	district.name AS district_name,
	station.name AS station_name,
	delivery_company.company_status,
	delivery_company.created_at,
	delivery_company.created_by,
	delivery_company.updated_at,
	delivery_company.updated_by
FROM delivery_company
LEFT JOIN country ON country.id = country_id
LEFT JOIN district ON delivery_company.id = district_id
LEFT JOIN station ON station.id = station_id
WHERE (delivery_company.id = $1 OR delivery_company.delivery_company_name = $1 OR delivery_company_phone_1 OR delivery_company_phone_2 OR delivery_company_email) AND  delivery_company.deleted_at IS NULL
`

func (s *Storage) GetDeliveryCompanyBy(ctx context.Context, idname string) (*storage.DeliveryCompany, error) {
	var dc storage.DeliveryCompany
	if err := s.db.Get(&dc, gdelivery_company, idname); err != nil {
		logger.Error("error while get delivery_company data. " + err.Error())
		return nil, fmt.Errorf("executing delivery_company details: %w", err)
	}
	return &dc, nil
}

const gtdelivery_companyBypn = `
SELECT
	id,
    company_name, 
	phone,
	email,
	company_address,
	country_id, 
	district_id, 
	station_id, 
	company_status,
	position,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM delivery_company
WHERE position = $1 AND deleted_at IS NULL
`

func (s *Storage) GetDeliveryCompanyByPosition(ctx context.Context, pos int32) (*storage.DeliveryCompany, error) {
	var dc storage.DeliveryCompany
	if err := s.db.Get(&dc, gtdelivery_companyBypn, pos); err != nil {
		logger.Error("error while get delivery_company data. " + err.Error())
		return nil, fmt.Errorf("executing delivery_company details: %w", err)
	}
	return &dc, nil
}

const updateDeliveryCompany = `
UPDATE delivery_company SET
	company_name = :company_name,
	company_status = :company_status,
	phone = :phone,
	email = :email,
	company_address = :company_address,
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

func (s *Storage) UpdateDeliveryCompany(ctx context.Context, dc storage.DeliveryCompany) (*storage.DeliveryCompany, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateDeliveryCompany)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&dc, dc); err != nil {
		return nil, fmt.Errorf("executing delivery_company update: %w", err)
	}

	return &dc, nil
}

const updateDeliveryCompanyStatus = `
	UPDATE 
		delivery_company 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateDeliveryCompanyStatus(ctx context.Context, dc storage.DeliveryCompany) (*storage.DeliveryCompany, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateDeliveryCompanyStatus)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&dc, dc); err != nil {
		return nil, fmt.Errorf("executing delivery_company status: %w", err)
	}

	return &dc, nil
}

const delDeliveryCompany = `
UPDATE
	delivery_company
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteDeliveryCompany(ctx context.Context, userid string, duid string) error {
	logger.Info("delete delivery_company")
	_, err := s.db.Exec(delDeliveryCompany, duid, userid)
	if err != nil {
		logger.Error("delete delivery_company")
		return err
	}
	return nil
}
