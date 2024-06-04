package models

type LoginInput struct {
	Username string
	Password string
}

type LoginOutput struct {
	Success bool
	Message string
}
