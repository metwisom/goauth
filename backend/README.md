Сервис аутентификации OAuth
Высокопроизводительный сервис аутентификации OAuth 2.0 с поддержкой традиционного входа по логину/паролю и аутентификации через Steam.
Возможности

🔐 Аутентификация, соответствующая стандарту OAuth 2.0
🚀 Высокопроизводительный HTTP-сервер на базе FastHTTP
🔑 Несколько методов аутентификации:
Логин/пароль
Steam OpenID


🛡️ Безопасное хеширование паролей с использованием bcrypt
📦 Поддержка базы данных PostgreSQL
🐳 Поддержка Docker для простого развертывания
🔄 Управление сессиями
🎯 Управление клиентами

Требования

Go 1.24 или новее
PostgreSQL 14 или новее
Docker и Docker Compose (опционально)

Установка

Клонируйте репозиторий:

git clone https://github.com/yourusername/goauth.git
cd goauth


Установите зависимости:

go mod download


Настройте переменные окружения:

cp .env.example .env
# Отредактируйте .env с вашими настройками

Конфигурация
Требуются следующие переменные окружения:
SESSION_COOKIE_NAME=session
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
STEAM_KEY=your_steam_api_key
APP_ENV=dev

Запуск с помощью Docker

Запустите сервисы:

docker-compose up -d


Сервис будет доступен по адресу http://localhost:8080

Локальный запуск

Запустите PostgreSQL (если не используете Docker):

# С помощью Docker
docker-compose up -d postgres

# Или запустите локальный экземпляр PostgreSQL


Запустите приложение:

go run cmd/main.go

API-эндпоинты
Аутентификация

POST /api/register — Регистрация нового пользователя
POST /api/login — Вход с использованием логина и пароля
POST /api/steam — Аутентификация через Steam
POST /api/logout — Выход из текущей сессии

OAuth 2.0

GET /api/authorize — Эндпоинт авторизации OAuth
POST /api/token — Эндпоинт получения токена OAuth

Управление пользователями

GET /api/getme — Получение информации о текущем пользователе

Структура проекта
.
├── cmd/                # Точки входа приложения
├── internal/           # Приватный код приложения
│   ├── app/           # Инициализация приложения
│   ├── config/        # Управление конфигурацией
│   ├── db/            # Операции с базой данных
│   ├── httpServer/    # Реализация HTTP-сервера
│   ├── libs/          # Интеграция с внешними библиотеками
│   ├── model/         # Модели данных
│   └── utils/         # Вспомогательные функции
├── db/                # Миграции базы данных
└── docker-compose.yml # Конфигурация Docker

Разработка
Запуск тестов
go test ./...

Сборка
go build -o goauth cmd/main.go

Вклад в проект

Сделайте форк репозитория
Создайте ветку для новой функции (git checkout -b feature/amazing-feature)
Зафиксируйте изменения (git commit -m 'Add some amazing feature')
Отправьте ветку в репозиторий (git push origin feature/amazing-feature)
Создайте Pull Request

Лицензия
Проект распространяется под лицензией MIT — подробности в файле LICENSE.
Безопасность
О любых проблемах безопасности сообщайте по адресу security@example.com
Благодарности

FastHTTP — Высокопроизводительный HTTP-сервер
PostgreSQL — База данных
Steam OpenID — Аутентификация через Steam

    