# Практическая работа № 29

Студент: Юркин В.И.

Группа: ПИМО-01-25

Тема: Подключение к RabbitMQ. Отправка и получение сообщений

Цель: Освоить базовую работу с брокером сообщений RabbitMQ в приложении на Go, научиться публиковать сообщения в очередь, реализовывать отдельного потребителя сообщений, подтверждать успешную обработку через ack и понимать назначение очередей сообщений в микросервисной архитектуре. 

## Что реализовано

- локальный запуск RabbitMQ через `docker compose`;
- durable queue `task_events`;
- persistent message при публикации;
- отдельный producer (`services/tasks`) и consumer (`services/worker`);
- ручное подтверждение обработки `Ack(false)`;
- `prefetch` через `PREFETCH_COUNT`, по умолчанию `1`;
- режим публикации `best-effort` по умолчанию и `strict` через env.

## Структура

```text
deploy/
  rabbit/
    docker-compose.yml

services/
  tasks/
    cmd/tasks/main.go
    internal/
      http/
      publisher/
      service/

  worker/
    cmd/worker/main.go
    internal/
      consumer/

internal/
  amqpclient/
  events/
```

## Переменные окружения

### tasks

- `TASKS_PORT=8082`
- `RABBIT_URL=amqp://guest:guest@localhost:5672/`
- `QUEUE_NAME=task_events`
- `PUBLISH_MODE=best-effort` или `strict`

### worker

- `RABBIT_URL=amqp://guest:guest@localhost:5672/`
- `QUEUE_NAME=task_events`
- `PREFETCH_COUNT=1`

## Запуск RabbitMQ

```powershell
cd deploy/rabbit
docker compose up -d
docker compose ps
```

RabbitMQ Management UI:

- URL: `http://localhost:15672`
- login: `guest`
- password: `guest`

![alt text](docs/image-1.png)

## Запуск worker

```powershell
cd services/worker
$env:RABBIT_URL="amqp://guest:guest@localhost:5672/"
$env:QUEUE_NAME="task_events"
$env:PREFETCH_COUNT="1"
go run ./cmd/worker
```

## Запуск tasks

```powershell
cd services/tasks
$env:TASKS_PORT="8082"
$env:RABBIT_URL="amqp://guest:guest@localhost:5672/"
$env:QUEUE_NAME="task_events"
$env:PUBLISH_MODE="best-effort"
go run ./cmd/tasks
```

## Проверка через API

```powershell
Invoke-WebRequest -Uri "http://localhost:8082/v1/tasks" `
    -Method Post `
    -Headers @{
        "Authorization" = "Bearer demo-token"
        "X-Request-ID"  = "pz13-001"
    } `
    -ContentType "application/json" `
    -Body '{"title":"Rabbit","description":"publish event"}'
```

Ожидаемо:

![alt text](docs/image.png)