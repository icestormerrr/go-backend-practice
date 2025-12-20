## Практическая работа № 12

Студент: Юркин В.И.

Группа: ПИМО-01-25

Тема: Подключение Swagger/OpenAPI. Автоматическая генерация документации

Цели:
1.	Освоить основы спецификации OpenAPI (Swagger) для REST API.
2.	Подключить автогенерацию документации к проекту из ПЗ 11 (notes-api).
3.	Научиться публиковать интерактивную документацию (Swagger UI / ReDoc) на эндпоинте GET /docs.
4.	Синхронизировать код и спецификацию (комментарии-аннотации → генерация) и/или «schema-first» (генерация кода из openapi.yaml).
5.	Подготовить процесс обновления документации (Makefile/скрипт).

Подход к созданию OAS: code first

## Структура проекта

```
pz12-openapi
├── docs
│   ├── docs.go # go файл для работы swagger
│   ├── swagger.json # open api схема в формате json
│   ├── swagger.yaml # open api схема в формате yaml
│   └── redoc.html # статика для работы redoc
├── cmd
│   └── server
│       └── main.go           # Точка входа приложения
├── internal
│   ├── core                   # Основные доменные интерфейсы и контракты
│   │   ├── domains.go         # Описание доменных сущностей и структур данных
│   │   ├── repos.go           # Интерфейсы репозиториев
│   │   └── service.go         # Интерфейсы сервисов
│   ├── delivery              # Слой взаимодействия с приложением
│   │   ├── http							 # Работа с HTTP запросами
│   │   │   ├── handlers        # HTTP handlers
│   │   │   │   ├── notes_handler.go  # Обработчики для /notes
│   │   │   └── middleware      # HTTP middleware
│   │   │       ├── logger.go    # Логирование запросов
│   │   │       └── recover.go   # Обработка паник в HTTP
│   │   └── router.go           # Маршрутизация HTTP запросов
│   ├── repos                  # Реализации репозиториев
│   │   └── notes_inmemory_repo.go   # Репозиторий notes в памяти
│   ├── services               # Реализация бизнес-логики
│   │   └── notes_service.go     # Сервис работы с notes
│   └── utils                  # Вспомогательные модули
│       ├── config
│       │   └── config.go      # Загрузка конфигурации из env/файлов
│       ├── http
│       │   └── http.go        # Общие HTTP-утилиты
```

## Команды

Установка зависимостей
```bash
make install
```

Запуск проекта
```bash
make run
```

Билд проекта
```bash
make build
```

Генерация Open api схемы
```bash
make gen-oas
```

Валидация Open api схемы
```bash
make validate-oas
```

## Конфигурация
.env
```
# Порт, на котором запускается приложение
APP_PORT=8080
```

## Запуск

Docker: 25.0.3
Golang: 1.24.0

### Локально
1. Создание .env файла (см. .env.example)
2. Установка зависимостей
```bash
make install
```
3. Запуск сервера
```bash
make run
```

### На сервере
1. Создание .env файла (см. .env.example)
2. Развёртывание сервера
```bash
docker-compose up --build -d
```

## Пример документации swagger
![alt text](image-13.png)

# Пример документации redoc
![alt text](image-14.png)

## Пример аннотаций над методами
```golang
// UpdateNote godoc
// @Summary      Обновить заметку (частично)
// @Tags         notes
// @Accept       json
// @Param        id     path   int        true  "Идентификатор заметки"
// @Param        input  body   NoteUpdateRequest true  "Поля для обновления"
// @Success      200    {object}  core.Note
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /notes/{id} [patch]
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid note id", err.Error())
		return
	}

	var req NoteUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	payload := core.NoteUpdatePayload{
		Title:   req.Title,
		Content: req.Content,
	}

	if err := h.service.UpdateNote(id, payload); err != nil {
		if err.Error() == "note not found" {
			http_utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		http_utils.WriteError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteNote godoc
// @Summary      Удалить заметку
// @Tags         notes
// @Param        id  path  int  true  "Bltynbabrfnjh pfvtnrb"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notes/{id} [delete]
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid note id", err.Error())
		return
	}

	if err := h.service.Delete(id); err != nil {
		http_utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(idStr string) (int64, error) {
	return strconv.ParseInt(idStr, 10, 64)
}
```