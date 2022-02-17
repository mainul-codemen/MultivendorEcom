package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const irole = `
INSERT INTO user_role(
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

func (s *Storage) CreateUserRole(con context.Context, des storage.UserRole) (string, error) {
	logger.Info("create user_role db")
	stmt, err := s.db.PrepareNamed(irole)
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

func (s *Storage) GetUserRole(ctx context.Context, sts bool) ([]storage.UserRole, error) {
	logger.Info("get all user_role")
	desq := `SELECT * from user_role WHERE deleted_at IS NULL`
	if sts {
		desq = desq + " AND status=1"
	}
	des := make([]storage.UserRole, 0)
	if err := s.db.Select(&des, desq); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return des, nil
}

const delur = `
UPDATE
	user_role
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteUserRole(ctx context.Context, userid string, duid string) error {
	logger.Info("delete user_role")
	_, err := s.db.Exec(delur, duid, userid)
	if err != nil {
		logger.Error("delete user_role")
		return err
	}
	return nil
}

const gdursquery = `
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
FROM user_role
WHERE (id = $1 OR name = $1) AND  deleted_at IS NULL
`

func (s *Storage) GetUserRoleBy(ctx context.Context, idname string) (*storage.UserRole, error) {
	var dpt storage.UserRole
	if err := s.db.Get(&dpt, gdursquery, idname); err != nil {
		logger.Error("error while get user_role data. " + err.Error())
		return nil, fmt.Errorf("executing user_role details: %w", err)
	}
	return &dpt, nil
}

const gdpusrPos = `
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
FROM user_role
WHERE position = $1 AND deleted_at IS NULL
`

func (s *Storage) GetUserRoleByPosition(ctx context.Context, pos int32) (*storage.UserRole, error) {
	var des storage.UserRole
	if err := s.db.Get(&des, gdpusrPos, pos); err != nil {
		logger.Error("error while get user_role data. " + err.Error())
		return nil, fmt.Errorf("executing user_role details: %w", err)
	}
	return &des, nil
}

const updateUser = `
UPDATE user_role SET
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

func (s *Storage) UpdateUserRole(ctx context.Context, p storage.UserRole) (*storage.UserRole, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateUser)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing user_role update: %w", err)
	}
	return &p, nil
}

const updateUsrStatus = `
	UPDATE 
		user_role 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateUserRoleStatus(ctx context.Context, p storage.UserRole) (*storage.UserRole, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateUsrStatus)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing user_role status: %w", err)
	}
	return &p, nil
}
