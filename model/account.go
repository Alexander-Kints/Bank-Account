package model

import (
	"context"
	"sync"
)

var accountCounter = 0

// Модель аккаунта
type Account struct {
	sync.Mutex
	ID        int
	Balance   float64
	InProcess bool
	Ctx       context.Context
	Ch        chan float64
}

// Конструктор аккаунта
func NewAccount() *Account {
	accountCounter++

	return &Account{
		Mutex:     sync.Mutex{},
		ID:        accountCounter,
		Balance:   0,
		InProcess: false,
		Ctx:       context.Background(),
		Ch:        make(chan float64, 10),
	}
}
