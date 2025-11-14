package products

import (
	_json "encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	repo "github.com/jan-paulus/go-api/internal/adapters/sqlite/sqlc"
	"github.com/jan-paulus/go-api/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (handler *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := handler.service.ListProducts(r.Context())
	if err != nil {
		slog.Error("Failed to list products.", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, products)
}

func (handler *handler) FindProductByID(w http.ResponseWriter, r *http.Request) {
	product, err := handler.service.FindProductByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		slog.Warn("No product found with the given ID.", "error", err)
		http.Error(w, "No product found with the given ID.", http.StatusNotFound)
		return
	}

	json.Write(w, http.StatusOK, product)
}

func (handler *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	params := repo.CreateProductParams{}
	if err := _json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	product, err := handler.service.CreateProduct(r.Context(), params)
	if err != nil {
		slog.Warn("Failed to create product.", "error", err, "params", params)
		http.Error(w, "Failed to create product.", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, product)
}
