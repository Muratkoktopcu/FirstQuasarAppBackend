package auth

import (
	"context" //request ömrü, iptal ve timeout yönetimi için.
	"errors"

	domain "myfirstbackend/internal/domain/user" //domain modeli (veritabanı gerçeğini temsil eden “User” nesnesi).
	"myfirstbackend/internal/model/dto"          //istek/yanıt modelleri (API ile konuşan dış katmana özel).
	"myfirstbackend/internal/repository/user"    //Repository arayüzü (veri erişim soyutlaması).
	"myfirstbackend/internal/security/jwt"       //JWT üretimi yapan servis arayüzü.

	"golang.org/x/crypto/bcrypt" //bcrypt kütüphanesi, şifreleri güvenli bir şekilde hashlemek için kullanılır.
)

// Dışarıya sunduğun sözleşme. Controller/Handler katmanı bu arayüzle konuşur.
// Geri dönüşlerde dto.AuthResponse (Access/Refresh token’lar) ve error var.
type Service interface {
	Register(ctx context.Context, in dto.RegisterRequest) (dto.AuthResponse, error)
	Login(ctx context.Context, in dto.LoginRequest) (dto.AuthResponse, error)
}

/*
service: İç implementasyon.
repo: Kullanıcıyı DB’de aramak/oluşturmak için.
tokens: JWT üretmek için.
*/
type service struct {
	repo   user.Repository
	tokens jwt.Service
}

// Yeni servis oluşturma fonksiyonu. Dışarıdan repo ve tokens alır.

func New(repo user.Repository, tokens jwt.Service) Service {
	return &service{repo: repo, tokens: tokens} //tokens a tokens servisi verir.
}

// Register: Yeni kullanıcı kaydı yapar.
func (s *service) Register(ctx context.Context, in dto.RegisterRequest) (dto.AuthResponse, error) {
	// Gerekli alanların dolu olup olmadığını kontrol et.
	if in.Email == "" || in.Password == "" || in.FullName == "" {
		return dto.AuthResponse{}, errors.New("missing fields")
	}

	// email already?
	if _, err := s.repo.FindByEmail(ctx, in.Email); err == nil {
		return dto.AuthResponse{}, errors.New("email already exists")
	}
	// Şifreyi hashle.
	//parametreyi byte olarak alma sebebi: bcrypt kütüphanesi byte dizisi ile çalışır.
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	// Yeni kullanıcı nesnesi oluştur.
	// ID veritabanı tarafından otomatik atanacak.
	u := &domain.User{
		Email:    in.Email,
		Password: string(hash),
		FullName: in.FullName,
	}
	// Kullanıcıyı veritabanına kaydet.
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
