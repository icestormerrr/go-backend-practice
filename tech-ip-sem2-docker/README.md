# Практическая работа № 23

Студент: Юркин В.И.

Группа: ПИМО-01-25

Тема: Написание Dockerfile и сборка контейнера

Цель: Освоить контейнеризацию backend-приложения на Go с помощью Docker, научиться писать Dockerfile, собирать Docker-образ и запускать контейнеризированный сервис в воспроизводимой среде.

## Что реализовано

- минимальный `tasks` HTTP-сервис с маршрутом `GET /health`
- отдельный `Dockerfile` для multi-stage сборки Go-приложения
- `.dockerignore` для чистого build context
- `docker-compose.yml` для запуска сервиса через Docker Compose
- env-конфигурация порта приложения через `TASKS_PORT`

## Структура

```text
tech-ip-sem2-docker/                  - корень проекта практической работы
├── services/
│   └── tasks/                        - учебный tasks-сервис
│       ├── cmd/
│       │   └── tasks/
│       │       └── main.go           - точка входа HTTP-сервиса
│       ├── .dockerignore             - исключения из Docker build context
│       ├── Dockerfile                - multi-stage сборка контейнера
│       └── go.mod                    - Go-модуль tasks-сервиса
├── deploy/
│   └── docker-compose.yml            - запуск tasks через Docker Compose
└── README.md                         - инструкция сборки и запуска
```

## Локальный запуск без Docker

Из каталога `services/tasks`:

```powershell
go run ./cmd/tasks
```

![alt text](docs/image-1.png)

Проверка:

![alt text](docs/image.png)


## Сборка образа

Из каталога `services/tasks`:

```powershell
docker build -t techip-tasks:0.1 .
```

![alt text](docs/image-2.png)


## Запуск контейнера

```powershell
docker run --rm -p 8082:8082 -e TASKS_PORT=8082 techip-tasks:0.1
```

![alt text](docs/image-3.png)

Проверка:

```powershell
Invoke-WebRequest `
  -Uri "http://localhost:8082/health" `
  -Method Get
```

![alt text](docs/image-4.png)

## Запуск через Docker Compose

Из каталога `deploy`:

```powershell
docker compose up -d --build
```

![alt text](docs/image-5.png)

Проверка статуса:

```powershell
docker compose ps
```

![alt text](docs/image-6.png)

Просмотр логов:

```powershell
docker compose logs -f
```

![alt text](docs/image-7.png)

После запуска через Compose проверка остаётся той же:

```powershell
Invoke-WebRequest `
  -Uri "http://localhost:8082/health" `
  -Method Get
```
![alt text](docs/image-8.png)

Остановка:

```powershell
docker compose down
```

![alt text](docs/image-9.png)

![alt text](docs/image-10.png)