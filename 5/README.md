# 5. Отправка уведомления при создании юзера через кафку

## Запуск

```bash
docker compose up -d

go run services/user/cmd/main.go

go run services/notification/cmd/main.go
```

**Интерфейсы:**
- Kafka UI: http://localhost:8080
- Mailpit: http://localhost:8025
