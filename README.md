# Document Review System

Учебная система проверки документов на Go: регистрация, логин, отправка документов,
просмотр статуса проверки, PostgreSQL, RabbitMQ, Prometheus и Grafana.

## Сценарий

1. Пользователь регистрируется.
2. Пользователь входит в систему и получает JWT.
3. Пользователь отправляет документ на проверку.
4. Документ сохраняется в PostgreSQL со статусом `pending_review`.
5. Сервис публикует событие `document-submitted` в RabbitMQ.
6. Worker читает событие из очереди и меняет статус документа:
   - `pending_review`
   - `in_review`
   - `approved`
7. Пользователь обновляет список документов и видит актуальный статус проверки.

## Архитектура

- `internal/transport/httpapi` - HTTP API, HTML-страница, JWT middleware, Prometheus metrics.
- `internal/usecase` - бизнес-логика регистрации, логина и проверки документов.
- `internal/repository/postgres` - SQL и работа с PostgreSQL.
- `internal/queue/rabbitmq` - RabbitMQ publisher/consumer для событий документов.
- `internal/worker` - worker проверки документов.
- `internal/config` - настройки из `.env`.

## Стек

- Go HTTP server
- PostgreSQL
- JWT + bcrypt
- RabbitMQ
- Prometheus
- Grafana
- Docker Compose

## Запуск

1. Заполните `.env` по примеру `.env.example`.

2. Поднимите инфраструктуру:

```powershell
docker compose up -d
```

3. Запустите Go-сервис:

```powershell
go run .
```

## Адреса

- Сайт/API: `http://localhost:8080`
- Метрики Go-сервиса: `http://localhost:8080/metrics`
- RabbitMQ Management: `http://localhost:15672` (`admin` / `admin`)
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000` (`admin` / `admin`)

## API

Публичные:

- `GET /` - HTML-интерфейс системы.
- `GET /test` - проверка сервера, возвращает `Hello!`.
- `POST /dbtest` - тестовая запись строки в PostgreSQL.
- `POST /auth/register` - регистрация.
- `POST /auth/login` - логин и получение JWT.

Защищенные, нужен `Authorization: Bearer <JWT>`:

- `POST /documents` - отправить документ на проверку.
- `GET /documents` - получить свои документы и статусы проверки.

## RabbitMQ

Используются:

- Exchange: `main_exchange`
- Queue: `document_submitted`
- Queue: `document_status`
- Routing key: `document-submitted`
- Routing key: `document-status-change`

В RabbitMQ можно смотреть, как событие о новом документе попадает в очередь.
Если worker быстро обработал сообщение, очередь может быть пустой, но в RabbitMQ будут видны
подключения, каналы и активность по очередям.

## Grafana и Prometheus

Prometheus собирает метрики с `/metrics`.

Полезные PromQL-запросы:

```promql
hse_http_requests_total
```

```promql
hse_http_requests_total{path="/documents"}
```

```promql
rate(hse_http_requests_total[1m])
```

```promql
rate(hse_http_request_duration_seconds_sum[1m])
/
rate(hse_http_request_duration_seconds_count[1m])
```

## Тесты

```powershell
go test ./...
```

Проверяется:

- повторная регистрация email запрещена;
- email нормализуется;
- неверный пароль не проходит;
- успешный логин возвращает JWT;
- защищенные действия требуют пользователя.
