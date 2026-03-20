package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"gitverse.ru/pavlov-systems/cache-lifecycle-go/internal/cache"
	"gitverse.ru/pavlov-systems/cache-lifecycle-go/internal/database"
	"gitverse.ru/pavlov-systems/cache-lifecycle-go/internal/warmup"
)

// Структура продукта (должна совпадать с той, что в БД)
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var (
	db  *pgx.Conn
	rdb *redis.Client
)

func main() {
	var err error
	ctx := context.Background()

	// 1. Подключения
	db, err = database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	rdb = cache.NewRedisClient()

	// 2. ПРОГРЕВ (Warm-up) при старте
	log.Println("Initializing application...")
	warmup.PreloadTopProducts(ctx, db, rdb)

	// 3. Маршрут API
	http.HandleFunc("/product/", productHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Извлекаем ID из URL (например, /product/5)
	idStr := r.URL.Path[len("/product/"):]
	id, _ := strconv.Atoi(idStr)

	start := time.Now()

	// Наша логика кэширования
	p, status := getProduct(ctx, db, rdb, id)

	duration := time.Since(start)

	// Формируем ответ
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache-Status", status) // Показываем статус в заголовках
	w.Header().Set("X-Response-Time", duration.String())

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":   p,
		"status": status,
		"time":   duration.String(),
	})
}

// Та самая функция getProduct (убедитесь, что она объявлена в main.go или вынесена)
func getProduct(ctx context.Context, db *pgx.Conn, rdb *redis.Client, id int) (*Product, string) {
	key := fmt.Sprintf("product:%d", id)

	// Проверка кэша (HOT)
	val, err := rdb.Get(ctx, key).Result()
	if err == nil {
		var p Product
		json.Unmarshal([]byte(val), &p)
		return &p, "HOT"
	}

	// Имитация тяжелой БД (COLD)
	time.Sleep(500 * time.Millisecond)

	var p Product
	err = db.QueryRow(ctx, "SELECT id, name, price FROM products WHERE id=$1", id).
		Scan(&p.ID, &p.Name, &p.Price)

	if err != nil {
		return nil, "NOT_FOUND"
	}

	// Сохранение в кэш
	data, _ := json.Marshal(p)
	rdb.Set(ctx, key, data, 5*time.Minute)

	return &p, "COLD"
}
