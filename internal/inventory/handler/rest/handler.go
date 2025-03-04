package rest

import (
	"encoding/json"
	"net/http"

	"github.com/kijevigombooc/inventory-manager/internal/inventory/handler/dto"
	"github.com/kijevigombooc/inventory-manager/internal/inventory/service"
)

func NewInventoryHandler(service service.Service) *inventoryHandler {
	return &inventoryHandler{service: service}
}

type inventoryHandler struct {
	service service.Service
}

func (h *inventoryHandler) RegisterRoutes(serveMux *http.ServeMux) {
	serveMux.HandleFunc("GET /warehouses", h.getWarehouses)
	serveMux.HandleFunc("POST /warehouses", h.createWarehouse)
	serveMux.HandleFunc("POST /insertProducts", h.insertProducts)
	serveMux.HandleFunc("POST /removeProducts", h.removeProducts)
}

func (h *inventoryHandler) getWarehouses(w http.ResponseWriter, r *http.Request) {
	warehouses, err := h.service.GetWarehouses()
	if err != nil {
		writeErrorMessageJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, warehouses, http.StatusOK)
}

func (h *inventoryHandler) createWarehouse(w http.ResponseWriter, r *http.Request) {
	warehouse := dto.Warehouse{}
	if err := json.NewDecoder(r.Body).Decode(&warehouse); err != nil {
		writeErrorMessageJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.CreateWarehouse(warehouse); err != nil {
		// TODO: return different error message if warehouse already exists
		writeErrorMessageJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, warehouse, http.StatusCreated)
}

func (h *inventoryHandler) insertProducts(w http.ResponseWriter, r *http.Request) {
	var req dto.InsertProductsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorMessageJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.ParseProduct(); err != nil {
		writeErrorMessageJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.InsertProducts(req.WarehouseName, req.ParsedProduct, req.Quantity); err != nil {
		writeErrorMessageJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, req, http.StatusOK)
}

func (h *inventoryHandler) removeProducts(w http.ResponseWriter, r *http.Request) {
	var req dto.RemoveProductsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorMessageJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.RemoveProducts(req.WarehouseName, req.Sku, req.Quantity); err != nil {
		writeErrorMessageJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, req, http.StatusOK)
}

func writeErrorMessageJSON(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error": "` + message + `"}`))
}

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.WriteHeader(statusCode)
	w.Write(response)
	return nil
}
