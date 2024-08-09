package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/AhmadMirza2023/ecom/service/cart"
	"github.com/AhmadMirza2023/ecom/service/order"
	"github.com/AhmadMirza2023/ecom/service/product"
	"github.com/AhmadMirza2023/ecom/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (a *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(a.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	productStore := product.NewStore(a.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(a.db)

	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	log.Println("Listening on port", a.addr)

	return http.ListenAndServe(a.addr, router)
}
