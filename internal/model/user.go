package model

type (
	RegisterRequest struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	AuthResponse struct {
		Email       string `json:"email"`
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	}
	LoginRequest struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		SecretKey string `json:"secretKey"`
	}

	FilterFindUser struct {
		Email *string
		ID    []uint64
	}
)
