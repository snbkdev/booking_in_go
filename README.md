# Fort Smyth Bed & Breakfast

Сайт бронирования небольшого гостевого дома с двумя номерами: подбор дат, бронь онлайн, письма гостю и владельцу, админка для обработки заявок.

## Стек

- Go 1.25, chi (роутинг), scs (сессии), nosurf (CSRF)
- PostgreSQL, миграции через [soda/buffalo-pop](https://gobuffalo.io/documentation/database/soda/) (`migrations/`)
- `html/template` + Bootstrap 5, собственные стили в `static/css/`
- Отправка почты — `go-simple-mail`, шаблоны писем в `email-templates/`

## Запуск

```bash
cp .env.example .env      # прописать доступы к БД и SMTP
soda migrate              # накатить миграции
./run.sh                  # http://localhost:8080
```

Для локальной почты удобен MailHog (`brew install mailhog && mailhog`), настройки по умолчанию в `.env.example` уже указывают на него.
MailHog можно поставить в docker если будут проблемы:
- docker run -d --name mailhog -p 1025:1025 -p 8025:8025 mailhog/mailhog

## Структура

```
cmd/web/        точка входа, роуты, middleware, отправка почты
internal/       handlers, render, models, forms, repository (БД и mock)
templates/      страницы и layout'ы
static/         css, js, изображения, вендорная тема админки
migrations/     схема БД
```

## Тесты

```bash
go vet ./...
go test ./...
```
