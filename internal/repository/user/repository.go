package user

import (
	"context"
	//"errors"
	//"time"

	domain "myfirstbackend/internal/domain/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, u *domain.User) (int64, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type pgRepository struct {
	pool *pgxpool.Pool
}

/*
Bu fonksiyon, dışarıdan pool alıp pgRepository döndürüyor.
Dönüş tipi Repository interface olduğu için, dışarıya sadece interface görünüyor (soyutlama).
Bu da katmanlar arası bağımlılığı azaltır (clean architecture mantığı).
*/
func NewPgRepository(pool *pgxpool.Pool) Repository {
	return &pgRepository{pool: pool}
}

func (r *pgRepository) Create(ctx context.Context, u *domain.User) (int64, error) {
	q := `
		INSERT INTO users (email, password, full_name, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id;
	`
	var id int64
	if err := r.pool.QueryRow(ctx, q, u.Email, u.Password, u.FullName).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *pgRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	q := `
		SELECT id, email, password, COALESCE(full_name,''), created_at, updated_at
		FROM users WHERE email=$1 LIMIT 1;
	`
	row := r.pool.QueryRow(ctx, q, email)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.FullName, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

// (Not: tabloyu oluşturmadıysan örnek SQL)
// CREATE TABLE IF NOT EXISTS users(
//   id BIGSERIAL PRIMARY KEY,
//   email TEXT UNIQUE NOT NULL,
//   password TEXT NOT NULL,
//   full_name TEXT NOT NULL,
//   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
//   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
// );
