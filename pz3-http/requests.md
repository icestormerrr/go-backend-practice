## Тестовые запросы (Windows)

Ниже приведён список тестовых запросов, для демонстрации возможностей сервера их нужно выполнять последовательно.

### 1. Проверка работы сервера
```bash
curl http://localhost:8080/health
```

### 2. Создание нового дела
```bash
curl -Method POST http://localhost:8080/tasks -Body '{"title":"Buy milk"}' -Headers @{"Content-Type"="application/json"} 
```

### 3. Создание второго нового дела
```bash
curl -Method POST http://localhost:8080/tasks -Body '{"title":"Send letter"}' -Headers @{"Content-Type"="application/json"}
```

### 4. Получение списка дел
```bash
curl http://localhost:8080/tasks
```

### 5. Получение списка дел (с фильтрацией по названию)
```bash
curl "http://localhost:8080/tasks?q=milk"
```

### 6. Обновление статуса выполнения дела
```bash
curl -Method PATCH http://localhost:8080/tasks/1 -Body '{"done":true}' -Headers @{"Content-Type"="application/json"}
```

### 7. Получение списка дел
```bash
curl http://localhost:8080/tasks
```

### 8. Удаление дела
```bash
curl -Method DELETE http://localhost:8080/tasks/1
```

### 9. Получение списка дел
```bash
curl http://localhost:8080/tasks
```

### 10. Получение дела по id (успешное)
```bash
curl http://localhost:8080/tasks/2
```

### 11. Получение дела по id (несуществующий)
```bash
curl http://localhost:8080/tasks/22
```

### 12. Создание нового дела (с пустым названием)
```bash
curl -Method POST http://localhost:8080/tasks -Body '{"title":""}' -Headers @{"Content-Type"="application/json"}
```


