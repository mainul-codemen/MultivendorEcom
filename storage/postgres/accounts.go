package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const insertact = `
INSERT INTO accounts(
	account_visualization, 
	account_name, 
	account_number,
	amount,
	status,
	created_by,
	updated_by
) VALUES (
	:account_visualization, 
	:account_name,
	:account_number, 
	:amount,
	:status,
	:created_by,
	:updated_by
) RETURNING
	id
`

func (s *Storage) CreateAccunt(con context.Context, act storage.Accounts) (string, error) {
	logger.Info("create accounts db")
	stmt, err := s.db.PrepareNamed(insertact)
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

func (s *Storage) GetAccounts(ctx context.Context, sts bool) ([]storage.Accounts, error) {
	logger.Info("get all accounts")
	actq := `SELECT * from accounts WHERE deleted_at IS NULL`
	if sts {
		actq = actq + " AND status=1"
	}
	act := make([]storage.Accounts, 0)
	if err := s.db.Select(&act, actq); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return act, nil
}

const delact = `
UPDATE
	accounts
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteAccounts(ctx context.Context, userid string, duid string) error {
	logger.Info("delete accounts")
	_, err := s.db.Exec(delact, duid, userid)
	if err != nil {
		logger.Error("delete accounts")
		return err
	}
	return nil
}

const gdactquery = `
SELECT
	id,
	account_visualization, 
	account_name, 
	account_number,
	amount,
	status,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM accounts
WHERE (id = $1 OR account_name = $1 OR account_number = $1) AND  deleted_at IS NULL
`

func (s *Storage) GetAccountsBy(ctx context.Context, idname string) (*storage.Accounts, error) {
	logger.Info("get accounts by")
	var dpt storage.Accounts
	if err := s.db.Get(&dpt, gdactquery, idname); err != nil {
		logger.Error("error while get accounts data. " + err.Error())
		return nil, fmt.Errorf("executing accounts details: %w", err)
	}
	return &dpt, nil
}

const updateAct = `
UPDATE accounts SET
	account_visualization = :account_visualization,
	account_name = :account_name,
	account_number = :account_number,
	amount = :amount,
	status = :status,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING *
`

func (s *Storage) UpdateAccounts(ctx context.Context, act storage.Accounts) (*storage.Accounts, error) {
	logger.Info("update accounts info")
	stmt, err := s.db.PrepareNamedContext(ctx, updateAct)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&act, act); err != nil {
		return nil, fmt.Errorf("executing accounts update: %w", err)
	}

	return &act, nil
}

const updateActStatus = `
	UPDATE 
		accounts 
	SET
		status = :status,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateAccountsStatus(ctx context.Context, act storage.Accounts) (*storage.Accounts, error) {
	logger.Info("update accounts status")
	stmt, err := s.db.PrepareNamedContext(ctx, updateActStatus)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&act, act); err != nil {
		return nil, fmt.Errorf("executing accounts status: %w", err)
	}

	return &act, nil
}

const addMony = `
	UPDATE 
		accounts 
	SET
		amount = :amount,
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) AddMoney(ctx context.Context, act storage.Accounts) (*storage.Accounts, error) {
	logger.Info("add money accounts")
	stmt, err := s.db.PrepareNamedContext(ctx, addMony)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&act, act); err != nil {
		return nil, fmt.Errorf("executing accounts add money: %w", err)
	}

	return &act, nil
}

const updateBlnc = `
	UPDATE 
		accounts 
	SET
		amount = :amount,
		updated_at = now(),
		updated_by = :updated_by
	WHERE 
		id = :id
	RETURNING *

`

func (s *Storage) UpdateBalance(ctx context.Context, actnt storage.Accounts) (*storage.Accounts, error) {
	logger.Info("update accounts balance")
	stmt, err := s.db.PrepareNamedContext(ctx, updateBlnc)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&actnt, actnt); err != nil {
		return nil, fmt.Errorf("executing accounts upate balance: %w", err)
	}

	return &actnt, nil
}
