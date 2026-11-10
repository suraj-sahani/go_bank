package main

import (
	"time"
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
}

type Account struct {
	ID            int       `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	AccountNumber int64     `json:"account_number"`
	Balance       int64     `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}
