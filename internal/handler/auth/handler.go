package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"myfirstbackend/internal/model/dto"
	"myfirstbackend/internal/security/jwt"
	authsvc "myfirstbackend/internal/service/auth"
)

type Handler struct {
	svc authsvc.Service
	jwt jwt.Service
}

func NewHandler(svc authsvc.Service, jwtSvc jwt.Service) *Handler {
	return &Handler{svc: svc, jwt: jwtSvc}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Access token’ı header’a, expiresIn’i de X-Expires-In header’ına yazar.
func (h *Handler) setAccessHeaders(w http.ResponseWriter, accessToken string) {
	w.Header().Set("Authorization", "Bearer "+accessToken)
	// exp'i token'dan okuyup kalan süreyi hesaplayalım:
	_, claims, err := h.jwt.Parse(accessToken)
	log.Printf("Parsed access token in setAccessHeaders: %v", claims)
	if err == nil {
		if expf, ok := claims["exp"].(float64); ok {
			exp := time.Unix(int64(expf), 0).UTC()
			remaining := time.Until(exp)
			if remaining < 0 {
				remaining = 0
			}
			w.Header().Set("X-Expires-In", strconv.FormatInt(int64(remaining.Seconds()), 10))
		}
	}
	// Bilgilendirici; istersen kaldır:
	w.Header().Set("X-Token-Type", "access")
}

// Refresh token’ı HttpOnly cookie olarak yazar.
func (h *Handler) setRefreshCookie(w http.ResponseWriter, refreshToken string) {
	_, claims, err := h.jwt.Parse(refreshToken)
	log.Printf("Parsed refresh token in setRefreshCookie: %v", claims)
	var exp time.Time
	if err == nil {
		if unix, ok := claims["exp"].(float64); ok {
			exp = time.Unix(int64(unix), 0).UTC()
		}
	}

	c := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/", // sadece refresh çağrısında gönderilsin
		HttpOnly: true,
		Secure:   true,                  // prod’da true; local HTTP geliştirirken false yapabilirsin
		SameSite: http.SameSiteNoneMode, // SPA farklı origin ise gerekli
	}
	if !exp.IsZero() {
		c.Expires = exp
	}
	http.SetCookie(w, c)
}

// Refresh cookie’yi temizler (logout).
func (h *Handler) clearRefreshCookie(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}
	http.SetCookie(w, c)
}

// ---- Handlers ----

// @Summary      Register user
// @Description  Yeni kullanıcı kaydı (email, password, fullName)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.RegisterRequest  true  "Register payload"
// @Success      201      {string}  string  "Access token header’da, refresh HttpOnly cookie’de"
// @Failure      400      {string}  string
// @Router       /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var in dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	out, err := h.svc.Register(r.Context(), in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.setRefreshCookie(w, out.RefreshToken) // refresh → cookie
	h.setAccessHeaders(w, out.AccessToken)  // access → headerlar

	// Body boş/çok küçük tutuyoruz; istersen kısa bir mesaj dön:
	writeJSON(w, http.StatusCreated, map[string]string{"status": "registered"})
}

// @Summary      Login
// @Description  Email + password ile giriş
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.LoginRequest  true  "Login payload"
// @Success      200      {string}  string  "Access token header’da, refresh HttpOnly cookie’de"
// @Failure      401      {string}  string
// @Router       /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var in dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	out, err := h.svc.Login(r.Context(), in)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	h.setRefreshCookie(w, out.RefreshToken)
	h.setAccessHeaders(w, out.AccessToken)

	writeJSON(w, http.StatusOK, map[string]string{"status": "logged_in"})
}

// @Summary      Refresh access token
// @Description  HttpOnly cookie’deki refresh token ile yeni access token üretir. (Opsiyonel rotation)
// @Tags         Auth
// @Produce      json
// @Success      200  {string}  string  "Access token header’da, refresh cookie opsiyonel yenilenir"
// @Failure      401  {string}  string
// @Router       /auth/refresh [post]
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("refresh_token")
	if err != nil || c.Value == "" {
		http.Error(w, "missing refresh token", http.StatusUnauthorized)
		return
	}
	_, claims, err := h.jwt.Parse(c.Value)
	log.Printf("Parsed refresh token in Refresh: %v", claims)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}
	if typ, _ := claims["type"].(string); typ != "refresh" {
		http.Error(w, "invalid token type", http.StatusUnauthorized)
		return
	}

	v, ok := claims["sub"]
	if !ok {
		http.Error(w, "missing sub", http.StatusUnauthorized)
		log.Printf("yeni versiyonda bi bak bakam: %v", v)
		return
	}
	//log.Printf("refresh fonk içerisindeki sub: %s", subStr)
	email, _ := claims["email"].(string)
	var userID int64
	switch vv := v.(type) {
	case string:
		id, err := strconv.ParseInt(vv, 10, 64)
		if err != nil {
			http.Error(w, "bad sub", http.StatusUnauthorized)
			return
		}
		userID = id
	case float64:
		userID = int64(vv)
	case json.Number: // ihtimale karşı
		id, err := vv.Int64()
		if err != nil {
			http.Error(w, "bad sub", http.StatusUnauthorized)
			return
		}
		userID = id
	default:
		http.Error(w, "invalid sub type", http.StatusUnauthorized)
		return
	}

	// (Opsiyonel) DB’de jti kontrolü & rotation yapılabilir (önerilir)
	// Şimdilik sadece yeni tokenlar üretelim:
	// Burada service yoksa jwtSvc.GenerateTokens da kullanabilirdik; fakat
	// rotation ileride service layer’a eklenirse oradan yürütmek mantıklı olur.
	//userID, _ := strconv.ParseInt(subStr, 10, 64)
	access, refresh, err := h.jwt.GenerateTokens(userID, email)
	if err != nil {
		// Eğer service’de böyle bir fonk yoksa: service değişmeden jwtSvc.GenerateTokens ile üret:
		// userID, _ := strconv.ParseInt(subStr, 10, 64)
		// access, refresh, err = h.jwt.GenerateTokens(userID, email)
		http.Error(w, "could not issue tokens", http.StatusInternalServerError)
		return
	}

	// Yeni refresh’i cookie’ye yaz (rotation yapıyorsan mutlaka)
	h.setRefreshCookie(w, refresh)
	// Yeni access’i header’a yaz
	h.setAccessHeaders(w, access)

	writeJSON(w, http.StatusOK, map[string]string{"status": "refreshed"})
}

// @Summary      Logout
// @Description  Refresh cookie’yi temizler (logout)
// @Tags         Auth
// @Produce      json
// @Success      200  {string}  string
// @Router       /auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	h.clearRefreshCookie(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "logged_out"})
}
