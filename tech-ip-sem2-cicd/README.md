# tech-ip-sem2-cicd

Практическая работа по настройке CI/CD для Go backend-проекта с автоматическими тестами, сборкой и Docker build.

## Что реализовано

- минимальный `tasks` сервис на Go с маршрутом `GET /health`
- unit-тест для обработчика, чтобы pipeline запускал не пустой `go test`
- `Dockerfile` с multi-stage build
- `.dockerignore` для чистого build context
- `docker-compose.yml` для локального запуска контейнера
- pipeline для GitHub Actions в `.github/workflows/ci.yml`
- альтернативный pipeline для GitLab CI в `.gitlab-ci.yml`

## Структура

```text
tech-ip-sem2-cicd/                    - корень проекта практической работы
├── .github/
│   └── workflows/
│       └── ci.yml                    - pipeline GitHub Actions
├── services/
│   └── tasks/                        - сервис, который проверяется и собирается в CI
│       ├── cmd/
│       │   └── tasks/
│       │       └── main.go           - точка входа HTTP-сервиса
│       ├── internal/
│       │   └── httpapi/
│       │       ├── handler.go        - health handler
│       │       └── handler_test.go   - unit-тест handler
│       ├── .dockerignore             - исключения из build context
│       ├── Dockerfile                - multi-stage сборка Docker-образа
│       └── go.mod                    - Go-модуль сервиса
├── deploy/
│   └── docker-compose.yml            - локальный запуск контейнера через Compose
├── .gitlab-ci.yml                    - альтернативный pipeline GitLab CI
└── README.md                         - описание практики и шагов pipeline
```

## CI и CD

- `CI` — Continuous Integration: автоматические тесты и сборка после изменения кода
- `CD` — Continuous Delivery / Deployment: подготовка артефакта к доставке, публикация образа и, при необходимости, деплой

В этом проекте обязательная часть pipeline:
- checkout репозитория
- настройка Go
- `go test ./...`
- `go build ./...`
- `docker build`

## Выбранная платформа

Основной вариант в проекте — `GitHub Actions`.

Файл:

```text
.github/workflows/ci.yml
```

Дополнительно для отчёта добавлен эквивалентный вариант:

```text
.gitlab-ci.yml
```

## Локальная проверка перед CI

Из каталога `services/tasks`:

```powershell
go test ./...
go build ./...
docker build -t techip-tasks:0.1 .
```

## Локальный запуск контейнера

```powershell
docker run --rm -p 8082:8082 -e TASKS_PORT=8082 techip-tasks:0.1
```

Проверка:

```powershell
Invoke-WebRequest `
  -Uri "http://localhost:8082/health" `
  -Method Get
```

Ожидаемый ответ:

```json
{"status":"ok","service":"tasks"}
```

## Pipeline GitHub Actions

`ci.yml` содержит два job:

### 1. test-and-build

Этот job:
- получает код репозитория
- настраивает Go 1.23
- выполняет `go mod tidy`
- запускает тесты
- выполняет сборку приложения

### 2. docker-build

Этот job:
- запускается только после успешного `test-and-build`
- настраивает Docker Buildx
- собирает Docker-образ сервиса

## Тег Docker-образа

В GitHub Actions используется:

```text
${{ github.sha }}
```

В GitLab CI используется:

```text
$CI_COMMIT_SHORT_SHA
```

Это позволяет точно понимать, какая версия приложения собрана.

## Secrets и переменные

Если pipeline должен публиковать образ в registry или делать деплой, секреты нужно хранить в CI secret storage:

- `REGISTRY_USERNAME`
- `REGISTRY_PASSWORD`
- `SSH_PRIVATE_KEY`

Их нельзя:
- коммитить в репозиторий
- писать прямо в YAML
- хранить в открытом `.env`, который попадает в Git

## Опциональная публикация образа

Пример логики для GitHub Actions:

```yaml
- name: Login to registry
  run: echo "${{ secrets.REGISTRY_PASSWORD }}" | docker login -u "${{ secrets.REGISTRY_USERNAME }}" --password-stdin ghcr.io

- name: Build image
  run: docker build -t ghcr.io/my-org/techip-tasks:${{ github.sha }} .
  working-directory: ./services/tasks

- name: Push image
  run: docker push ghcr.io/my-org/techip-tasks:${{ github.sha }}
```

## Опциональный деплой

Минимальная идея деплоя:
- pipeline подключается к серверу по SSH
- выполняет `docker pull`
- затем запускает `docker compose up -d`

Пример логики:

```text
docker pull ghcr.io/my-org/techip-tasks:<tag>
docker compose up -d
```

## Что важно понять

- pipeline должен работать на проекте, который уже собирается локально
- в multi-service репозитории важно правильно указать `working-directory`
- Docker build в CI подтверждает, что образ собирается не только на машине разработчика
- секреты должны храниться только в CI variables / secrets
- автоматический деплой удобен, но требует аккуратной работы с доступами и откатами
