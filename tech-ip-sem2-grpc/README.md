# Практическая работа № 18

Студент: Юркин В.И.
Группа: ПИМО-01-25

Тема: gRPC, создание простого микросервиса, вызовы методов

## Структура

```text
tech-ip-sem2-grpc/                     - отдельный проект для gRPC-версии практики
├── proto/                             - protobuf-контракт сервиса авторизации
│   └── auth.proto                     - описание AuthService и метода Verify
│   └── authpb/                        - сгенерированные Go-файлы клиента и сервера gRPC
├── services/
│   ├── auth/                          - Auth service с HTTP login и gRPC Verify
│   │   ├── cmd/auth/main.go           - запуск HTTP и gRPC серверов
│   │   └── internal/
│   │       ├── grpc/                  - реализация gRPC-сервера AuthService
│   │       ├── http/                  - HTTP login и совместимый verify endpoint
│   │       └── service/               - бизнес-логика проверки токена
│   └── tasks/                         - Tasks service с HTTP API и gRPC-клиентом
│       ├── cmd/tasks/main.go          - запуск HTTP сервера и gRPC-клиента
│       └── internal/
│           ├── client/authgrpc/       - gRPC-клиент Tasks -> Auth
│           ├── http/                  - HTTP handlers для CRUD задач
│           └── service/               - логика задач и in-memory хранилище
├── shared/
│   └── middleware/                    - request-id и базовое логирование HTTP-запросов
├── docs/
│   ├── pz18_api.md                    - API, gRPC-контракт и маппинг ошибок
│   └── pz18_logs.md                   - примеры логов для отчёта
├── go.mod                             - Go-модуль проекта
└── README.md                          - инструкция запуска и проверки
```

## Генерация protobuf-кода

Сгенерированные файлы лежат в `proto/authpb/`.

Команда, использованная в проекте:

```powershell
protoc `
  --go_out=. --go_opt=module=tech-ip-sem2-grpc `
  --go-grpc_out=. --go-grpc_opt=module=tech-ip-sem2-grpc `
  proto/auth.proto
```

Использованные инструменты:
- `protoc 28.2`
- `protoc-gen-go v1.34.1`
- `protoc-gen-go-grpc v1.5.1`


## Запуск

### 1. Auth service

```powershell
cd tech-ip-sem2-grpc
$env:AUTH_HTTP_PORT="8081"
$env:AUTH_GRPC_PORT="50051"
go run ./services/auth/cmd/auth
```

### 2. Tasks service

```powershell
cd tech-ip-sem2-grpc
$env:TASKS_PORT="8082"
$env:AUTH_GRPC_ADDR="localhost:50051"
$env:AUTH_GRPC_TIMEOUT_MS="1500"
go run ./services/tasks/cmd/tasks
```

Если нужно вернуть значения по умолчанию, достаточно открыть новый терминал или очистить переменные окружения текущей сессии.

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
  -Headers @{ "X-Request-ID" = "req-101" } `
  -ContentType "application/json" `
  -Body $loginBody
```

Создать задачу:

```powershell
$taskBody = @{
  title = "Learn gRPC"
  description = "Replace HTTP verify with gRPC verify"
  due_date = "2026-03-25"
} | ConvertTo-Json -Compress

Invoke-RestMethod `
  -Uri "http://localhost:8082/v1/tasks" `
  -Method Post `
  -Headers @{
    Authorization = "Bearer demo-token"
    "X-Request-ID" = "req-101"
  } `
  -ContentType "application/json" `
  -Body $taskBody
```

## Материалы для отчёта

- Подробное описание API и gRPC-контракта: [docs/pz18_api.md](./docs/pz18_api.md)
- Пример реальных запросов/ответов/логов: [docs/pz18_real_requests.md](./docs/pz18_real_requests.md)