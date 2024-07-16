package handler

import (
	"BankAccount/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Обработчик для создания аккаунта
func (h *Handler) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		id := h.services.BankAccount.CreateAccount()
		jsonResponse(w, fmt.Sprintf("account with id %d was created", id), http.StatusOK)

		log.Printf("account with id %d was created", id)
	}(&wg)

	wg.Wait()
}

// Обработчик для пополнения аккаунта
func (h *Handler) DepositHandler(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Запуск горутины с wg
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		// Тело POST-запроса. Параметры: amount
		var operationBody model.OperationBody

		// Валидация slug id
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			jsonResponse(w, "no valid id", http.StatusBadRequest)
			return
		}

		// Распаршивание тела запроса
		if err := json.NewDecoder(r.Body).Decode(&operationBody); err != nil {
			jsonResponse(w, "server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		// Вызов сервиса, ответственного за пополнение
		if err := h.services.BankAccount.Deposit(id, operationBody.Amount); err != nil {
			jsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Получение аккаунта из хранилища
		account, _ := h.services.BankAccount.GetAccount(id)

		// Если аккаунт в процессе снятия, отправляем в канал аккаунта сумму пополнения
		account.Lock()
		if account.InProcess {
			account.Ch <- operationBody.Amount
		}
		account.Unlock()

		jsonResponse(w, fmt.Sprintf("account id %d deposit amount %.2f", id, operationBody.Amount), http.StatusOK)
		log.Printf("account id %d deposit amount %.2f", id, operationBody.Amount)
	}(&wg)

	wg.Wait()
}

// Обработчик для снятия денег с аккаунта
func (h *Handler) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		// Тело POST-запроса. Параметры: amount
		var operationBody model.OperationBody

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			jsonResponse(w, "no valid id", http.StatusBadRequest)
			return
		}

		account, err := h.services.BankAccount.GetAccount(id)
		if err != nil {
			jsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Если аккаунт уже в процессе снятия, делаем возврат соотв. сообщения
		if account.InProcess {
			jsonResponse(w, "account already in process", http.StatusOK)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&operationBody); err != nil {
			jsonResponse(w, "server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		// Функция для отправки сообщения об успешном снятии денег
		successFunc := func() {
			jsonResponse(w, fmt.Sprintf("account id %d withdraw amount %.2f", id, operationBody.Amount),
				http.StatusOK)
			log.Printf("account id %d withdraw amount %.2f", id, operationBody.Amount)
		}

		// Если не хватает денег на балансе, программа будет ждать 20 секунд, пока на балансе не наберется
		// нужная сумма. В противном случае возврат сообщения о недостатке средств.
		if err := h.services.BankAccount.Withdraw(id, operationBody.Amount); err != nil {
			account.Lock()
			account.InProcess = true
			account.Unlock()
			ctx, _ := context.WithTimeout(account.Ctx, 20*time.Second)
			for {
				select {
				case <-ctx.Done():
					jsonResponse(w, err.Error(), http.StatusBadRequest)
					account.Lock()
					account.InProcess = false
					account.Unlock()
					return
				case _, ok := <-account.Ch:
					if ok {
						fmt.Println("пополнение")
						if err := h.services.BankAccount.Withdraw(id, operationBody.Amount); err == nil {
							account.Lock()
							account.InProcess = false
							account.Unlock()
							successFunc()
							return
						}
					}
				default:
					time.Sleep(2 * time.Second)
					log.Printf("withdraw accout id %d:waiting money...", id)
				}
			}
		}

		successFunc()
	}(&wg)

	wg.Wait()
}

// Обработчик для проверки баланса.
func (h *Handler) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			jsonResponse(w, "no valid id", http.StatusBadRequest)
			return
		}

		balance, err := h.services.BankAccount.GetBalance(id)
		if err != nil {
			jsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		balanceBody := struct {
			ID      int    `json:"id"`
			Balance string `json:"balance"`
		}{id, fmt.Sprintf("%.2f", balance)}

		balanceResponse, err := json.Marshal(&balanceBody)
		if err != nil {
			jsonResponse(w, "server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if _, err := w.Write(balanceResponse); err != nil {
			jsonResponse(w, "server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		log.Printf("balance request account id %d (balance: %.2f)", id, balance)
	}(&wg)

	wg.Wait()
}
