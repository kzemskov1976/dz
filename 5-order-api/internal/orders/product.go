package orders

import (
	"demo/order/5-order-api/configs"
	"demo/order/5-order-api/pkg/middleware"
	"demo/order/5-order-api/pkg/req"
	"demo/order/5-order-api/pkg/res"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type ProductHandler struct {
	*ProductRepository
}

type ProductHandlerDeps struct {
	*configs.Config
	*ProductRepository
}

func (handler *ProductHandler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(middleware.PhoneKey).(string)
		if !ok {
			http.Error(w, "", http.StatusInternalServerError)
		}
		body, err := req.HandleBody[ProductCreateRequest](&w, r)
		if err != nil {
			res.Json(w, nil, http.StatusBadRequest)
			return
		}
		product := NewProduct(*body)
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
		product, err := handler.ProductRepository.UpdateProduct(&Product{
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

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		ProductRepository: deps.ProductRepository,
	}
	router.Handle("POST /product", middleware.Auth(handler.CreateProduct(), deps.Config))
	router.Handle("GET /product/{id}", middleware.Auth(handler.GetProduct(), deps.Config))
	router.Handle("PATCH /product/{id}", middleware.Auth(handler.UpdateProduct(), deps.Config))
	router.Handle("DELETE /product/{id}", middleware.Auth(handler.DeleteProduct(), deps.Config))
}
