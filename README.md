# 🚗 GoRent

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?logo=postgresql)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker)](https://www.docker.com/)
[![Swagger](https://img.shields.io/badge/Swagger-85EA2D?logo=swagger)](https://swagger.io/)

Сервис аренды автомобилей на Go с автодокументацией через Swagger.

## 🚀 Быстрый старт

```bash
Клонирование и запуск
git clone https://github.com/scwern/gorent.git
cd gorent
docker-compose up -d
Доступ после запуска:

API: http://localhost:8080

Swagger Docs: http://localhost:8080/swagger/index.html

База данных: 5432 (PostgreSQL)

API возможности
Аутентификация (JWT через pkg/jwt)

Управление автомобилями (добавление, поиск, бронь)

История аренды пользователя

Миграции БД (из папки migrations)

Конфигурация
Настройки через переменные окружения (internal/config)
