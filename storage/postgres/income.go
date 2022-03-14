package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const insertInc = `
INSERT INTO income(
	title, 
	income_amount, 
	account_id,
	note,
	status,
	income_date,
	created_by,
	updated_by
) VALUES (
	:title, 
	:income_amount,
	:account_id, 
	:note,
	:status,
	:income_date,
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) CreateIncome(con context.Context, act storage.Income) (string, error) {
	logger.Info("create income db")
	stmt, err := s.db.PrepareNamed(insertInc)
	if err != nil {
		logger.Error(ewpq + err.Error())
		return "", err
	}

	var id string
	if err := stmt.Get(&id, act); err != nil {
		logger.Error(ewmq + err.Error())
		return "", err
	}
	return id, nil
}

func (s *Storage) GetIncome(ctx context.Context, sts bool) ([]storage.Income, error) {
	logger.Info("get all income")
	actq := `SELECT * from income WHERE deleted_at IS NULL`
	if sts {
		actq = actq + " AND status=1"
	}
	act := make([]storage.Income, 0)
	if err := s.db.Select(&act, actq); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return act, nil
}

const delInc = `
UPDATE
	income
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteIncome(ctx context.Context, userid string, duid string) error {
	logger.Info("delete income")
	_, err := s.db.Exec(delInc, duid, userid)
	if err != nil {
		logger.Error("delete income")
		return err
	}
	return nil
}

const gdinc = `
SELECT
	id,
	title, 
	income_amount, 
	account_id,
	note,
	status,
	income_date,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM income
WHERE (id = $1 OR title = $1) AND  deleted_at IS NULL
`

func (s *Storage) GetIncomeBy(ctx context.Context, idname string) (*storage.Income, error) {
	logger.Info("get income by")
	var dpt storage.Income
	if err := s.db.Get(&dpt, gdinc, idname); err != nil {
		logger.Error("error while get income data. " + err.Error())
		return nil, fmt.Errorf("executing income details: %w", err)
	}

	return &dpt, nil
}

const updateInc = `
UPDATE income SET
	title = :title,
	income_amount = :income_amount,
	account_id = :account_id,
	note = :note,
	status = :status,
	income_date = :income_date,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING *
`

func (s *Storage) UpdateIncome(ctx context.Context, act storage.Income) (*storage.Income, error) {
	logger.Info("update income info")
	stmt, err := s.db.PrepareNamedContext(ctx, updateInc)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&act, act); err != nil {
		return nil, fmt.Errorf("executing income update: %w", err)
	}

	return &act, nil
}

const updateIncStatus = `
	UPDATE 
		income 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateIncomeStatus(ctx context.Context, act storage.Income) (*storage.Income, error) {
	logger.Info("update income status")
	stmt, err := s.db.PrepareNamedContext(ctx, updateIncStatus)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&act, act); err != nil {
		return nil, fmt.Errorf("executing income status: %w", err)
	}

	return &act, nil
}
