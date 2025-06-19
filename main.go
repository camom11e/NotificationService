package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"NotificationService/rabbitmq" 
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
	SMS      bool      `json:"sms"`
	Message  string    `json:"message"`
	DateTime time.Time `json:"date_time"`
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
	Email    bool      `json:"email"`
	Message  string    `json:"message"`
	DateTime time.Time `json:"date_time"`
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
	InApp    bool      `json:"inApp"`
	Message  string    `json:"message"`
	DateTime time.Time `json:"date_time"`
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
	Push     bool      `json:"push"`
	Message  string    `json:"message"`
	DateTime time.Time `json:"date_time"`
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

type Notification interface {
	GetType() string
	GetMessage() string
	GetDateTime() time.Time
}

// Инициализация шедулера
var DeferredMessages []Notification

func Add(n Notification) {
	DeferredMessages = append(DeferredMessages, n)
}

func processNotification(n Notification) {
	// Здесь будет логика обработки уведомления
	fmt.Printf("Processing notification: %s\n", n.GetType())
}

func serializeNotificationToJson(n Notification) ([]byte, error) {
	return json.Marshal(n)
}

func checkReady(n Notification) bool {
	return n.GetDateTime().Before(time.Now()) || n.GetDateTime().Equal(time.Now())
}

// Запуск шедулера
func Start() {
	// Канал для очереди уведомлений
	queue := make(chan Notification, 100)
	
	// Запускаем worker'ов для обработки уведомлений
	for i := 0; i < 5; i++ {
		go func() {
			for n := range queue {
				jsonData, err := serializeNotificationToJson(n)
				if err != nil {
					log.Printf("Сериализация неудалась: %v", err)
					continue
				}
				
				if err := client.Publish(jsonData); err != nil {
					log.Printf("Publish error: %v", err)
					continue
				}
				
				fmt.Printf("Sent: %s\n", n.GetMessage())
				time.Sleep(1 * time.Second)
			}
		}()
	}
	
	// Шедулер, который проверяет уведомления
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			var remaining []Notification
			for _, n := range DeferredMessages {
				if checkReady(n) {
					queue <- n
				} else {
					remaining = append(remaining, n)
				}
			}
			DeferredMessages = remaining
		}
	}
}

func Convert(req NotificationRequest) (email *NotificationEmail, sms *NotificationSMS, push *NotificationPush, inApp *NotificationInApp) {
	var dateTime time.Time
	var err error
	
	if req.DateTime != "" {
		formats := []string{
			time.RFC3339,                    
			"2006-01-02T15:04:05",           
			"2006-01-02T15:04",               
			"2006-01-02 15:04:05",           
			"2006-01-02 15:04",               
		}
		for _, format := range formats {
		dateTime, err = time.Parse(format, req.DateTime)
		if err == nil {
			break
		
		}
	}

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

var client *rabbitmq.Client

func main() {
	// Инициализация клиента RabbitMQ
	var err error
	client, err = rabbitmq.New(rabbitmq.Config{
		URL:       "amqp://guest:guest@localhost:5672/",
		QueueName: "test_queue",
	})
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Error closing RabbitMQ connection: %v", err)
		}
	}()

	// Раздача статики
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	// Обработчик API
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

		email, sms, push, inApp := Convert(req)
		if email != nil {
			Add(email)
		}
		if sms != nil {
			Add(sms)
		}
		if push != nil {
			Add(push)
		}
		if inApp != nil {
			Add(inApp)
		}

		// Запуск шедулера в отдельной горутине
		go Start()

		// Отправляем успешный ответ
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	})

	// Запуск сервера
	fmt.Println("Запуск сервера на localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}