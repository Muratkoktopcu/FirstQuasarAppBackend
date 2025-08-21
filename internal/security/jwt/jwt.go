package jwt

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// bu paket iki temel iş yapıyor.
type Config struct {
	Secret          string //HS256 ile imzalama anahtarı.string tutulma sebebi time durationa parse etmesi kolay olsun diye.
	AccessTokenTTL  string // "15m"
	RefreshTokenTTL string // "720h"
	Issuer          string // token'ı üreten uygulamanın adı
}

type Service interface {
	GenerateTokens(userID int64, email string) (access, refresh string, err error)
	Parse(token string) (*jwtv5.Token, jwtv5.MapClaims, error)
}

// configten gelen değerlerin run time karşılıkları
type service struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	issuer     string
}

func New(cfg Config) Service {
	acc, _ := time.ParseDuration(cfg.AccessTokenTTL)
	ref, _ := time.ParseDuration(cfg.RefreshTokenTTL)
	return &service{
		secret:     []byte(cfg.Secret),
		accessTTL:  acc,
		refreshTTL: ref,
		issuer:     cfg.Issuer,
	}
}

/*
Neden iki token?
Access token kısa ömürlüdür; her istekte Authorization: Bearer ile gelir.
Refresh token uzun ömürlüdür; access süresi bittiğinde yeni access üretmek için kullanılır (genellikle sunucu tarafında bir “refresh endpoint” ile)
*/
func (s *service) GenerateTokens(userID int64, email string) (string, string, error) {
	now := time.Now()
	//claims, JWT token'ın içeriğini tutar.
	accessClaims := jwtv5.MapClaims{
		"sub":   userID, // subject, yani kullanıcı ID'si
		"email": email,
		"iss":   s.issuer,                    //kim üretti
		"iat":   now.Unix(),                  //ne zaman üretildi
		"exp":   now.Add(s.accessTTL).Unix(), //ne zaman süresi dolacak
		"type":  "access",                    // token türü
	}
	//HS256 ile imzalanmış JWT token oluştururuz
	accessToken := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, accessClaims)
	//secret ile imzalanır,string'e çevrilir
	accStr, err := accessToken.SignedString(s.secret)
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwtv5.MapClaims{
		"sub":   userID,
		"email": email,
		"iss":   s.issuer,
		"iat":   now.Unix(),
		"exp":   now.Add(s.refreshTTL).Unix(),
		"type":  "refresh",
	}
	refreshToken := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, refreshClaims)
	refStr, err := refreshToken.SignedString(s.secret)
	if err != nil {
		return "", "", err
	}

	return accStr, refStr, nil
}

// Parse token'ı alır, doğrular ve içindeki bilgileri çıkarır.
// Eğer token geçerli değilse hata döner.
func (s *service) Parse(tokenStr string) (*jwtv5.Token, jwtv5.MapClaims, error) {
	/**jwtv5.Token
	  Parse edilmiş token nesnesi.
	  İçinde header, method, claims gibi tüm token yapısı bulunur.*/
	/**jwtv5.MapClaims
	  Token’ın payload kısmındaki claims (id, email, expireTime gibi key-value bilgiler).
	  Tipi: map[string]interface{} benzeri.
	  Mesela claims["sub"] veya claims["exp"] gibi değerleri buradan alırsın.*/

	//NewParser fonksiyonu, JWT token’larını okumak ve doğrulamak için bir parser nesnesi oluşturur.
	parser := jwtv5.NewParser(jwtv5.WithValidMethods([]string{jwtv5.SigningMethodHS256.Alg()})) //alg saldırılarına karşı koruma.Token sadece HS256 ile imzalanmışsa kabul edilir.alg : none veya farklı bir algoritma ise hata döner.
	claims := jwtv5.MapClaims{}                                                                 //map[string]interface{} olarak tanımlanır.
	tkn, err := parser.ParseWithClaims(tokenStr, claims, func(token *jwtv5.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, nil, err
	}
	if !tkn.Valid {
		return nil, nil, errors.New("invalid token")
	}
	return tkn, claims, nil
}
