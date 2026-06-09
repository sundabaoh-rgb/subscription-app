# 📦 Subscription Service

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

Swagger UI доступен по адресу: `http://localhost:8080/swagger/`

### Эндпоинты

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/api/v1/subscriptions` | Создать подписку |
| `GET` | `/api/v1/subscriptions` | Получить список подписок |
| `GET` | `/api/v1/subscriptions/{id}` | Получить подписку по ID |
| `PUT` | `/api/v1/subscriptions/{id}` | Обновить подписку |
| `DELETE` | `/api/v1/subscriptions/{id}` | Удалить подписку |
| `GET` | `/api/v1/subscriptions/total-cost` | Подсчёт суммарной стоимости |

### Создание подписки

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

Ответ:
```json
{
  "ID": "b6620558-8de0-45af-8f0e-ca4e35f9daa0",
  "ServiceName": "Yandex Plus",
  "Price": 400,
  "UserID": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "StartDate": "2025-07-01T00:00:00Z",
  "EndDate": null,
  "CreatedAt": "2026-06-09T00:00:00Z"
}
```

### Список подписок с фильтрацией

```bash
GET /api/v1/subscriptions?user_id=60601fee-...&service_name=Yandex Plus&page=1&limit=20
```

### Подсчёт суммарной стоимости

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
    service_name TEXT        NOT NULL,
    price        INTEGER     NOT NULL CHECK (price > 0),
    user_id      UUID        NOT NULL,
    start_date   DATE        NOT NULL,
    end_date     DATE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Миграции

Миграции применяются автоматически при старте сервиса.

---

## Архитектура

Проект построен по принципу слоистой архитектуры:

```
HTTP Request
     ↓
  Handler        — парсинг запроса, валидация формата
     ↓
  Service        — бизнес логика, валидация правил
     ↓
  Repository     — SQL запросы к БД
     ↓
  PostgreSQL
```

Каждый слой зависит только от интерфейсов домена — не от конкретных реализаций. Это позволяет легко менять реализацию любого слоя.