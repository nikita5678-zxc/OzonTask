# GraphQL API

Система для добавления и чтения постов и комментариев с использованием GraphQL, аналогичная комментариям к постам на популярных платформах, таких как Хабр или Reddit.

## Описание

Проект реализует GraphQL API с использованием следующих технологий:
- Go (основной язык программирования)
- PostgreSQL (база данных)
- GraphQL (API интерфейс)
- Docker (контейнеризация)

## Структура проекта

```
OzonTask/
├── api/           # GraphQL API реализация
├── db/            # Работа с базой данных
├── main.go        # Точка входа в приложение
├── Dockerfile     # Конфигурация Docker
└── go.mod         # Зависимости проекта
```

## Требования

- Go 1.23 или выше
- PostgreSQL
- Docker (опционально)

## Установка и запуск

### Локальная установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/nikita5678-zxc/OzonTask.git
cd OzonTask
```

2. Установите зависимости:
```bash
go mod download
```

3. Создайте файл .env с настройками базы данных:
```
DB_HOST=localhost
DB_PORT=6432
DB_USER=postgres
DB_PASSWORD=123
DB_NAME=postgres
```

4. Запустите приложение:
```bash
go run main.go
```

### Запуск через Docker

```bash
docker-compose up --build
```

## Использование

После запуска приложения доступны следующие эндпоинты:

- `http://localhost:3000/` - GraphQL Playground (интерактивная документация)
- `http://localhost:3000/graphql` - GraphQL API

## API Endpoints

### GraphQL

Основной эндпоинт для работы с API:
```
POST http://localhost:3000/graphql
```

### GraphQL Playground

Интерактивная документация и тестирование API:
```
GET http://localhost:3000/
```

## Разработка

Для разработки рекомендуется использовать GraphQL Playground для тестирования запросов и мутаций.

