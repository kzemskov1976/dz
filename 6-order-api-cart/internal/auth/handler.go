package auth

import (
	"demo/order/6-order-api-cart/configs"
	"demo/order/6-order-api-cart/internal/user"
	"demo/order/6-order-api-cart/pkg/jwt"
	"demo/order/6-order-api-cart/pkg/req"
	"demo/order/6-order-api-cart/pkg/res"
	"net/http"
)

type AuthHandler struct {
	*user.UserRepository
	*configs.Config
}

type AuthHandlerDeps struct {
	*user.UserRepository
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		UserRepository: deps.UserRepository,
		Config:         deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/verify", handler.Verify())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[AuthRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sessionId, err := handler.UserRepository.GetSessionId(body.Phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := AuthResponse{
			SessionId: sessionId,
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *AuthHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[VerifyRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user, err := handler.UserRepository.VerifyCode(body.SessionId, body.Code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := jwt.NewJWT(handler.Config.Secret).Create(jwt.JWTData{
			Phone: user.Phone,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := VerifyResponse{
			Token: token,
		}
		res.Json(w, data, http.StatusOK)

	}
}
