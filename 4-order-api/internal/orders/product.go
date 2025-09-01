package orders

import (
	"demo/order/4-order-api/pkg/req"
	"demo/order/4-order-api/pkg/res"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type ProductHandler struct {
	*ProductRepository
}

type ProductHandlerDeps struct {
	*ProductRepository
}

func (handler *ProductHandler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	router.HandleFunc("POST /product", handler.CreateProduct())
	router.HandleFunc("GET /product/{id}", handler.GetProduct())
	router.HandleFunc("PATCH /product/{id}", handler.UpdateProduct())
	router.HandleFunc("DELETE /product/{id}", handler.DeleteProduct())
}
