package handler

import "BankAccount/service"

// Главный обработчик. Работает с сервисами
type Handler struct {
	services *service.Service
}

// Конструктор главного обработчика
func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}
