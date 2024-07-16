package repository

import (
	"BankAccount/model"
	"BankAccount/storage"
)

// Репозиторий для операций над аккаунтами (интерфейс)
type Account interface {
	CreateAccount(account *model.Account)
	GetAccountById(id int) (*model.Account, error)
	DepositMoneyById(id int, amount float64) error
	WithdrawMoneyById(id int, amount float64) error
	GetBalanceById(id int) (float64, error)
}

// Главный репозиторий
type Repository struct {
	Account
}

// Конструктор главного репозитория
func NewRepository(db *storage.AccountDB) *Repository {
	return &Repository{
		Account: NewAccountRepository(db),
	}
}
