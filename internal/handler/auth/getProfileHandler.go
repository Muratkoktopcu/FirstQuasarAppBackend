package auth

import (
	"encoding/json"
	"net/http"

	"myfirstbackend/internal/http/middleware"
	usersvc "myfirstbackend/internal/service/auth"
)

type HandlerGetProfile struct {
	svc usersvc.ServiceProfile
}

func NewHandlerGetProfile(svc usersvc.ServiceProfile) *HandlerGetProfile {
	return &HandlerGetProfile{svc: svc}
}

func writeJSONProfile(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Me godoc
// @Summary      Kullanıcı profilini getir
// @Description  JWT ile doğrulanmış kullanıcının profilini döndürür
// @Tags         profile
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  dto.UserProfileResponse
// @Failure      400  {string}  string "invalid user id"
// @Failure      401  {string}  string "unauthorized"
// @Failure      404  {string}  string "user not found"
// @Failure      500  {string}  string "internal error"
// @Router       /profile [get]
func (h *HandlerGetProfile) Me(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserIDFromCtx(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	out, err := h.svc.GetProfile(r.Context(), uid)
	if err != nil {
		switch err {
		case usersvc.ErrInvalidID:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usersvc.ErrNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	writeJSONProfile(w, http.StatusOK, out)
}
