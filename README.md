```
docker-compose up -d
```

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

Из горячего кэша, потому что "Продукт 50" нет в кэше
```http
http://localhost:8080/product/50
```
