package warmup

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

// Продукт для внутреннего использования в warmup
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func PreloadTopProducts(ctx context.Context, db *pgx.Conn, rdb *redis.Client) error {
	log.Println("--- Starting Cache Warm-up ---")

	// Выбираем топ-10 товаров (например, самые дешевые или просто первые)
	rows, err := db.Query(ctx, "SELECT id, name, price FROM products LIMIT 10")
	if err != nil {
		return err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			continue
		}

		// Кладем в Redis
		key := fmt.Sprintf("product:%d", p.ID)
		data, _ := json.Marshal(p)
		rdb.Set(ctx, key, data, 10*time.Minute)
		count++
	}

	log.Printf("--- Warm-up complete: %d products cached ---\n", count)
	return nil
}
