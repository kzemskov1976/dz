package orders

import (
	"demo/order/6-order-api-cart/pkg/db"
	"demo/order/6-order-api-cart/pkg/middleware"
	"demo/order/6-order-api-cart/pkg/req"
	"demo/order/6-order-api-cart/pkg/res"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func (handler *ProductHandler) AddProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(middleware.PhoneKey).(string)
		if !ok {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		body, err := req.HandleBody[ProductCreateRequest](&w, r)
		if err != nil {
			res.Json(w, nil, http.StatusBadRequest)
			return
		}
		product := db.NewProduct(body.Name, body.Description)
		Createdproduct, err := handler.ProductRepository.AddProduct(product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, Createdproduct, http.StatusCreated)
	}
}

func (handler *ProductHandler) GetProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(middleware.PhoneKey).(string)
		if !ok {
			http.Error(w, "", http.StatusInternalServerError)
		}
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		product, err := handler.GetProductById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		res.Json(w, product, http.StatusOK)
	}
}

func (handler *ProductHandler) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(middleware.PhoneKey).(string)
		if !ok {
			http.Error(w, "", http.StatusInternalServerError)
		}
		body, err := req.HandleBody[ProductUpdateRequest](&w, r)
		if err != nil {
			res.Json(w, nil, http.StatusBadRequest)
			return
		}
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		product, err := handler.ProductRepository.UpdateProduct(&db.Product{
			Model:       gorm.Model{ID: uint(id)},
			Name:        body.Name,
			Description: body.Description,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, product, http.StatusOK)
	}
}

func (handler *ProductHandler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(middleware.PhoneKey).(string)
		if !ok {
			http.Error(w, "", http.StatusInternalServerError)
		}
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = handler.GetProductById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		err = handler.DeleteProductById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, nil, http.StatusOK)
	}
}
