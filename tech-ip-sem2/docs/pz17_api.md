# PZ17 API

## Сервисы и границы ответственности

`Auth service` отвечает только за учебную авторизацию:
- принимает логин и пароль;
- возвращает упрощённый Bearer-токен;
- проверяет валидность токена;
- не хранит задачи и не знает о CRUD-операциях.

`Tasks service` отвечает только за задачи:
- хранит задачи в памяти;
- реализует CRUD по REST;
- перед каждой операцией вызывает `Auth service`;
- прокидывает `X-Request-ID` во внутренний HTTP-вызов.

## Переменные окружения

### Auth service

- `AUTH_PORT` - порт сервиса, по умолчанию `8081`

### Tasks service

- `TASKS_PORT` - порт сервиса, по умолчанию `8082`
- `AUTH_BASE_URL` - базовый адрес Auth service, по умолчанию `http://localhost:8081`
- `AUTH_TIMEOUT_MS` - таймаут межсервисного вызова в миллисекундах, по умолчанию `2500`

## Auth service

### POST /v1/auth/login

Запрос:

```json
{
  "username": "student",
  "password": "student"
}
```

Успешный ответ `200 OK`:

```json
{
  "access_token": "demo-token",
  "token_type": "Bearer"
}
```

Ошибки:
- `400 Bad Request` - неверный JSON или не переданы поля
- `401 Unauthorized` - неверные учётные данные

Пример:

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

### GET /v1/auth/verify

Заголовки:
- `Authorization: Bearer demo-token`
- `X-Request-ID: req-002`

Успешный ответ `200 OK`:

```json
{
  "valid": true,
  "subject": "student"
}
```

Ошибка `401 Unauthorized`:

```json
{
  "valid": false,
  "error": "unauthorized"
}
```

Пример:

```powershell
Invoke-WebRequest `
  -Uri "http://localhost:8081/v1/auth/verify" `
  -Method Get `
  -Headers @{
    Authorization = "Bearer demo-token"
    "X-Request-ID" = "req-002"
  }
```

## Tasks service

Во всех запросах нужен заголовок `Authorization: Bearer demo-token`.
Опционально можно передавать `X-Request-ID`.

### POST /v1/tasks

Запрос:

```json
{
  "title": "Do PZ17",
  "description": "split services",
  "due_date": "2026-01-10"
}
```

Ответ `201 Created`:

```json
{
  "id": "t_001",
  "title": "Do PZ17",
  "description": "split services",
  "due_date": "2026-01-10",
  "done": false
}
```

Пример:

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

### GET /v1/tasks

Ответ `200 OK`:

```json
[
  {
    "id": "t_001",
    "title": "Do PZ17",
    "description": "split services",
    "due_date": "2026-01-10",
    "done": false
  }
]
```

Пример:

```powershell
Invoke-WebRequest `
  -Uri "http://localhost:8082/v1/tasks" `
  -Method Get `
  -Headers @{
    Authorization = "Bearer demo-token"
    "X-Request-ID" = "req-005"
  }
```

### GET /v1/tasks/{id}

Пример:

```powershell
Invoke-WebRequest `
  -Uri "http://localhost:8082/v1/tasks/t_001" `
  -Method Get `
  -Headers @{
    Authorization = "Bearer demo-token"
    "X-Request-ID" = "req-006"
  }
```

### PATCH /v1/tasks/{id}

Запрос:

```json
{
  "title": "Read lecture (updated)",
  "done": true
}
```

Пример:

```powershell
$patchBody = @{
  title = "Read lecture (updated)"
  done = $true
} | ConvertTo-Json -Compress

Invoke-RestMethod `
  -Uri "http://localhost:8082/v1/tasks/t_001" `
  -Method Patch `
  -Headers @{
    Authorization = "Bearer demo-token"
    "X-Request-ID" = "req-007"
  } `
  -ContentType "application/json" `
  -Body $patchBody
```

### DELETE /v1/tasks/{id}

Пример:

```powershell
Invoke-WebRequest `
  -Uri "http://localhost:8082/v1/tasks/t_001" `
  -Method Delete `
  -Headers @{
    Authorization = "Bearer demo-token"
    "X-Request-ID" = "req-008"
  }
```

Ожидаемый ответ: `204 No Content`

## Ошибки Tasks service

- `400 Bad Request` - невалидный JSON или неверные поля
- `401 Unauthorized` - токен отсутствует или отклонён Auth service
- `404 Not Found` - задача не найдена
- `503 Service Unavailable` - Auth service недоступен или превысил таймаут
- `502 Bad Gateway` - Auth service вернул неожиданный ответ
