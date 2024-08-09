package product

import (
	"net/http"

	"github.com/AhmadMirza2023/ecom/types"
	"github.com/AhmadMirza2023/ecom/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.HandleCreateProduct).Methods(http.MethodGet)
}

func (h *Handler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	utils.WriteJson(w, http.StatusOK, products)
}
