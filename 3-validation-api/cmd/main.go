package main

import (
	"demo/validation/3-validation-api/configs"
	"demo/validation/internal/verify"
	"flag"
	"fmt"
	"net/http"
)

func main() {
	port := flag.Int("port", 8081, "Listenting port")
	flag.Parse()

	conf := configs.LoadConfig()

	router := http.NewServeMux()
	verify.NewVerifyHandler(router, verify.VerifyHandlerDeps{
		Config: conf,
	})

	addr := fmt.Sprintf(":%d", *port)

	server := http.Server{
		Addr:    addr,
		Handler: router,
	}
	fmt.Printf("Server listening on port %d\n", *port)
	server.ListenAndServe()
}
