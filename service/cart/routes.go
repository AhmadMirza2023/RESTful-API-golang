package cart

import (
	"fmt"
	"net/http"

	"github.com/AhmadMirza2023/ecom/service/auth"
	"github.com/AhmadMirza2023/ecom/types"
	"github.com/AhmadMirza2023/ecom/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, productStore: productStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userId := auth.GetUserIDFromContext(r.Context())
	var cart types.CartCheckoutPayload
	// get JSON payload
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// get products
	productIDs, err := getCartItemsIDs(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	products, err := h.productStore.GetProductsByID(productIDs)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	orderId, totalPrice, err := h.createOrder(products, cart.Items, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"orderId":     orderId,
	})
}
