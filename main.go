package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Хранилище сообщений (пока в памяти)
type NotificationRequest struct {
	Notifications struct {
		SMS   bool `json:"sms"`
		Email bool `json:"email"`
		InApp bool `json:"inApp"`
		Push  bool `json:"Push"`
	} `json:"notifications"`
	Message       string `json:"message"`
	IsDebug       bool   `json:"isDebug"`
	ScheduledTime string `json:"scheduledTime"`
	IsDelayed     bool   `json:"isDelayed"`
}

func main() {
	// Раздача статики
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Главная страница
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	// Обработчик POST /send (приём данных формы)
	http.HandleFunc("/api/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}

		// Читаем тело запроса
		var req NotificationRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			return
		}

		// Логируем полученные данные (для демонстрации)
		log.Printf("Получены данные: %+v\n", req)

		// Отправляем успешный ответ
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	})

	// Обработчик GET /api/messages (отдача данных в JSON)
	// http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != "GET" {
	// 		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
	// 		return
	// 	}

	// 	// Указываем, что ответ будет в JSON
	// 	w.Header().Set("Content-Type", "application/json")

	// 	// Преобразуем messages в JSON и отправляем
	// 	json.NewEncoder(w).Encode(map[string][]string{
	// 		"messages": messages,
	// 	})
	// })

	// Запуск сервера
	fmt.Println("Запуска сервера на localhost:8080")
	http.ListenAndServe(":8080", nil)
}
