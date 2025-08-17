package auth

import (
	"context"
	"errors"

	domain "myfirstbackend/internal/domain/user"
	"myfirstbackend/internal/model/dto"
	"myfirstbackend/internal/repository/user"
	"myfirstbackend/internal/security/jwt"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, in dto.RegisterRequest) (dto.AuthResponse, error)
	Login(ctx context.Context, in dto.LoginRequest) (dto.AuthResponse, error)
}

type service struct {
	repo   user.Repository
	tokens jwt.Service
}

func New(repo user.Repository, tokens jwt.Service) Service {
	return &service{repo: repo, tokens: tokens}
}

func (s *service) Register(ctx context.Context, in dto.RegisterRequest) (dto.AuthResponse, error) {
	if in.Email == "" || in.Password == "" || in.FullName == "" {
		return dto.AuthResponse{}, errors.New("missing fields")
	}

	// email already?
	if _, err := s.repo.FindByEmail(ctx, in.Email); err == nil {
		return dto.AuthResponse{}, errors.New("email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	u := &domain.User{
		Email:    in.Email,
		Password: string(hash),
		FullName: in.FullName,
	}

	id, err := s.repo.Create(ctx, u)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	acc, ref, err := s.tokens.GenerateTokens(id, in.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	return dto.AuthResponse{AccessToken: acc, RefreshToken: ref}, nil
}

func (s *service) Login(ctx context.Context, in dto.LoginRequest) (dto.AuthResponse, error) {
	u, err := s.repo.FindByEmail(ctx, in.Email)
	if err != nil {
		return dto.AuthResponse{}, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(in.Password)); err != nil {
		return dto.AuthResponse{}, errors.New("invalid credentials")
	}
	acc, ref, err := s.tokens.GenerateTokens(u.ID, u.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	return dto.AuthResponse{AccessToken: acc, RefreshToken: ref}, nil
}
