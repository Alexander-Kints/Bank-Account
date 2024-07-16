package repository

import (
	"BankAccount/model"
	"BankAccount/storage"
	"errors"
)

// Имплементация репозитория для операций над аккантами в локальном хранилище (AccountDB)
type AccountRepository struct {
	db *storage.AccountDB
}

func (ar *AccountRepository) CreateAccount(account *model.Account) {
	ar.db.Lock()
	defer ar.db.Unlock()

	ar.db.DB[account.ID] = account
}

func (ar *AccountRepository) GetAccountById(id int) (*model.Account, error) {
	ar.db.RLock()
	defer ar.db.RUnlock()

	account, ok := ar.db.DB[id]
	if !ok {
		return nil, errors.New("account not exist")
	}

	return account, nil
}

func (ar *AccountRepository) DepositMoneyById(id int, amount float64) error {
	ar.db.RLock()
	defer ar.db.RUnlock()

	account, ok := ar.db.DB[id]
	if !ok {
		return errors.New("account not exist")
	}

	account.Lock()
	defer account.Unlock()

	account.Balance += amount

	return nil
}

func (ar *AccountRepository) WithdrawMoneyById(id int, amount float64) error {
	ar.db.RLock()
	defer ar.db.RUnlock()

	account, ok := ar.db.DB[id]
	if !ok {
		return errors.New("account not exist")
	}

	account.Lock()
	defer account.Unlock()

	account.Balance -= amount

	return nil
}

func (ar *AccountRepository) GetBalanceById(id int) (float64, error) {
	ar.db.RLock()

	account, ok := ar.db.DB[id]
	if !ok {
		ar.db.RUnlock()
		return 0, errors.New("account not exist")
	}
	ar.db.RUnlock()

	account.Lock()
	defer account.Unlock()

	balance := account.Balance

	return balance, nil
}

func NewAccountRepository(db *storage.AccountDB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}
