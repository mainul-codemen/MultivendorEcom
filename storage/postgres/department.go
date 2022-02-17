package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const insertdpt = `
INSERT INTO department(
	name, 
	description, 
	status,
	position,
	created_by,
	updated_by
) VALUES (
	:name, 
	:description,
	:status, 
	:position,
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) CreateDepartment(con context.Context, des storage.Department) (string, error) {
	logger.Info("create department db")
	stmt, err := s.db.PrepareNamed(insertdpt)
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

func (s *Storage) GetDepartment(ctx context.Context, sts bool) ([]storage.Department, error) {
	logger.Info("get all department")
	desq := `SELECT * from department WHERE deleted_at IS NULL`
	if sts {
		desq = desq + " AND status=1"
	}
	des := make([]storage.Department, 0)
	if err := s.db.Select(&des, desq); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return des, nil
}

const deldpt = `
UPDATE
	department
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteDepartment(ctx context.Context, userid string, duid string) error {
	logger.Info("delete department")
	_, err := s.db.Exec(deldpt, duid, userid)
	if err != nil {
		logger.Error("delete department")
		return err
	}
	return nil
}

const gdptquery = `
SELECT
	id,
	name, 
	description, 
	status,
	position,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM department
WHERE (id = $1 OR name = $1) AND  deleted_at IS NULL
`

func (s *Storage) GetDepartmentBy(ctx context.Context, idname string) (*storage.Department, error) {
	var dpt storage.Department
	if err := s.db.Get(&dpt, gdptquery, idname); err != nil {
		logger.Error("error while get department data. " + err.Error())
		return nil, fmt.Errorf("executing department details: %w", err)
	}
	return &dpt, nil
}

const gdptPos = `
SELECT
	id,
	name, 
	description, 
	status,
	position,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM department
WHERE position = $1 AND deleted_at IS NULL
`

func (s *Storage) GetDepartmentByPosition(ctx context.Context, pos int32) (*storage.Department, error) {
	var des storage.Department
	if err := s.db.Get(&des, gdptPos, pos); err != nil {
		logger.Error("error while get department data. " + err.Error())
		return nil, fmt.Errorf("executing department details: %w", err)
	}
	return &des, nil
}

const updateDept = `
UPDATE department SET
	name = :name,
	description = :description,
	status = :status,
	position = :position,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING *
`

func (s *Storage) UpdateDepartment(ctx context.Context, p storage.Department) (*storage.Department, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateDept)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing department update: %w", err)
	}
	return &p, nil
}

const updateDptStatus = `
	UPDATE 
		department 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateDepartmentStatus(ctx context.Context, p storage.Department) (*storage.Department, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateDptStatus)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing department status: %w", err)
	}
	return &p, nil
}
