package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const itrnstyp = `
INSERT INTO transaction_types(
	transaction_type_name, 
	status,
	created_by,
	updated_by
) VALUES (
	:transaction_type_name, 
	:status, 
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) CreateTransactionTypes(con context.Context, des storage.TransactionTypes) (string, error) {
	logger.Info("create transaction_types db")
	stmt, err := s.db.PrepareNamed(itrnstyp)
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

func (s *Storage) GetTransactionTypes(ctx context.Context, sts bool) ([]storage.TransactionTypes, error) {
	logger.Info("get all transaction_types")
	desq := `SELECT * from transaction_types WHERE deleted_at IS NULL`
	if sts {
		desq = desq + " AND status=1"
	}
	des := make([]storage.TransactionTypes, 0)
	if err := s.db.Select(&des, desq); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return des, nil
}

const deltrt = `
UPDATE
	transaction_types
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteTransactionTypes(ctx context.Context, userid string, duid string) error {
	logger.Info("delete transaction_types")
	_, err := s.db.Exec(deltrt, duid, userid)
	if err != nil {
		logger.Error("delete transaction_types")
		return err
	}
	return nil
}

const gttq = `
SELECT
	id,
	transaction_type_name, 
	status,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM transaction_types
WHERE (id = $1 OR transaction_type_name = $1) AND  deleted_at IS NULL
`

func (s *Storage) GetTransactionTypesBy(ctx context.Context, idname string) (*storage.TransactionTypes, error) {
	var dpt storage.TransactionTypes
	if err := s.db.Get(&dpt, gttq, idname); err != nil {
		logger.Error("error while get transaction_types data. " + err.Error())
		return nil, fmt.Errorf("executing transaction_types details: %w", err)
	}
	return &dpt, nil
}

const uptr = `
UPDATE transaction_types SET
	transaction_type_name = :transaction_type_name,
	status = :status,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING *
`

func (s *Storage) UpdateTransactionTypes(ctx context.Context, p storage.TransactionTypes) (*storage.TransactionTypes, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, uptr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing transaction_types update: %w", err)
	}
	return &p, nil
}

const updatetrsts = `
	UPDATE 
		transaction_types 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateTransactionTypesStatus(ctx context.Context, p storage.TransactionTypes) (*storage.TransactionTypes, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updatetrsts)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing transaction_types status: %w", err)
	}
	return &p, nil
}
