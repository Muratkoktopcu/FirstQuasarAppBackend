package auth

import (
	"context"
	"errors"

	"myfirstbackend/internal/model/dto"
	"myfirstbackend/internal/repository/user"
)

var (
	ErrInvalidID = errors.New("invalid user id")
	ErrNotFound  = errors.New("user not found")
)

type ServiceProfile interface {
	// İş mantığı: id doğrula, repo'dan getir, DTO'ya map et
	GetProfile(ctx context.Context, id int64) (dto.UserProfileResponse, error)
}

type serviceProfile struct {
	repo user.Repository
}

func NewServiceProfile(repo user.Repository) ServiceProfile {
	return &serviceProfile{repo: repo}
}

func (s *serviceProfile) GetProfile(ctx context.Context, id int64) (dto.UserProfileResponse, error) {
	if id <= 0 {
		return dto.UserProfileResponse{}, ErrInvalidID
	}

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// repo hata türüne göre granular davranabilirsin.
		// örn: pgx.ErrNoRows → ErrNotFound
		return dto.UserProfileResponse{}, ErrNotFound
	}

	return dto.UserProfileResponse{
		ID:       u.ID,
		Email:    u.Email,
		FullName: u.FullName,
	}, nil
}
