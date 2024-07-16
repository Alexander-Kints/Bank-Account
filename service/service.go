package service

import (
	"BankAccount/model"
	"BankAccount/repository"
)

// Сервис-интерфейс для операций над аккаунтами
type BankAccount interface {
	CreateAccount() int
	GetAccount(id int) (*model.Account, error)
	Deposit(id int, amount float64) error
	Withdraw(id int, amount float64) error
	GetBalance(id int) (float64, error)
}

// Главный сервис
type Service struct {
	BankAccount
}

// Конструктор главного сервиса
func NewService(repos *repository.Repository) *Service {
	return &Service{
		BankAccount: NewAccountService(repos.Account),
	}
}
