package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginDTO struct {
	AccountNumber int64  `json:"account_number"`
	Password      string `json:"password"`
}

type TransferDTO struct {
	ToAccount int     `json:"to_account"`
	Amount    float64 `json:"amount"`
}

type CreateAccountDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type Account struct {
	ID            int       `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	AccountNumber int64     `json:"account_number"`
	PasswordHash  string    `json:"-"`
	Balance       int64     `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewAccount creates an Account with a bcrypt-hashed password.
// The account number is assigned later by the storage layer via DB sequence + Feistel cipher.
func NewAccount(firstName, lastName, password string) (*Account, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: string(hash),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}, nil
}

// VerifyPassword checks if the given plaintext password matches the stored hash.
func (a *Account) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password)) == nil
}
