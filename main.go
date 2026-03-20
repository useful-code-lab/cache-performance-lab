package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"gitverse.ru/pavlov-systems/cache-lifecycle-go/internal/cache"
	"gitverse.ru/pavlov-systems/cache-lifecycle-go/internal/database"
	"gitverse.ru/pavlov-systems/cache-lifecycle-go/internal/warmup"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func getProduct(ctx context.Context, db *pgx.Conn, rdb *redis.Client, id int) (*Product, string) {
	key := fmt.Sprintf("product:%d", id)

	// 1. Попытка получить из кэша (HOT)
	val, err := rdb.Get(ctx, key).Result()
	if err == nil {
		var p Product
		json.Unmarshal([]byte(val), &p)
		return &p, "Горячий (из кэша)"
	}

	// 2. Если в кэше нет (COLD), идем в БД
	// Имитируем задержку БД, чтобы разница была заметна
	time.Sleep(500 * time.Millisecond)

	var p Product
	query := "SELECT id, name, price FROM products WHERE id=$1"
	err = db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		log.Println("Product not found")
		return nil, "ERROR"
	}

	// 3. Сохраняем в кэш на 1 минуту, чтобы следующий запрос был HOT
	data, _ := json.Marshal(p)
	rdb.Set(ctx, key, data, time.Minute)

	return &p, "Холодный (из базы данных)"
}

func main() {
	// Инициализация (пропустим обработку ошибок для краткости)
	db, _ := database.ConnectDB()
	rdb := cache.NewRedisClient()
	ctx := context.Background()

	// НОВАЯ ЧАСТЬ: Прогреваем кэш при старте
	err := warmup.PreloadTopProducts(ctx, db, rdb)
	if err != nil {
		log.Fatalf("Warmup failed: %v", err)
	}

	// Теперь проверим запрос к ID=5 (который попал под прогрев)
	start := time.Now()
	p, status := getProduct(ctx, db, rdb, 5)
	fmt.Printf("[%s] ID: %d, Time: %v\n", status, p.ID, time.Since(start))

	// Демонстрация
	targetID := 10

	// Первый вызов — COLD

	p1, status1 := getProduct(ctx, db, rdb, targetID)
	fmt.Printf("[%s] ID: %d, Time: %v\n", status1, p1.ID, time.Since(start))

	// Второй вызов — HOT
	start = time.Now()
	p2, status2 := getProduct(ctx, db, rdb, targetID)
	fmt.Printf("[%s] ID: %d, Time: %v\n", status2, p2.ID, time.Since(start))
}
