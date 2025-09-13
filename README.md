BankKibikov

Учебное банковское приложение на Go с двухфакторной аутентификацией через email (временно OTP печатается в консоль для отладки).

---

## 🚀 Запуск проекта

### 1. Запуск базы данных PostgreSQL
В корне проекта есть `docker-compose.yaml`. Запусти:

```bash
docker-compose up -d
Проверка:

bash
Копировать код
docker ps
Должен быть контейнер bank-postgres.

2. Применить миграции
Создаём таблицу users:

bash
Копировать код
docker exec -i bank-postgres psql -U bankuser -d bankdb < migrations/001_create_users.sql
3. Запуск приложения
bash
Копировать код
go run ./cmd/server
Сервер стартует на порту :8082.

📌 Роуты
Healthcheck
http
Копировать код
GET /ping
Ответ:

json
Копировать код
{ "status": "ok" }
Создание пользователя
http
Копировать код
POST /users
Content-Type: application/json

{
  "username": "max",
  "password": "12345",
  "email": "test@example.com"
}
Ответ:

json
Копировать код
{
  "id": "uuid",
  "status": "user created"
}
Логин (шаг 1: проверка пароля + генерация OTP)
http
Копировать код
POST /login
Content-Type: application/json

{
  "username": "max",
  "password": "12345"
}
Ответ:

json
Копировать код
{ "status": "OTP sent to email" }
⚠️ OTP сейчас не отправляется на email, а печатается в консоль приложения:

yaml
Копировать код
DEBUG OTP for user max: 482190 (expires 2025-09-13T14:25:00Z)
Подтверждение OTP (шаг 2)
http
Копировать код
POST /verify-otp
Content-Type: application/json

{
  "username": "max",
  "otp": "482190"
}
Ответ:

json
Копировать код
{ "status": "login successful" }
Получить всех пользователей
http
Копировать код
GET /users
Получить пользователя по ID
http
Копировать код
GET /users/:id
🗂 Структура проекта
bash
Копировать код
BankKibikov/
├── cmd/server/              # main.go — точка входа, запуск HTTP сервера
├── configs/                 # конфиги (config.yaml)
├── internal/
│   ├── db/                  # postgres.go — инициализация соединения с БД
│   ├── handler/             # обработчики HTTP-запросов
│   │   ├── handler.go       # маршруты
│   │   ├── user_handler.go  # логика /users
│   │   ├── auth_handler.go  # логика /login и /verify-otp
│   │   └── error.go         # единая обработка ошибок
│   ├── logger/              # logger.go — инициализация zap-логгера
│   ├── models/              # user.go — описание структур (User)
│   ├── repository/          # user_repository.go — работа с БД
│   ├── security/            # auth_service.go — 2FA логика (OTP)
│   └── service/             # user_service.go — бизнес-логика (создание юзеров)
├── migrations/              # SQL миграции для создания таблиц
│   └── 001_create_users.sql
└── docker-compose.yaml      # контейнер с Postgres
⚙️ Роль файлов
cmd/server/main.go → запускает сервер, подключает логгер, БД, репозитории, сервисы, хендлеры.

internal/db/postgres.go → инициализация PostgreSQL (pgxpool).

internal/logger/logger.go → zap-логгер.

internal/models/user.go → структура пользователя.

internal/repository/user_repository.go → SQL-запросы (INSERT, SELECT, UPDATE).

internal/service/user_service.go → бизнес-логика пользователей (валидация, дефолтный пароль).

internal/security/auth_service.go → логика 2FA: генерация OTP, проверка, временно печатает код в консоль.

internal/handler/user_handler.go → HTTP API для /users.

internal/handler/auth_handler.go → HTTP API для /login и /verify-otp.

internal/handler/handler.go → собирает все маршруты.

internal/handler/error.go → единый обработчик ошибок (возвращает JSON).

migrations/001_create_users.sql → SQL для таблицы users.

docker-compose.yaml → запуск PostgreSQL.

📝 Как работает 2FA
Пользователь регистрируется через /users.

При логине (/login) проверяется username + password.

Если всё верно → генерируется одноразовый код (OTP), сохраняется в БД и печатается в консоль.

Пользователь вводит этот код в /verify-otp.

Если код совпал и не просрочен → логин успешный.
