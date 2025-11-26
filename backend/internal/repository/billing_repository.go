package repository

import (
	"database/sql"

	"github.com/dejavu/backend/pkg/database"
)

type BillingRepository struct {
	db *database.DB
}

type BillingAccount struct {
	ID      string
	UserID  string
	Credits float64
}

type UsageRecord struct {
	ID           string
	UserID       string
	DeploymentID sql.NullString
	Type         string
	Amount       float64
	Description  string
}

func NewBillingRepository(db *database.DB) *BillingRepository {
	return &BillingRepository{db: db}
}

func (r *BillingRepository) GetOrCreateAccount(userID string) (*BillingAccount, error) {
	account := &BillingAccount{}
	
	// Try to get existing account
	err := r.db.QueryRow(
		"SELECT id, user_id, credits FROM billing_accounts WHERE user_id = $1",
		userID,
	).Scan(&account.ID, &account.UserID, &account.Credits)

	if err == sql.ErrNoRows {
		// Create new account with initial credits
		err = r.db.QueryRow(
			"INSERT INTO billing_accounts (user_id, credits) VALUES ($1, $2) RETURNING id",
			userID, 100.00, // $100 initial credits
		).Scan(&account.ID)
		if err != nil {
			return nil, err
		}
		account.UserID = userID
		account.Credits = 100.00
		return account, nil
	}

	return account, err
}

func (r *BillingRepository) AddCredits(userID string, amount float64) error {
	_, err := r.db.Exec(
		"UPDATE billing_accounts SET credits = credits + $1 WHERE user_id = $2",
		amount, userID,
	)
	return err
}

func (r *BillingRepository) DeductCredits(userID string, amount float64) error {
	result, err := r.db.Exec(
		"UPDATE billing_accounts SET credits = credits - $1 WHERE user_id = $2 AND credits >= $1",
		amount, userID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows // Insufficient credits
	}

	return nil
}

func (r *BillingRepository) RecordUsage(record *UsageRecord) error {
	query := `
		INSERT INTO usage_records (user_id, deployment_id, type, amount, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return r.db.QueryRow(
		query,
		record.UserID,
		record.DeploymentID,
		record.Type,
		record.Amount,
		record.Description,
	).Scan(&record.ID)
}

func (r *BillingRepository) GetUsageHistory(userID string, limit int) ([]*UsageRecord, error) {
	query := `
		SELECT id, user_id, deployment_id, type, amount, description
		FROM usage_records
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*UsageRecord
	for rows.Next() {
		record := &UsageRecord{}
		if err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.DeploymentID,
			&record.Type,
			&record.Amount,
			&record.Description,
		); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

