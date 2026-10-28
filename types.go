package main

import "math/rand/v2"

type Account struct {
	ID            int    `json:"id"`
	Firstname     string `json:"firstName"`
	LastName      string `json:"lastName"`
	AccountNumber int64  `json:"accountNumber"`
	Balance       int64  `json:"balance"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:            rand.IntN(10000),
		Firstname:     firstName,
		LastName:      lastName,
		AccountNumber: int64(rand.IntN(1000000)),
	}
}
