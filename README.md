# Subscription Service

REST-сервис для агрегации данных об онлайн подписках пользователей.

---

## Быстрый старт

### Требования
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### Запуск

```bash
git clone https://github.com/your-username/subscription-service
cd subscription-service
cp .env.example .env
docker compose up --build
```

Сервис будет доступен на `http://localhost:8080`

Swagger UI: `http://localhost:8080/swagger/`

---

## Стек технологий

| Слой | Технология |
|------|-----------|
| Язык | Go 1.25 |
| База данных | PostgreSQL 16 |
| Драйвер БД | pgx v5 |
| Миграции | golang-migrate |
| Логгер | Uber Zap |
| Контейнеризация | Docker + Docker Compose |

---


## Конфигурация

Создай `.env` файл на основе `.env.example`:

```env
APP_PORT=8080
APP_ENV=local
LOG_LEVEL=info

DB_HOST=postgres
DB_PORT=5432
DB_NAME=subscriptions
DB_USER=postgres
DB_PASSWORD=postgres
```

---

## API

### Эндпоинты

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/api/v1/subscriptions` | Создать подписку |
| `GET` | `/api/v1/subscriptions` | Получить список подписок |
| `GET` | `/api/v1/subscriptions/{id}` | Получить подписку по ID |
| `PUT` | `/api/v1/subscriptions/{id}` | Обновить подписку |
| `DELETE` | `/api/v1/subscriptions/{id}` | Удалить подписку |
| `GET` | `/api/v1/subscriptions/total-cost` | Подсчёт суммарной стоимости |

### Примеры запросов

**Создание подписки**
```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

**Список подписок с фильтрацией и пагинацией**
```bash
GET /api/v1/subscriptions?user_id=60601fee-...&service_name=Yandex Plus&page=1&limit=20
```

**Подсчёт суммарной стоимости за период**
```bash
GET /api/v1/subscriptions/total-cost?user_id=60601fee-...&from=01-2025&to=12-2025
```

Ответ:
```json
{
  "total_cost": 2400,
  "currency": "RUB",
  "period": {
    "from": "01-2025",
    "to": "12-2025"
  }
}
```

---

## База данных

### Схема таблицы

```sql
CREATE TABLE subscriptions (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name TEXT        NOT NULL CHECK (length(trim(service_name)) > 0),
    price        INTEGER     NOT NULL CHECK (price > 0),
    user_id      UUID        NOT NULL,
    start_date   DATE        NOT NULL,
    end_date     DATE        CHECK (end_date IS NULL OR end_date > start_date),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Миграции применяются автоматически при старте сервиса.

---

## Тесты

```bash
go test ./...
```

Покрыты юнит тестами:
- Успешное создание подписки
- Валидация: цена <= 0
- Валидация: пустое название сервиса
- Валидация: дата окончания раньше даты начала
- GetByID: запись не найдена
- Delete: успешное удаление
- Delete: запись не найдена