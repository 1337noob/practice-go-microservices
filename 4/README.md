# 4. WebSocket Chat

## Фичи

- Широковещательные сообщения (broadcast)
- Приватные сообщения `/w <username> <text>`
- Системные уведомления (join/leave)

## Запуск

```bash
go run cmd/main.go
```

Открыть в браузере: `http://localhost:8080?name=alice`

Открыть в браузере: `http://localhost:8080?name=bob`

Для приватного сообщения: `/w bob hello`
