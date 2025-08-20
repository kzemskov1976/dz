package verify

import (
	"demo/validation/3-validation-api/configs"
	"demo/validation/3-validation-api/email"
	"fmt"
	"net/http"
)

type VerifyHandler struct {
	*configs.Config
}

type VerifyHandlerDeps struct {
	*configs.Config
}

func (handler *VerifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// fmt.Println("send email")
		email.SendEmail(handler.Email, handler.Password, handler.Address)
	}
}

func (handler *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("verify")
	}
}

func NewVerifyHandler(router *http.ServeMux, deps VerifyHandlerDeps) {
	handler := &VerifyHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /send", handler.Send())
	router.HandleFunc("GET /verify/{hash}", handler.Verify())
}
