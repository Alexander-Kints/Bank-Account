package model

// Модель для тела json (снятие и пополнение аккаунта)
type OperationBody struct {
	ID     int     `json:"id"`
	Amount float64 `json:"amount"`
}
