# 🚗 GoRent

Сервис аренды автомобилей на Go с автодокументацией через Swagger.

## 🚀 Быстрый старт

```bash
# Клонирование и запуск
git clone https://github.com/scwern/gorent.git
cd gorent
docker-compose up -d

Доступ после запуска:

API: http://localhost:8080

Swagger Docs: http://localhost:8080/swagger/index.html

База данных: 5432 (PostgreSQL)

📖 API возможности
Аутентификация (JWT через pkg/jwt)

Управление автомобилями (добавление, поиск, бронь)

История аренды пользователя

Миграции БД (из папки migrations)

⚙️ Конфигурация
Настройки через переменные окружения (см. internal/config):

DB_URL - подключение к PostgreSQL

JWT_SECRET - секретный ключ

SERVER_PORT - порт приложения
