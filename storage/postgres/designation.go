package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const insertdes = `
INSERT INTO designation(
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

func (s *Storage) CreateDesignation(con context.Context, des storage.Designation) (string, error) {
	logger.Info("create designation db")
	stmt, err := s.db.PrepareNamed(insertdes)
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

func (s *Storage) GetDesignation(ctx context.Context, sts bool) ([]storage.Designation, error) {
	logger.Info("get all designation")
	desq := `SELECT * from designation WHERE deleted_at IS NULL`
	if sts {
		desq = desq + " AND status=1"
	}
	des := make([]storage.Designation, 0)
	if err := s.db.Select(&des, desq); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return des, nil
}

const deldes = `
UPDATE
	designation
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteDesignation(ctx context.Context, userid string, duid string) error {
	logger.Info("delete designation")
	_, err := s.db.Exec(deldes, duid, userid)
	if err != nil {
		logger.Error("delete designation")
		return err
	}
	return nil
}

const gdquery = `
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
FROM designation
WHERE (id = $1 OR name = $1) AND  deleted_at IS NULL
`

func (s *Storage) GetDesignationBy(ctx context.Context, idname string) (*storage.Designation, error) {
	var des storage.Designation
	if err := s.db.Get(&des, gdquery, idname); err != nil {
		logger.Error("error while get designation data. " + err.Error())
		return nil, fmt.Errorf("executing designation details: %w", err)
	}
	return &des, nil
}

const gpq = `
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
FROM designation
WHERE position = $1 AND deleted_at IS NULL
`

func (s *Storage) GetDesignationByPosition(ctx context.Context, pos int32) (*storage.Designation, error) {
	var des storage.Designation
	if err := s.db.Get(&des, gpq, pos); err != nil {
		logger.Error("error while get designation data. " + err.Error())
		return nil, fmt.Errorf("executing designation details: %w", err)
	}
	return &des, nil
}

const updateDesg = `
UPDATE designation SET
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

func (s *Storage) UpdateDesignation(ctx context.Context, p storage.Designation) (*storage.Designation, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateDesg)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing designation update: %w", err)
	}
	return &p, nil
}

const updateStatus = `
	UPDATE 
		designation 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateDesignationStatus(ctx context.Context, p storage.Designation) (*storage.Designation, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateStatus)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing designation status: %w", err)
	}
	return &p, nil
}
