package verify

import (
	"crypto/sha256"
	"demo/validation/3-validation-api/configs"
	"demo/validation/3-validation-api/pkg/email"
	"demo/validation/3-validation-api/pkg/res"
	"encoding/json"
	"fmt"
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

	h := sha256.New()
	h.Write([]byte(email))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	newItem := HashItem{
		Email: email,
		Hash:  hash,
	}

	fp := "verify.json"
	file, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var hashArr []HashItem

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if info.Size() != 0 {
		if err := json.NewDecoder(file).Decode(&hashArr); err != nil {
			return nil, err
		}
		if len(hashArr) > 0 {
			for _, item := range hashArr {
				if item.Email == email {
					return &newItem, nil
				}
			}
		}
		file.Truncate(0)
		file.Seek(0, 0)
	}

	hashArr = append(hashArr, newItem)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(hashArr); err != nil {
		return nil, err
	}

	return &newItem, nil
}

func CheckHash(hash string) bool {
	fp := "verify.json"
	file, err := os.OpenFile(fp, os.O_RDONLY, 0644)
	if err != nil {
		return false
	}
	defer file.Close()

	var hashArr []HashItem

	if err := json.NewDecoder(file).Decode(&hashArr); err != nil {
		return false
	}

	for _, item := range hashArr {
		if item.Hash == hash {
			return true
		}
	}
	return false
}
