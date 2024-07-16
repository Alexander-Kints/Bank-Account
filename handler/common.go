package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// Универсальная функция для отправки json вида {"message": "text"}, устанавливает statusCode
func jsonResponse(w http.ResponseWriter, msg string, status int) {
	resp, err := json.Marshal(struct {
		Msg string `json:"message"`
	}{msg})

	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(status)

	if _, err := w.Write(resp); err != nil {
		log.Println(err)
		return
	}
}
