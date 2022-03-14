package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const insertactrn = `
INSERT INTO accounts_transactions(
	from_account_id, 
	to_account_id, 
	user_id,
	transaction_amount,
	transaction_type_id,
	transaction_source_id,
	reference,
	note,
	status,
	created_by,
	updated_by
) VALUES (
	:from_account_id, 
	:to_account_id, 
	:user_id,
	:transaction_amount,
	:transaction_type_id,
	:transaction_source_id,
	:reference,
	:note,
	:status,
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) CreateAccountsTransaction(con context.Context, act storage.AccountsTransaction) (string, error) {
	logger.Info("create accounts_transactions db")
	stmt, err := s.db.PrepareNamed(insertactrn)
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

func (s *Storage) GetAccountsTransaction(ctx context.Context, sts bool) ([]storage.AccountsTransaction, error) {
	logger.Info("get all accounts_transactions")
	actq := `
	SELECT
	at.id,
	from_account_id,
	a.account_name AS from_account_name,
	ab.account_name AS to_account_name,
	ts.transaction_source_name,
	tt.transaction_type_name,
	to_account_id, 
	user_id,
	transaction_amount,
	transaction_type_id,
	transaction_source_id,
	reference,
	note,
	at.status,
	at.created_at,
	at.created_by,
	at.updated_at,
	at.updated_by
FROM 
	accounts_transactions as at
	LEFT JOIN transaction_source AS ts ON ts.id = transaction_source_id 
	LEFT JOIN transaction_types AS tt ON tt.id = transaction_type_id 
	LEFT JOIN accounts AS a ON a.id = from_account_id
	LEFT JOIN accounts AS ab ON ab.id = to_account_id
WHERE 
	at.deleted_at IS NULL`
	if sts {
		actq = actq + " AND status=1"
	}
	act := make([]storage.AccountsTransaction, 0)
	if err := s.db.Select(&act, actq); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return act, nil
}

const delactrns = `
UPDATE
	accounts_transactions
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteAccountsTransaction(ctx context.Context, userid string, duid string) error {
	logger.Info("delete accounts_transactions")
	_, err := s.db.Exec(delactrns, duid, userid)
	if err != nil {
		logger.Error("delete accounts_transactions")
		return err
	}
	return nil
}

const getacttrn = `
SELECT
	at.id,
	from_account_id,
	a.account_name AS from_account_name,
	to_account_id, 
	ab.account_name AS to_account_name,
	user_id,
	ts.transaction_source_name,
	tt.transaction_type_name,
	transaction_amount,
	transaction_type_id,
	transaction_source_id,
	reference,
	note,
	at.status,
	at.created_at,
	at.created_by,
	at.updated_at,
	at.updated_by
FROM 
accounts_transactions as at
	LEFT JOIN transaction_source AS ts ON ts.id = transaction_source_id 
	LEFT JOIN transaction_types AS tt ON tt.id = transaction_type_id 
	LEFT JOIN accounts AS a ON a.id = from_account_id
	LEFT JOIN accounts AS ab ON ab.id = to_account_id
WHERE (at.id = $1 OR user_id = $1) AND  at.deleted_at IS NULL
`

func (s *Storage) GetAccountsTransactionBy(ctx context.Context, idname string) (*storage.AccountsTransaction, error) {
	logger.Info("get accounts_transactions by")
	var dpt storage.AccountsTransaction
	if err := s.db.Get(&dpt, getacttrn, idname); err != nil {
		logger.Error("error while get accounts_transactions data. " + err.Error())
		return nil, fmt.Errorf("executing accounts_transactions details: %w", err)
	}
	return &dpt, nil
}

const updateActtrns = `
UPDATE accounts_transactions SET
	from_account_id = :from_account_id,
	to_account_id = :to_account_id,
	user_id = :user_id,
	transaction_amount = :transaction_amount,
	transaction_type_id = :transaction_type_id,
	reference = :reference,
	note = :note,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING *
`

func (s *Storage) UpdateAccountsTransaction(ctx context.Context, act storage.AccountsTransaction) (*storage.AccountsTransaction, error) {
	logger.Info("update accounts_transactions info")
	stmt, err := s.db.PrepareNamedContext(ctx, updateActtrns)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&act, act); err != nil {
		return nil, fmt.Errorf("executing accounts_transactions update: %w", err)
	}

	return &act, nil
}

const updateActtrStatus = `
	UPDATE 
		accounts_transactions 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateAccountsTransactionStatus(ctx context.Context, act storage.AccountsTransaction) (*storage.AccountsTransaction, error) {
	logger.Info("update accounts_transactions status")
	stmt, err := s.db.PrepareNamedContext(ctx, updateActtrStatus)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&act, act); err != nil {
		return nil, fmt.Errorf("executing accounts_transactions status: %w", err)
	}

	return &act, nil
}
