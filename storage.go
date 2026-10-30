package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// ErrAccountNotFound is returned when an account lookup yields no results.
var ErrAccountNotFound = errors.New("account not found")

// AccountStore defines the interface for account persistence operations.
type AccountStore interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	DeleteAccount(int) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
}

// PostgresAccountStore implements AccountStore using a PostgreSQL database.
type PostgresAccountStore struct {
	db *sql.DB
}

// NewPostgresAccountStore connects to PostgreSQL using the DATABASE_URL
// environment variable and returns a ready-to-use store.
func NewPostgresAccountStore() (*PostgresAccountStore, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresAccountStore{
		db: db,
	}, nil
}

// Init creates the required database tables if they don't already exist.
func (s *PostgresAccountStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresAccountStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
			id             SERIAL PRIMARY KEY,
			first_name     VARCHAR(50) NOT NULL,
			last_name      VARCHAR(50) NOT NULL,
			account_number BIGINT NOT NULL UNIQUE,
			balance        BIGINT NOT NULL DEFAULT 0,
			created_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresAccountStore) CreateAccount(a *Account) error {
	query := `
			INSERT INTO accounts (
				first_name,
				last_name,
				account_number,
				balance
			)
			VALUES ($1, $2, $3, $4)`

	_, err := s.db.Exec(
		query,
		a.FirstName,
		a.LastName,
		a.AccountNumber,
		a.Balance,
	)

	return err
}

func (s *PostgresAccountStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresAccountStore) DeleteAccount(id int) error {
	query := `DELETE FROM accounts WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrAccountNotFound
	}

	return nil
}

func (s *PostgresAccountStore) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM accounts`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*Account, 0)
	for rows.Next() {
		acc, err := scanRowToAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (s *PostgresAccountStore) GetAccountByID(id int) (*Account, error) {
	query := `SELECT * FROM accounts WHERE id = $1`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return scanRowToAccount(rows)
	}
	return nil, ErrAccountNotFound
}

func scanRowToAccount(rows *sql.Rows) (*Account, error) {
	acc := new(Account)

	err := rows.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.AccountNumber,
		&acc.Balance,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)

	return acc, err
}
