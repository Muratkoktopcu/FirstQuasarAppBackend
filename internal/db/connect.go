package db

import (
	"context"
	"log"  //loglama için kullanılır.Hataları ve bilgileri konsola yazdırmak için.
	"os"   //ortam değişkenlerini okumak için kullanılır.
	"time" //zaman aşımı ve timeout için kullanılır.

	"github.com/jackc/pgx/v5/pgxpool" //PostgreSQL için bağlantı havuzu sağlayan kütüphane. Havuz (pool), birden fazla bağlantıyı yöneterek performansı artırır.
	"github.com/joho/godotenv"        //.env dosyasını okuyup ortam değişkenlerine (os.Getenv) aktarmaya yarar.
)

/*
Connect() → Veritabanına bağlanmak için yazılmış fonksiyon.
Dönüş tipi *pgxpool.Pool → PostgreSQL bağlantı havuzu nesnesinin adresi.
Bu fonksiyon çağırıldığında her seferinde hazır ve test edilmiş bir veritabanı bağlantısı elde etmiş olursun.
*/
func Connect() *pgxpool.Pool {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("DATABASE_URL not set")
	}

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal("pgx parse config:", err)
	}

	// küçük bir pool ayarı
	cfg.MaxConns = 5

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal("pgx new pool:", err)
	}

	// test
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("db ping:", err)
	}
	return pool
}
