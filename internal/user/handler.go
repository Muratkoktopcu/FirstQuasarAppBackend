package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	repo *Repository
}
type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// writeHeader status kodu yazar ve sadece 1 kez çağrılabilir
// Header.Set http cevabının header kısmına bilgi eklemek veya değiştirmek için kullanılır
// writeJSON  HTTP cevabını JSON formatında döndürür
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// Register godoc
// @Summary     Kullanıcı kaydı
// @Description Yeni bir kullanıcı oluşturur.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       payload body     RegisterRequest true "Kayıt bilgileri"
// @Success     201     {object} RegisterResponse
// @Failure     400     {string} string "invalid json or validation error"
// @Failure     409     {string} string "email already registered"
// @Failure     500     {string} string "db error"
// @Router      /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_json"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Name == "" || req.Email == "" || len(req.Password) < 6 {
		http.Error(w, "name/email required, password >= 6 chars", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	exists, err := h.repo.EmailExists(ctx, req.Email)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "email already registered", http.StatusConflict)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "hash error", http.StatusInternalServerError)
		return
	}

	u := &User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	id, err := h.repo.Create(ctx, u)
	if err != nil {
		http.Error(w, "insert error", http.StatusInternalServerError)
		return
	}

	resp := RegisterResponse{ID: id, Name: u.Name, Email: u.Email}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
