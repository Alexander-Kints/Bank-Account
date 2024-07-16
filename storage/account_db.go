package storage

import (
	"BankAccount/model"
	"sync"
)

// Локальное хранилище аккаунтов
type AccountDB struct {
	sync.RWMutex // Потокобезопасность
	DB           map[int]*model.Account
}

func NewDB() *AccountDB {
	return &AccountDB{
		RWMutex: sync.RWMutex{},
		DB:      make(map[int]*model.Account),
	}
}
