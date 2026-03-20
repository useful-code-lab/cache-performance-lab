# cache-lifecycle-go 🚀

Проект на **Go** для демонстрации принципов **Cold** (холодного) и **Hot** (горячего) кэширования с использованием **PostgreSQL** и **Redis**.

## 📌 Описание проекта
Цель проекта — наглядно показать разницу в производительности системы при наличии прогретого кэша и при его отсутствии. 

### Основные возможности:
1. **Cache-Aside Strategy**: Приложение сначала проверяет Redis, при промахе (miss) идет в Postgres и обновляет кэш.
2. **Cold Start Simulation**: Имитация «тяжелых» запросов к БД с задержкой 500мс.
3. **Warm-up Mechanism**: Автоматический прогрев кэша популярными данными при старте сервера.
4. **Latency Analysis**: Вывод времени обработки запроса в заголовках и теле ответа.

## 🛠 Стек технологий
- **Язык**: Go 1.21+
- **База данных**: PostgreSQL 15
- **Кэш**: Redis 7
- **Инфраструктура**: Docker Compose

## 🚀 Быстрый запуск

### 1. Поднятие инфраструктуры
```bash
docker-compose up -d
```
## Этапы реализации и тестирования проекта cache-lifecycle-go

Создать БД
```bash
docker exec -i $(docker ps -qf "name=postgres") psql -U user -d inventory < scripts/init.sql

```

```bash
aleksey@aleksey:~/cache-lifecycle-go$ go run main.go
[Холодный (из базы данных)] ID: 10, Time: 503.42724ms
[Горячий (из кэша)] ID: 10, Time: 183.619µs
```

Получение данных из БД напрямую
```bash
docker exec -it $(docker ps -qf "name=postgres") psql -U user -d inventory -c "SELECT id, name, price FROM products LIMIT 10;"
```
```bash
aleksey@aleksey:~/cache-lifecycle-go$ go run main.go
```

Получение данных
Из горячего кэша, потому что "Продукт 1" в кэше
```http
http://localhost:8080/product/1
```

```JSON
{"data":{"id":1,"name":"Product 1","price":435.37081066506},"status":"HOT","time":"562.908µs"}
```


Из горячего кэша, потому что "Продукт 50" нет в кэше
```http
http://localhost:8080/product/50
```

```JSON
{"data":{"id":50,"name":"Product 50","price":120.305695569964},"status":"COLD","time":"502.455967ms"}
```

Сравнительная таблица замеров
- Cold (Холодный),	PostgreSQL (Диск),	503.42 ms,	1x (Базовое)
- Hot (Горячий),	Redis (RAM),	0.18 ms (183µs),	~2700x быстрее
