package service

import (
	"BankAccount/model"
	"BankAccount/repository"
	"errors"
)

// Имплементация сервиса-интерфейса для операций над аккаунтами
type Account struct {
	repos repository.Account // Работает с репозиторием для операций над аккаунтами
}

func (acc *Account) CreateAccount() int {
	newAccount := model.NewAccount()

	acc.repos.CreateAccount(newAccount)

	return newAccount.ID
}

func (acc *Account) GetAccount(id int) (*model.Account, error) {
	account, err := acc.repos.GetAccountById(id)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (acc *Account) Deposit(id int, amount float64) error {
	if err := acc.repos.DepositMoneyById(id, amount); err != nil {
		return err
	}

	return nil
}

func (acc *Account) Withdraw(id int, amount float64) error {
	balance, err := acc.repos.GetBalanceById(id)
	if err != nil {
		return err
	}

	if balance-amount < 0 {
		return errors.New("not enough money")
	}

	if err := acc.repos.WithdrawMoneyById(id, amount); err != nil {
		return err
	}

	return nil
}

func (acc *Account) GetBalance(id int) (float64, error) {
	balance, err := acc.repos.GetBalanceById(id)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func NewAccountService(repos repository.Account) *Account {
	return &Account{
		repos: repos,
	}
}
