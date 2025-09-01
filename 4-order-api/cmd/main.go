package main

import (
	"demo/order/4-order-api/configs"
	"demo/order/4-order-api/internal/orders"
	"demo/order/4-order-api/internal/verify"
	"demo/order/4-order-api/pkg/db"
	"demo/order/4-order-api/pkg/middleware"
	"flag"
	"fmt"
	"net/http"
)

func main() {
	port := flag.Int("port", 8081, "Listenting port")
	flag.Parse()

	conf := configs.LoadConfig()
	db := db.NewDB(conf)
	repo := orders.NewProductRepository(db)

	router := http.NewServeMux()

	verify.NewVerifyHandler(router, verify.VerifyHandlerDeps{
		Config: conf,
	})
	orders.NewProductHandler(router, orders.ProductHandlerDeps{
		ProductRepository: repo,
	})

	addr := fmt.Sprintf(":%d", *port)

	server := http.Server{
		Addr:    addr,
		Handler: middleware.Logging(router),
	}
	fmt.Printf("Server listening on port %d\n", *port)
	server.ListenAndServe()
}
