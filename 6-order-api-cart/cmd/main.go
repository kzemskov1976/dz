package main

import (
	"demo/order/6-order-api-cart/configs"
	"demo/order/6-order-api-cart/internal/auth"
	"demo/order/6-order-api-cart/internal/orders"
	"demo/order/6-order-api-cart/internal/user"
	"demo/order/6-order-api-cart/internal/verify"
	"demo/order/6-order-api-cart/pkg/db"
	"demo/order/6-order-api-cart/pkg/middleware"
	"flag"
	"fmt"
	"net/http"
)

func main() {
	port := flag.Int("port", 8081, "Listenting port")
	flag.Parse()

	conf := configs.LoadConfig()
	db := db.NewDB(conf)
	repoProduct := orders.NewProductRepository(db)
	repoUser := user.NewUserRepository(db)

	router := http.NewServeMux()

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		UserRepository: repoUser,
		Config:         conf,
	})
	verify.NewVerifyHandler(router, verify.VerifyHandlerDeps{
		Config: conf,
	})
	orders.NewProductHandler(router, orders.ProductHandlerDeps{
		ProductRepository: repoProduct,
		Config:            conf,
		IntUserRepository: repoUser,
	})

	addr := fmt.Sprintf(":%d", *port)

	server := http.Server{
		Addr:    addr,
		Handler: middleware.Logging(router),
	}
	fmt.Printf("Server listening on port %d\n", *port)
	server.ListenAndServe()
}
