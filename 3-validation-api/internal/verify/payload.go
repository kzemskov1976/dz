package verify

type EmailVerify struct {
	Email string `json:"email" validate:"required,email"`
}

type HashItem struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
}
