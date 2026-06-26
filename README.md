# Medical Task Tracker

REST API для трекера задач медицинского персонала.

В репозитории реализованы:
- CRUD задач;
- теги задач;
- периодические задачи (weekly, monthly, yearly, biweekly, shift, parity);
- отметка выполнения обычных задач и отдельных экземпляров периодических задач.

## Требования

- Go `1.26+`
- Docker + Docker Compose (для запуска с PostgreSQL)

## Переменные окружения

Пример значений находится в `.env.example`.

Минимально необходимые:
- `APP_ENV` (`dev`, `test`, `prod`)
- `APP_PORT` (по умолчанию `8000`)
- `LOG_LEVEL` (`debug`, `info`)
- `POSTGRES_HOST`
- `POSTGRES_PORT`
- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`

## Локальный запуск (без Docker)

1. Скопировать пример env:
   ```bash
   cp .env.example .env
   ```
2. Поднять PostgreSQL (локально или в контейнере).
3. Применить миграции:
   ```bash
   go run ./cmd migration up
   ```
4. Запустить API:
   ```bash
   go run ./cmd run
   ```

API будет доступен на `http://localhost:8000`.

## Запуск через Docker Compose

```bash
cp .env.example .env
docker compose up --build
```

Контейнер `migrate` применяет миграции перед запуском `backend`.

## Миграции

Применить:
```bash
go run ./cmd migration up
```

Откатить последнюю:
```bash
go run ./cmd migration down
```

Создать новую:
```bash
go run ./cmd migration create add_new_feature
```

Для `create` нужен установленный `goose`:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Тестирование

Интеграционный smoke-скрипт:
```bash
bash ./test.bash
```

По умолчанию скрипт использует `BASE_URL=http://localhost:8000`.

## API документация

- OpenAPI: `docs/swagger.yaml`
- Основной префикс: `/api/v1`

Ключевые endpoints:
- `POST /api/v1/tasks/create`
- `GET /api/v1/tasks`
- `GET /api/v1/tasks/{id}`
- `PUT /api/v1/tasks/{id}`
- `DELETE /api/v1/tasks/{id}`
- `POST /api/v1/tasks/{id}/tags`
- `DELETE /api/v1/tasks/{id}/tags`
- `POST /api/v1/tasks/{id}/done`

## Принятые допущения

- Проблема бесконечности решена lazy-подходом: экземпляры периодических задач генерируются на лету для запрошенного окна, а не материализуются заранее в БД.
- Для периодических задач выполнение фиксируется по конкретной дате через `occurrence_date`.
- "Задачи на конкретные даты" моделируются как набор обычных (не периодических) задач.
- Четность/нечетность (`parity`) рассчитывается по числу дня месяца.
- Бонус-сценарий "исключения из правил" (перенос/отмена одного экземпляра серии) пока не реализован.
