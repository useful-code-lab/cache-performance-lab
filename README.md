```
docker-compose up -d
```

Создать БД
```bash
docker exec -i $(docker ps -qf "name=postgres") psql -U user -d inventory < scripts/init.sql

```

aleksey@aleksey:~/cache-lifecycle-go$ go run main.go
[Холодный (из базы данных)] ID: 10, Time: 503.42724ms
[Горячий (из кэша)] ID: 10, Time: 183.619µs