package orders

import (
	"demo/order/6-order-api-cart/configs"
	di "demo/order/6-order-api-cart/pkg/interfaces"
	"demo/order/6-order-api-cart/pkg/middleware"
	"net/http"
)

type ProductHandler struct {
	*ProductRepository
	di.IntUserRepository
}

type ProductHandlerDeps struct {
	*configs.Config
	*ProductRepository
	di.IntUserRepository
}

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		ProductRepository: deps.ProductRepository,
		IntUserRepository: deps.IntUserRepository,
	}
	router.Handle("POST /product", middleware.Auth(handler.AddProduct(), deps.Config))
	router.Handle("GET /product/{id}", middleware.Auth(handler.GetProduct(), deps.Config))
	router.Handle("PATCH /product/{id}", middleware.Auth(handler.UpdateProduct(), deps.Config))
	router.Handle("DELETE /product/{id}", middleware.Auth(handler.DeleteProduct(), deps.Config))

	router.Handle("POST /order", middleware.Auth(handler.CreateOrder(), deps.Config))
	router.Handle("GET /order/{id}", middleware.Auth(handler.GetOrder(), deps.Config))
	router.Handle("GET /my-orders", middleware.Auth(handler.GetMyOrders(), deps.Config))

}
