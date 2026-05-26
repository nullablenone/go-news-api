package user

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,max=100"`
	Email    string `json:"email" binding:"required,email,max=150"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ToRegisterResponse(u User) RegisterResponse {
	return RegisterResponse{Name: u.Name, Email: u.Email}
}
