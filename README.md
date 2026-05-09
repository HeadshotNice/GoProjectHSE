# HSE Labs (LR1–LR5) — Go HTTP Server

Проект: простой HTTP сервер на Go с чистой архитектурой (3 слоя), PostgreSQL, JWT и middleware.

## Быстрый старт

1. Заполните переменные окружения (есть пример в [.env](./.env)).
2. Запуск:
   - из консоли: `go run .`
   - из GoLand: Run конфигурация для `main` (важно: Working Directory = корень проекта, чтобы `.env` подхватился)

Сервер по умолчанию слушает `:8080` (переменная `HTTP_ADDR`).

## Эндпоинты

- `GET /test` -> `Hello!`
- `POST /dbtest` -> пишет строку тела запроса в БД, ответ `{ "ok": true }`
- `POST /auth/register` -> регистрация пользователя
- `POST /auth/login` -> логин, возвращает JWT `{ "token": "..." }`
- `POST /orders` -> создать заказ (нужен `Authorization: Bearer <JWT>`)
- `GET /orders` -> список заказов пользователя (нужен `Authorization: Bearer <JWT>`)

## Чистая архитектура (3 слоя)

**Transport (HTTP)**
- Роутинг/парсинг запросов/формирование ответов, middleware.

**Usecase (Business logic)**
- Валидация данных, сценарии: тестовый hello, запись в БД, регистрация/логин, заказы, JWT issuance.

**Repository (DB isolation)**
- Только SQL и работа с БД (PostgreSQL).

Связь слоев сделана через интерфейсы, которые объявлены на уровне usecase:
см. [internal/usecase/usecase.go](./internal/usecase/usecase.go) (интерфейсы `TestRepo`, `DBTestRepo`, `UsersRepo`, `OrdersRepo`).

## ЛР1 — Пустой HTTP сервер + /test + graceful shutdown + 3 слоя

Что требуется:
- HTTP сервер на Go
- `GET /test` -> `"Hello!"`
- обработка запроса задействует все 3 слоя (transport → usecase → repository)
- graceful shutdown по сигналу завершения

Где реализовано:
- Точка входа и graceful shutdown: [main.go](./main.go)
  - `signal.NotifyContext(...)` ловит `SIGTERM`/`Ctrl+C`
  - `srv.Shutdown(ctx)` завершает сервер корректно
- HTTP слой и `GET /test`: [internal/transport/httpapi/handler.go](./internal/transport/httpapi/handler.go)
  - `handleTest` вызывает `uc.TestHello(...)`
- Usecase слой для теста: [internal/usecase/usecase.go](./internal/usecase/usecase.go)
  - `TestHello` вызывает `test.Hello(...)`
- Repository слой для теста: [internal/repository/postgres/test_repo.go](./internal/repository/postgres/test_repo.go)
  - `Hello` делает `select 1` и возвращает `"Hello!"` (чтобы запрос точно прошел через repository)

## ЛР2 — Подключение PostgreSQL + инициализация таблиц + /dbtest

Что требуется:
- подключение к PostgreSQL
- инициализация таблиц
- `POST /dbtest`: записывает строку из тела запроса в БД
- работа с БД на слое repository

Где реализовано:
- Конфиг БД из env: [internal/config/config.go](./internal/config/config.go)
  - переменные `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- Открытие подключения: [internal/repository/postgres/db.go](./internal/repository/postgres/db.go)
  - драйвер `pgx` через `database/sql`
- Инициализация таблиц: [internal/repository/postgres/schema.go](./internal/repository/postgres/schema.go)
  - создаются таблицы `dbtest_logs`, `users`, `orders`
- Repository для /dbtest: [internal/repository/postgres/dbtest_repo.go](./internal/repository/postgres/dbtest_repo.go)
  - `InsertLine(...)` делает `insert into dbtest_logs(...)`
- Usecase для /dbtest: [internal/usecase/usecase.go](./internal/usecase/usecase.go)
  - `DBTestInsert(...)` валидирует и вызывает repository
- HTTP хэндлер для /dbtest: [internal/transport/httpapi/handler.go](./internal/transport/httpapi/handler.go)
  - `handleDBTest` читает тело и вызывает usecase

## ЛР3 — Пользователь: регистрация + логин + JWT

Что требуется:
- хэндлер регистрации (создает пользователя в БД)
- хэндлер авторизации (логин) и выдачи JWT

Где реализовано:
- Таблица пользователей создается в: [internal/repository/postgres/schema.go](./internal/repository/postgres/schema.go)
  - `users(email unique, password_hash, created_at)`
- Repository пользователей: [internal/repository/postgres/users_repo.go](./internal/repository/postgres/users_repo.go)
  - `Create(...)`, `FindByEmail(...)`
- Usecase регистрации/логина: [internal/usecase/usecase.go](./internal/usecase/usecase.go)
  - `Register(...)` делает bcrypt-хеш и создает пользователя
  - `Login(...)` проверяет пароль и вызывает JWT issuance
- JWT менеджер: [internal/usecase/authjwt/jwt.go](./internal/usecase/authjwt/jwt.go)
  - `Issue(userID)` создает токен (HS256)
  - `ParseUserID(token)` парсит токен и возвращает `userID`
- HTTP хэндлеры:
  - `POST /auth/register`: [internal/transport/httpapi/handler.go](./internal/transport/httpapi/handler.go)
  - `POST /auth/login`: [internal/transport/httpapi/handler.go](./internal/transport/httpapi/handler.go)

Переменные JWT:
- `JWT_SECRET`, `JWT_ISSUER`, `JWT_TTL` (см. [.env](./.env) и [internal/config/config.go](./internal/config/config.go))

## ЛР4 — Заказы: создание + получение статусов всех заказов пользователя

Что требуется:
- создать новый заказ
- вернуть статусы всех заказов пользователя

Где реализовано:
- Таблица заказов создается в: [internal/repository/postgres/schema.go](./internal/repository/postgres/schema.go)
  - `orders(user_id -> users, status default 'created', created_at)`
- Repository заказов: [internal/repository/postgres/orders_repo.go](./internal/repository/postgres/orders_repo.go)
  - `Create(userID)`, `ListByUser(userID)`
- Usecase:
  - `CreateOrder(...)`, `ListOrders(...)`: [internal/usecase/usecase.go](./internal/usecase/usecase.go)
- HTTP:
  - `POST /orders`, `GET /orders`: [internal/transport/httpapi/handler.go](./internal/transport/httpapi/handler.go)

## ЛР5 — Middleware: JWT -> userID в хэндлерах заказов

Что требуется:
- middleware, который принимает JWT и передает в хэндлеры заказов ID пользователя

Где реализовано:
- Middleware: [internal/transport/httpapi/handler.go](./internal/transport/httpapi/handler.go)
  - `authMiddleware(...)`:
    - читает `Authorization: Bearer <token>`
    - парсит токен (`ParseUserID`)
    - кладет `userID` в `context.Context`
  - `userIDFromCtx(...)` достает `userID` для `handleCreateOrder` / `handleListOrders`

