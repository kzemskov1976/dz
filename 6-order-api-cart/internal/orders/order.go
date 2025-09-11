package orders

import (
	"demo/order/6-order-api-cart/pkg/middleware"
	"demo/order/6-order-api-cart/pkg/req"
	"demo/order/6-order-api-cart/pkg/res"
	"net/http"
	"strconv"
)

func (handler *ProductHandler) CreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.PhoneKey).(string)
		if !ok {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		body, err := req.HandleBody[OrderCreateRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user, err := handler.IntUserRepository.FindByPhone(phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		order, err := handler.ProductRepository.AddOrder(user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, product := range body.Products {
			p, err := handler.ProductRepository.GetProductByName(product)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = handler.ProductRepository.AddProductToOrder(order, p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		res.Json(w, OrderCreateResponse{OrderId: order.ID}, http.StatusCreated)
	}
}

func (handler *ProductHandler) GetMyOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.PhoneKey).(string)
		if !ok {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		user, err := handler.IntUserRepository.FindByPhone(phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orders, err := handler.ProductRepository.GetOrdersByUserId(user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}
		data := MyOrdersResponse{}
		for _, o := range *orders {
			data.OrderId = append(data.OrderId, o.ID)
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *ProductHandler) GetOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.PhoneKey).(string)
		if !ok {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		user, err := handler.IntUserRepository.FindByPhone(phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		order, err := handler.ProductRepository.GetOrderByOrderId(user.ID, uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data := OrderResponse{
			OrderId: order.ID,
		}
		for _, item := range order.Products {
			data.Products = append(data.Products, item.Name)
		}
		res.Json(w, data, http.StatusOK)
	}
}
