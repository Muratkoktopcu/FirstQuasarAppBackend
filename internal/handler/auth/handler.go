package auth

import (
	"encoding/json"
	"net/http"

	"myfirstbackend/internal/model/dto"
	"myfirstbackend/internal/service/auth"
)

type Handler struct {
	svc auth.Service
}

func NewHandler(svc auth.Service) *Handler {
	return &Handler{svc: svc}
}

// Register godoc
// @Summary      Register user
// @Description  Yeni kullanıcı kaydı (email, password, fullName)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.RegisterRequest  true  "Register payload"
// @Success      201      {object}  dto.AuthResponse
// @Failure      400      {string}  string
// @Router       /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var in dto.RegisterRequest
	//HTTP body’sindeki JSON’u in struct’ına map etmeye çalışır.
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	//r.Context(): HTTP request’in context’i. Timeout, cancellation, trace gibi şeyler bu context’te tutulur. DB işlemleri de iptal edilebilir hale gelir.
	out, err := h.svc.Register(r.Context(), in) //
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated) //response’un durum kodunu 201 Created olarak ayarlar.
	//out: DTO (Data Transfer Object) olarak tanımlanmış bir yapıdır.
	_ = json.NewEncoder(w).Encode(out)
}

// Login godoc
// @Summary      Login
// @Description  Email + password ile giriş
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.LoginRequest  true  "Login payload"
// @Success      200      {object}  dto.AuthResponse
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
	_ = json.NewEncoder(w).Encode(out)
}
