package postgres

import (
	"context"
	"fmt"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
)

const resetpass = `
INSERT INTO password_change_requests(
	user_id,
	password_reset_token
) VALUES (
	:user_id,
	:password_reset_token
) RETURNING
	id
`

func (s *Storage) PassResRequest(con context.Context, prr storage.PassResRequest) (string, error) {
	logger.Info("create password reset request")
	stmt, err := s.db.PrepareNamed(resetpass)
	if err != nil {
		logger.Error(ewpq + err.Error())
		return "", err
	}

	var id string
	if err := stmt.Get(&id, prr); err != nil {
		logger.Error(ewmq + err.Error())
		return "", err
	}

	return id, nil
}

func (s *Storage) GetPassResRequestInfo(ctx context.Context, usrid, token string) (*storage.PassResRequest, error) {
	logger.Info("match password info while match request")
	stmt := "SELECT * from password_change_requests WHERE user_id=$1 AND password_reset_token=$2"
	var prer storage.PassResRequest
	err := s.db.Get(&prer, stmt, usrid, token)
	if err != nil {
		logger.Error("error while get password reset data. " + err.Error())
		return nil, fmt.Errorf("executing passeword reset table data: %w", err)
	}
	return &prer, nil
}

const saveresetpass = `
UPDATE 
	users 
SET
	password = :password
WHERE
	id = :user_id
RETURNING id
`

func (s *Storage) SavePasswordReset(ctx context.Context, stn storage.PassResRequest) (string, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, saveresetpass)
	if err != nil {
		return "", err
	}

	defer stmt.Close()
	if err := stmt.Get(&stn, stn); err != nil {
		return "", fmt.Errorf(" executing password reset: %w", err)
	}
	return stn.ID, nil
}
