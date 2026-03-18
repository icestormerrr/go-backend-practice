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
tech-ip-sem2/
  services/
    auth/
      cmd/auth/main.go
      internal/http/
      internal/service/
    tasks/
      cmd/tasks/main.go
      internal/http/
      internal/service/
      internal/client/authclient/
  shared/
    middleware/
    httpx/
  docs/
    pz17_api.md
    pz17_diagram.md
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
