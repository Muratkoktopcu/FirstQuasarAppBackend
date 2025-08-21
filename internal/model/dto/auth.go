package dto

type RegisterRequest struct {
	Email    string `json:"email" example:"alice@example.com"`
	Password string `json:"password" example:"P@ssw0rd!"`
	FullName string `json:"fullName" example:"Alice Doe"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"alice@example.com"`
	Password string `json:"password" example:"P@ssw0rd!"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserProfileResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
}
