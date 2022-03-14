package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const insertIncT = `
INSERT INTO income_tax(
	account_id, 
	tax_receipt_number, 
	status,
	income_tax_date,
	tax_amount,
	created_by,
	updated_by
) VALUES (
	:account_id, 
	:tax_receipt_number,
	:status, 
	:income_tax_date,
	:tax_amount,
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) CreateIncomeTax(con context.Context, act storage.IncomeTax) (string, error) {
	logger.Info("create income db")
	stmt, err := s.db.PrepareNamed(insertIncT)
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

func (s *Storage) GetIncomeTax(ctx context.Context, sts bool) ([]storage.IncomeTax, error) {
	logger.Info("get all income tax")
	actq := `SELECT * from income_tax WHERE deleted_at IS NULL`
	if sts {
		actq = actq + " AND status=1"
	}
	act := make([]storage.IncomeTax, 0)
	if err := s.db.Select(&act, actq); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return act, nil
}

const delIncT = `
UPDATE
	income_tax
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteIncomeTax(ctx context.Context, userid string, duid string) error {
	logger.Info("delete income tax")
	_, err := s.db.Exec(delIncT, duid, userid)
	if err != nil {
		logger.Error("delete income tax")
		return err
	}
	return nil
}

const gdincT = `
SELECT
	id,
	account_id, 
	tax_receipt_number, 
	status,
	income_tax_date,
	tax_amount,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM income_tax
WHERE (id = $1 OR account_id = $1) AND  deleted_at IS NULL
`

func (s *Storage) GetIncomeTaxBy(ctx context.Context, idname string) (*storage.IncomeTax, error) {
	logger.Info("get income tax by")
	var dpt storage.IncomeTax
	if err := s.db.Get(&dpt, gdincT, idname); err != nil {
		logger.Error("error while get income data. " + err.Error())
		return nil, fmt.Errorf("executing income details: %w", err)
	}
	return &dpt, nil
}

const updateIncT = `
UPDATE income_tax SET
	account_id = :account_id,
	tax_receipt_number = :tax_receipt_number,
	status = :status,
	income_tax_date = :income_tax_date,
	tax_amount = :tax_amount,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING *
`

func (s *Storage) UpdateIncomeTax(ctx context.Context, act storage.IncomeTax) (*storage.IncomeTax, error) {
	logger.Info("update income tax info")
	stmt, err := s.db.PrepareNamedContext(ctx, updateIncT)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&act, act); err != nil {
		return nil, fmt.Errorf("executing income tax update: %w", err)
	}

	return &act, nil
}

const updateIncTStatus = `
	UPDATE 
		income_tax 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateIncomeTaxStatus(ctx context.Context, act storage.IncomeTax) (*storage.IncomeTax, error) {
	logger.Info("update income tax status")
	stmt, err := s.db.PrepareNamedContext(ctx, updateIncTStatus)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&act, act); err != nil {
		return nil, fmt.Errorf("executing income tax status: %w", err)
	}

	return &act, nil
}
