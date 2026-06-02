# 6. Метрики, Логирование

### filebeat нормально запустится только на linux !!!

## Запуск

```bash
docker compose up --build -d

## Зайти в Grafana
## Management -> Stack Management -> Data Views
## Добавить новый с паттерном go-app-filebeat-*
```

**Доступ к сервисам:**
- Grafana: http://localhost:3000
- Kibana: http://localhost:5601
