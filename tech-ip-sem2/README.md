# Практическая работа № 17

Студент: Юркин В.И.
Группа: ПИМО-01-25

Тема: Практическая работа по разделению монолита на 2 микросервиса с синхронным HTTP-взаимодействием.

## Что реализовано

- `Auth service` для учебной авторизации и проверки токена
- `Tasks service` для CRUD задач в памяти
- межсервисный вызов `Tasks -> Auth` по HTTP
- `X-Request-ID` middleware и базовое логирование
- env-конфигурация портов и адреса Auth
- таймаут на межсервисный вызов

## Структура

```text
tech-ip-sem2/                          - корень проекта практической работы
├── services/                          - каталог всех микросервисов
│   ├── auth/                          - Auth service
│   │   ├── cmd/
│   │   │   └── auth/
│   │   │       └── main.go            - точка входа Auth service
│   │   └── internal/
│   │       ├── http/                  - HTTP-роуты и handlers Auth service
│   │       └── service/               - бизнес-логика авторизации
│   └── tasks/                         - Tasks service
│       ├── cmd/
│       │   └── tasks/
│       │       └── main.go            - точка входа Tasks service
│       └── internal/
│           ├── client/
│           │   └── authclient/        - HTTP-клиент Tasks -> Auth
│           ├── http/                  - HTTP-роуты и handlers Tasks service
│           └── service/               - бизнес-логика и in-memory хранилище задач
├── shared/                            - общий код для обоих сервисов
│   ├── middleware/                    - middleware логирование
│   └── httpx/                         - создание HTTP-клиента с таймаутами
├── docs/                              - документация по API и схеме взаимодействия
│   ├── pz17_api.md                    - описание эндпоинтов, кодов ответов и примеров запросов
│   └── pz17_diagram.md                - Mermaid-диаграмма взаимодействия сервисов
│   └── pz17_real_requests.md          - Пример реальных запросов, ответов, логов
├── go.mod                             - Go-модуль проекта
```

## Запуск

### 1. Auth service

```powershell
cd tech-ip-sem2
$env:AUTH_PORT="8081"
go run ./services/auth/cmd/auth
```

### 2. Tasks service

```powershell
cd tech-ip-sem2
$env:TASKS_PORT="8082"
$env:AUTH_BASE_URL="http://localhost:8081"
$env:AUTH_TIMEOUT_MS="2500"
go run ./services/tasks/cmd/tasks
```

## Быстрая проверка (Windows)

Получить токен:

```powershell
$loginBody = @{
  username = "student"
  password = "student"
} | ConvertTo-Json -Compress

Invoke-RestMethod `
  -Uri "http://localhost:8081/v1/auth/login" `
  -Method Post `
  -Headers @{ "X-Request-ID" = "req-001" } `
  -ContentType "application/json" `
  -Body $loginBody
```

Создать задачу:

```powershell
$taskBody = @{
  title = "Do PZ17"
  description = "split services"
  due_date = "2026-01-10"
} | ConvertTo-Json -Compress

Invoke-RestMethod `
  -Uri "http://localhost:8082/v1/tasks" `
  -Method Post `
  -Headers @{
    Authorization = "Bearer demo-token"
    "X-Request-ID" = "req-003"
  } `
  -ContentType "application/json" `
  -Body $taskBody
```

Запрос без токена:

```powershell
Invoke-WebRequest `
  -Uri "http://localhost:8082/v1/tasks" `
  -Method Get `
  -Headers @{ "X-Request-ID" = "req-004" }
```

## Материалы для отчёта

1. Границы сервисов описаны в [docs/pz17_api.md](./docs/pz17_api.md)
2. Mermaid-схема находится в [docs/pz17_diagram.md](./docs/pz17_diagram.md)
3. Список эндпоинтов и curl-примеры находятся в [docs/pz17_api.md](./docs/pz17_api.md)
4. Пример реальных запросов, логов с `request-id` приведён в [docs/pz17_real_requests](./docs/pz17_real_requests.md)
