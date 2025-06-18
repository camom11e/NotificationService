package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Хранилище сообщений
type NotificationRequest struct {
	Notifications struct {
		SMS   bool `json:"sms"`
		Email bool `json:"email"`
		InApp bool `json:"inApp"`
		Push  bool `json:"Push"`
	} `json:"notifications"`
	Message  string `json:"message"`
	DateTime string `json:"dateTime"`
}

type NotificationSMS struct {
	SMS      bool
	Message  string
	DateTime time.Time
}

func (n NotificationSMS) GetType() string {
	return "NotificationSMS"
}

func (n NotificationSMS) GetMessage() string {
	return n.Message
}

func (n NotificationSMS) GetDateTime() time.Time {
	return n.DateTime
}

type NotificationEmail struct {
	Email    bool
	Message  string
	DateTime time.Time
}

func (n NotificationEmail) GetType() string {
	return "NotificationEmail"
}

func (n NotificationEmail) GetMessage() string {
	return n.Message
}

func (n NotificationEmail) GetDateTime() time.Time {
	return n.DateTime
}

type NotificationInApp struct {
	InApp    bool
	Message  string
	DateTime time.Time
}

func (n NotificationInApp) GetType() string {
	return "NotificationInApp"
}

func (n NotificationInApp) GetMessage() string {
	return n.Message
}

func (n NotificationInApp) GetDateTime() time.Time {
	return n.DateTime
}

type NotificationPush struct {
	Push     bool
	Message  string
	DateTime time.Time
}

func (n NotificationPush) GetType() string {
	return "NotificationPush"
}

func (n NotificationPush) GetMessage() string {
	return n.Message
}

func (n NotificationPush) GetDateTime() time.Time {
	return n.DateTime
}

// type QueueNotification struct {
// 	Notifications []interface{}
// 	mu            sync.Mutex
// }

// type Notification interface {
// }

func Convert(req NotificationRequest) (email *NotificationEmail, sms *NotificationSMS, push *NotificationPush, inApp *NotificationInApp) {
	var dateTime time.Time
	var err error
	if req.DateTime != "" {
		dateTime, err = time.Parse(time.RFC3339, req.DateTime)
		if err != nil {
			fmt.Printf("Ошибка парсинга даты: %v\n", err)
			return
		}

	} else {
		dateTime = time.Now()
	}
	if req.Notifications.Email {
		email = &NotificationEmail{
			Email:    true,
			Message:  req.Message,
			DateTime: dateTime,
		}
	}
	if req.Notifications.SMS {
		sms = &NotificationSMS{
			SMS:      true,
			Message:  req.Message,
			DateTime: dateTime,
		}
	}
	if req.Notifications.Push {
		push = &NotificationPush{
			Push:     true,
			Message:  req.Message,
			DateTime: dateTime,
		}
	}

	if req.Notifications.InApp {
		inApp = &NotificationInApp{
			InApp:    true,
			Message:  req.Message,
			DateTime: dateTime,
		}
	}
	return
}

func main() {
	// Раздача статики
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	// Обработчик
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

		// Логируем полученные данные
		log.Printf("Получены данные: %+v\n", req)

		// Отправляем успешный ответ
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	})

	// Запуск сервера
	fmt.Println("Запуска сервера на localhost:8080")
	http.ListenAndServe(":8080", nil)
}
