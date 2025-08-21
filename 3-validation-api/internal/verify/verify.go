package verify

import (
	"crypto/rand"
	"demo/validation/3-validation-api/configs"
	"demo/validation/3-validation-api/pkg/email"
	"demo/validation/3-validation-api/pkg/res"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
)

type VerifyHandler struct {
	*configs.Config
}

type VerifyHandlerDeps struct {
	*configs.Config
}

func (handler *VerifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var payload EmailVerify
		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}
		hash, err := SaveHash(payload.Email)
		if err != nil {
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// fmt.Println(hash)
		err = email.SendEmail(handler.Email, handler.Password, handler.Address, hash.Email, hash.Hash)
		if err != nil {
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, "Письмо для подтверждения отправлено", http.StatusOK)
	}
}

func (handler *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		result := CheckHash(hash)
		res.Json(w, result, http.StatusOK)
	}
}

func NewVerifyHandler(router *http.ServeMux, deps VerifyHandlerDeps) {
	handler := &VerifyHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /send", handler.Send())
	router.HandleFunc("GET /verify/{hash}", handler.Verify())
}

func SaveHash(email string) (*HashItem, error) {

	hash, _ := GenerateRandomHash(16)
	newItem := HashItem{
		Email: email,
		Hash:  hash,
	}

	fp := "verify.json"
	file, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(newItem); err != nil {
		return nil, err
	}

	return &newItem, nil
}

func CheckHash(hash string) bool {
	result := false
	fp := "verify.json"
	file, err := os.OpenFile(fp, os.O_RDWR, 0644)
	if err != nil {
		result = false
	}
	defer file.Close()

	var data HashItem

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		result = false
	}

	if data.Hash == hash {
		result = true
	}

	file.Truncate(0)
	file.Seek(0, 0)
	return result
}

func GenerateRandomHash(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
