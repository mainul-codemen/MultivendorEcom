package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const itrnsstyp = `
INSERT INTO transaction_source(
	transaction_source_name, 
	status,
	created_by,
	updated_by
) VALUES (
	:transaction_source_name, 
	:status, 
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) CreateTransactionSource(con context.Context, des storage.TransactionSource) (string, error) {
	logger.Info("create transaction_source db")
	stmt, err := s.db.PrepareNamed(itrnsstyp)
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

func (s *Storage) GetTransactionSource(ctx context.Context, sts bool) ([]storage.TransactionSource, error) {
	logger.Info("get all transaction_source")
	desq := `SELECT * from transaction_source WHERE deleted_at IS NULL`
	if sts {
		desq = desq + " AND status=1"
	}
	des := make([]storage.TransactionSource, 0)
	if err := s.db.Select(&des, desq); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return des, nil
}

const deltrts = `
UPDATE
	transaction_source
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteTransactionSource(ctx context.Context, userid string, duid string) error {
	logger.Info("delete transaction_source")
	_, err := s.db.Exec(deltrts, duid, userid)
	if err != nil {
		logger.Error("delete transaction_source")
		return err
	}
	return nil
}

const gttsq = `
SELECT
	id,
	transaction_source_name, 
	status,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM transaction_source
WHERE (id = $1 OR transaction_source_name = $1) AND  deleted_at IS NULL
`

func (s *Storage) GetTransactionSourceBy(ctx context.Context, idname string) (*storage.TransactionSource, error) {
	var dpt storage.TransactionSource
	if err := s.db.Get(&dpt, gttsq, idname); err != nil {
		logger.Error("error while get transaction_source data. " + err.Error())
		return nil, fmt.Errorf("executing transaction_source details: %w", err)
	}
	return &dpt, nil
}

const uptrs = `
UPDATE transaction_source SET
	transaction_source_name = :transaction_source_name,
	status = :status,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING *
`

func (s *Storage) UpdateTransactionSource(ctx context.Context, p storage.TransactionSource) (*storage.TransactionSource, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, uptrs)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing transaction_source update: %w", err)
	}
	return &p, nil
}

const updatetrst = `
	UPDATE 
		transaction_source 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateTransactionSourceStatus(ctx context.Context, p storage.TransactionSource) (*storage.TransactionSource, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updatetrst)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing transaction_source status: %w", err)
	}
	return &p, nil
}
