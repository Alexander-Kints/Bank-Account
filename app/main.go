package main

import (
	"BankAccount/config"
	"BankAccount/handler"
	"BankAccount/repository"
	"BankAccount/service"
	"BankAccount/storage"
	"log"
	"net/http"
)

func main() {
	cfg := config.Config{
		Host: "localhost",
		Port: "9000",
	}
	db := storage.NewDB()
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	// Инициализация роутов
	http.HandleFunc("/accounts", handlers.CreateAccountHandler)
	http.HandleFunc("/accounts/{id}/deposit", handlers.DepositHandler)
	http.HandleFunc("/accounts/{id}/withdraw", handlers.WithdrawHandler)
	http.HandleFunc("/accounts/{id}/balance", handlers.BalanceHandler)

	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
