package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const insertbranch = `
INSERT INTO branch(
	branch_name, 
	status,
	branch_phone_1,
	branch_phone_2,
	branch_email,
	branch_address,
	country_id,
	branch_id,
	station_id,
	position,
	created_by,
	updated_by
) VALUES (
	:branch_name, 
	:status,
	:branch_phone_1,
	:branch_phone_2,
	:branch_email,
	:branch_address,
	:country_id,
	:branch_id,
	:station_id,
	:position,
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) CreateBranch(con context.Context, des storage.Branch) (string, error) {
	logger.Info("create branch db")
	stmt, err := s.db.PrepareNamed(insertbranch)
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

const branchlist = `
SELECT 
  branch.id,
  branch.branch_name,
  branch.position,
  branch.branch_phone_1,
  branch.branch_phone_2,
  branch.branch_email,
  branch.branch_address,
  country.name AS country_id,
  district.name AS district_id,
  station.name AS station_id,
  branch.branch_status,
  branch.created_at,
  branch.created_by,
  branch.updated_at,
  branch.updated_by
FROM branch
LEFT JOIN country ON country.id = country_id
LEFT JOIN district ON district.id = district_id
LEFT JOIN station ON station.id = station_id
WHERE branch.deleted_at IS NULL
`

func (s *Storage) GetBranchList(ctx context.Context) ([]storage.Branch, error) {
	logger.Info("get all branch")
	branch := make([]storage.Branch, 0)
	if err := s.db.Select(&branch, branchlist); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return branch, nil
}

const gbranch = `
SELECT
	branch.id,
	branch.branch_name,
	branch.position,
	branch.country_id,
	branch.branch_id,
	branch.station_id,
	branch.branch_phone_1,
	branch.branch_phone_2,
	branch.branch_email,
	branch.branch_address,
	country.name AS country_name,
	district.name AS district_name,
	station.name AS station_name,
	branch.status,
	branch.created_at,
	branch.created_by,
	branch.updated_at,
	branch.updated_by
FROM branch
LEFT JOIN country ON country.id = country_id
LEFT JOIN district ON branch.id = district_id
LEFT JOIN station ON station.id = station_id
WHERE (branch.id = $1 OR branch.branch_name = $1 OR branch_phone_1 OR branch_phone_2 OR branch_email) AND  branch.deleted_at IS NULL
`

func (s *Storage) GetBranchBy(ctx context.Context, idname string) (*storage.Branch, error) {
	var branch storage.Branch
	if err := s.db.Get(&branch, gbranch, idname); err != nil {
		logger.Error("error while get branch data. " + err.Error())
		return nil, fmt.Errorf("executing branch details: %w", err)
	}
	return &branch, nil
}

const gtbranchBypn = `
SELECT
	id,
	branch_name, 
	branch_phone_1,
	branch_phone_2,
	branch_email,
	branch_address,
	country_id, 
	district_id, 
	station_id, 
	status,
	position,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM branch
WHERE position = $1 AND deleted_at IS NULL
`

func (s *Storage) GetBranchByPosition(ctx context.Context, pos int32) (*storage.Branch, error) {
	var branch storage.Branch
	if err := s.db.Get(&branch, gtbranchBypn, pos); err != nil {
		logger.Error("error while get branch data. " + err.Error())
		return nil, fmt.Errorf("executing branch details: %w", err)
	}
	return &branch, nil
}

const updateBranch = `
UPDATE branch SET
	name = :name,
	status = :status,
	branch_phone_1 = :branch_phone_1,
	branch_phone_2 = :branch_phone_2,
	branch_email = :branch_email,
	branch_address = :branch_address,
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

func (s *Storage) UpdateBranch(ctx context.Context, branch storage.Branch) (*storage.Branch, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateBranch)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&branch, branch); err != nil {
		return nil, fmt.Errorf("executing branch update: %w", err)
	}

	return &branch, nil
}

const updateBranchStatus = `
	UPDATE 
		branch 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateBranchStatus(ctx context.Context, branch storage.Branch) (*storage.Branch, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateBranchStatus)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&branch, branch); err != nil {
		return nil, fmt.Errorf("executing branch status: %w", err)
	}

	return &branch, nil
}

const delBranch = `
UPDATE
	branch
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteBranch(ctx context.Context, userid string, duid string) error {
	logger.Info("delete branch")
	_, err := s.db.Exec(delBranch, duid, userid)
	if err != nil {
		logger.Error("delete branch")
		return err
	}
	return nil
}
