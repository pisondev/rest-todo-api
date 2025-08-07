package web

type UserAuthRequest struct {
	Username string `validate:"required,min=3,max=20" json:"username"`
	Password string `validate:"required,min=6,max=64" json:"password"`
}
