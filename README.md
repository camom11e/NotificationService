# Notification Service

Микросервис для обработки и рассылки уведомлений через RabbitMQ

## 🚀 Запуск сервиса

- Docker
- Go 1.24.1

### 1. Запуск RabbitMQ
```bash
docker run -d --name rabbitmq \
  -p 5672:5672 -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=guest \
  -e RABBITMQ_DEFAULT_PASS=guest \
  rabbitmq:3-management
```


```bash
git clone https://github.com/camom11e/NotificationService
cd NotificationService
```
```bash
go run main.go
```