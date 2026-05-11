# Практическая работа № 30

Студент: Юркин В.И.

Группа: ПИМО-01-25

Тема: Реализация очереди задач (producer–consumer)

Цель: Освоить построение очереди задач по модели producer–consumer с использованием RabbitMQ, научиться организовывать повторные попытки обработки, настраивать очередь проблемных сообщений (DLQ), а также реализовывать базовую идемпотентность обработчика для защиты от повторной обработки одного и того же сообщения. Логика занятия продолжает предыдущую работу с RabbitMQ: теперь вместо простого события рассматривается именно задача, которая может обрабатываться долго, завершаться ошибкой и требовать повторной попытки.

## Что реализовано

- отдельный HTTP producer `services/tasks`
- отдельный worker consumer `services/worker`
- основная очередь `task_jobs`
- очередь проблемных сообщений `task_jobs_dlq`
- ограничение повторных попыток до `3`
- перевод сообщения в DLQ после превышения лимита попыток
- идемпотентная проверка по `message_id`
- имитация долгой обработки через `time.Sleep(2 * time.Second)`
- управляемая ошибка для `task_id == "t_fail"`

## Структура проекта

```text
tech-ip-sem2-task-queue/
├── deploy/
│   └── rabbit/
│       └── docker-compose.yml
├── internal/
│   ├── amqpclient/
│   ├── jobs/
│   ├── publisher/
│   └── rabbitsetup/
├── services/
│   ├── tasks/
│   │   ├── cmd/tasks/main.go
│   │   └── internal/
│   │       └── http/
│   └── worker/
│       ├── cmd/worker/main.go
│       └── internal/
│           ├── consumer/
│           └── store/
└── README.md
```

## Запуск RabbitMQ

```powershell
cd deploy/rabbit
docker compose up -d
docker compose ps
```

RabbitMQ Management UI:

- `http://localhost:15672`
- login: `guest`
- password: `guest`

## Запуск worker

```powershell
cd services/worker
$env:RABBIT_URL="amqp://guest:guest@localhost:5672/"
$env:PREFETCH_COUNT="1"
go run ./cmd/worker
```

## Запуск tasks

```powershell
cd services/tasks
$env:TASKS_PORT="8082"
$env:RABBIT_URL="amqp://guest:guest@localhost:5672/"
go run ./cmd/tasks
```

## Проверка

```powershell
go test ./...
```

## API

Маршрут:

`POST /v1/jobs/process-task`

Тело запроса:

```json
{
  "task_id": "t_001"
}
```


## Логика worker

Worker выполняет такой сценарий:

1. получает сообщение из `task_jobs`
2. проверяет `message_id` в in-memory store
3. если сообщение уже обрабатывалось, делает `ack` и не выполняет работу повторно
4. если сообщение новое, запускает обработку
5. если обработка успешна, сохраняет `message_id` и делает `ack`
6. если обработка неуспешна, увеличивает `attempt`
7. если `attempt <= 3`, публикует сообщение повторно в `task_jobs`
8. если `attempt > 3`, публикует сообщение в `task_jobs_dlq`


## Проверка успешной обработки

```powershell
$body = @{
  task_id = "t_001"
} | ConvertTo-Json -Compress

(Invoke-WebRequest `
  -UseBasicParsing `
  -Uri http://localhost:8082/v1/jobs/process-task `
  -Method Post `
  -ContentType "application/json" `
  -Body $body).Content
```

Ожидаемое поведение:

![alt text](docs/image.png)

Логи worker:
![alt text](docs/image-1.png)

## Проверка retries и DLQ

```powershell
$body = @{
  task_id = "t_fail"
} | ConvertTo-Json -Compress

(Invoke-WebRequest `
  -UseBasicParsing `
  -Uri http://localhost:8082/v1/jobs/process-task `
  -Method Post `
  -ContentType "application/json" `
  -Body $body).Content
```

Ожидаемое поведение:

![alt text](docs/image-2.png)

- первая попытка: ошибка
- вторая попытка: ошибка
- третья попытка: ошибка
- после превышения лимита сообщение публикуется в `task_jobs_dlq`

Логи worker:

![alt text](docs/image-3.png)

## Проверка через RabbitMQ Management UI

![alt text](docs/image-4.png)

- существует очередь `task_jobs`
- существует очередь `task_jobs_dlq`


## Итог

В проекте реализована классическая очередь задач по модели producer-consumer. HTTP-сервис принимает запрос и быстро ставит задачу в очередь, а тяжёлая обработка выполняется отдельным worker. При временной или искусственно смоделированной ошибке сообщение не теряется, а повторно публикуется с увеличенным `attempt`. После превышения лимита задача уходит в DLQ, что позволяет не блокировать основную очередь и не зацикливать проблемное сообщение. Дополнительно показан базовый принцип идемпотентности, который защищает от повторной обработки одного и того же `message_id`.
