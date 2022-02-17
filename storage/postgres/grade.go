package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const insertgrd = `
INSERT INTO grade(
	name, 
	basic_salary, 
	lunch_allowance, 
	transportation, 
	rent_allowance, 
	absent_penalty, 
	total_salary, 
	status,
	position,
	created_by,
	updated_by
) VALUES (
	:name, 
	:basic_salary,
	:lunch_allowance,
	:transportation,
	:rent_allowance,
	:absent_penalty,
	:total_salary,
	:status, 
	:position,
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) CreateGrade(con context.Context, des storage.Grade) (string, error) {
	logger.Info("create grade db")
	stmt, err := s.db.PrepareNamed(insertgrd)
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

func (s *Storage) GetGrade(ctx context.Context, sts bool) ([]storage.Grade, error) {
	logger.Info("get all grade")
	desq := `SELECT * from grade WHERE deleted_at IS NULL`
	if sts {
		desq = desq + " AND status=1"
	}
	des := make([]storage.Grade, 0)
	if err := s.db.Select(&des, desq); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return des, nil
}

const delgrade = `
UPDATE
	grade
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteGrade(ctx context.Context, userid string, duid string) error {
	logger.Info("delete grade")
	_, err := s.db.Exec(delgrade, duid, userid)
	if err != nil {
		logger.Error("delete grade")
		return err
	}
	return nil
}

const gdGrdquery = `
SELECT
	id,
	name, 
	basic_salary, 
	lunch_allowance, 
	transportation, 
	rent_allowance, 
	absent_penalty, 
	total_salary, 
	status,
	position,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM grade
WHERE (id = $1 OR name = $1) AND  deleted_at IS NULL
`

func (s *Storage) GetGradeBy(ctx context.Context, idname string) (*storage.Grade, error) {
	var dpt storage.Grade
	if err := s.db.Get(&dpt, gdGrdquery, idname); err != nil {
		logger.Error("error while get grade data. " + err.Error())
		return nil, fmt.Errorf("executing grade details: %w", err)
	}
	return &dpt, nil
}

const ggrdPos = `
SELECT
	id,
	name, 
	basic_salary, 
	lunch_allowance, 
	transportation, 
	rent_allowance, 
	absent_penalty, 
	total_salary, 
	status,
	position,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM grade
WHERE position = $1 AND deleted_at IS NULL
`

func (s *Storage) GetGradeByPosition(ctx context.Context, pos int32) (*storage.Grade, error) {
	var des storage.Grade
	if err := s.db.Get(&des, ggrdPos, pos); err != nil {
		logger.Error("error while get grade data. " + err.Error())
		return nil, fmt.Errorf("executing grade details: %w", err)
	}
	return &des, nil
}

const updateGrd = `
UPDATE grade SET
	name = :name,
	basic_salary = :basic_salary,
	lunch_allowance = :lunch_allowance,
	transportation = :transportation,
	rent_allowance = :rent_allowance,
	absent_penalty = :absent_penalty,
	total_salary = :total_salary,
	status = :status,
	position = :position,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING *
`

func (s *Storage) UpdateGrade(ctx context.Context, p storage.Grade) (*storage.Grade, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateGrd)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing grade update: %w", err)
	}
	return &p, nil
}

const updateGrdStatus = `
	UPDATE 
		grade 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateGradeStatus(ctx context.Context, p storage.Grade) (*storage.Grade, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateGrdStatus)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing grade status: %w", err)
	}
	return &p, nil
}
