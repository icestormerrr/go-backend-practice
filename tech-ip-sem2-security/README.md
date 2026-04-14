# Практическая работа № 21

Студент: Юркин В.И.

Группа: ПИМО-01-25

Тема: Реализация HTTPS (TLS-сертификаты). Защита от SQL-инъекций

Цель: Освоить базовые практические подходы к защите backend-приложения на Go за счёт включения HTTPS с использованием TLS-сертификата и предотвращения SQL-инъекций при работе с базой данных


## Структура

```text
tech-ip-sem2-security/                - корень проекта практической работы
├── cmd/
│   ├── gencert/
│   │   └── main.go                   - генерация локального self-signed сертификата
│   └── server/
│       └── main.go                   - запуск HTTPS-сервера и подключение PostgreSQL
├── certs/
│   ├── server.crt                    - TLS-сертификат сервера
│   └── server.key                    - приватный ключ сервера
├── internal/
│   ├── config/
│   │   └── config.go                 - env-конфигурация HTTPS и PostgreSQL
│   ├── httpapi/
│   │   └── handler.go                - handlers /health и /students
│   └── student/
│       ├── model.go                  - структура Student
│       └── repo.go                   - безопасные и небезопасные SQL-примеры
├── sql/
│   └── init.sql                      - создание таблицы students и тестовых данных
├── go.mod                            - Go-модуль проекта
└── README.md                         - инструкция запуска и проверки
```

## Подготовка

### Переменные окружения

- `HTTPS_ADDR` - адрес HTTPS-сервера, по умолчанию `:8443`
- `TLS_CERT_FILE` - путь к сертификату, по умолчанию `certs/server.crt`
- `TLS_KEY_FILE` - путь к приватному ключу, по умолчанию `certs/server.key`
- `DATABASE_DSN` - строка подключения PostgreSQL, по умолчанию `postgres://postgres:postgres@localhost:5432/study_security?sslmode=disable`

### Подготовка базы данных

1. Запустите postgresql
> Если у вас нет postgres поднимите через docker
```
docker run --name my-postgres-container -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
docker exec -it my-postgres-container psql -U postgres
```

2. Создайте базу данных `study_security` и выполните скрипт из `sql/init.sql`.

```powershell
CREATE DATABASE study_security;
\c study_security
<вставьте скрипт `sql/init.sql`>
```

![alt text](docs/image-3.png)

### Генерация сертификата

В проекте используется локальный self-signed сертификат для учебной среды. Файлы лежат в каталоге `certs/`.

Быстрый способ без внешних инструментов:

```powershell
go run ./cmd/gencert
```

## Запуск приложения

Из корня проекта:

```powershell
go run ./cmd/server
```

Ожидаемое сообщение:

```text
HTTPS server started on https://localhost:8443
```

## Проверка HTTPS

1. Проверка health:

```powershell
curl.exe -k https://localhost:8443/health
```

Ожидаемый ответ:

![alt text](docs/image.png)

2. Проверка чтения студента:

```powershell
curl.exe -k "https://localhost:8443/students?id=1"
```

Ожидаемый ответ:

![alt text](docs/image-1.png)

3. Проверка отсутствующего студента:

```powershell
curl.exe -k "https://localhost:8443/students?id=999"
```

Ожидаемый ответ:

![alt text](docs/image-2.png)

## SQL-инъекция: опасный и безопасный подход

Небезопасный подход:

```go
query := "SELECT id, full_name, study_group, email FROM students WHERE id = " + rawID
```

Такой код опасен, потому что пользовательский ввод становится частью SQL-строки.

Безопасный подход:

```go
row := db.QueryRow(
    "SELECT id, full_name, study_group, email FROM students WHERE id = $1",
    id,
)
```

В проекте основной маршрут использует подготовленный запрос:

```go
stmt, err := db.Prepare("SELECT id, full_name, study_group, email FROM students WHERE id = $1")
```

