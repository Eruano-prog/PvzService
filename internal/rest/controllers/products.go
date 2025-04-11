package controllers

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/rest"
	"AvitoPvz/internal/rest/middleware"
	"encoding/json"
	"log/slog"
	"net/http"
)

type ProductController struct {
	log            *slog.Logger
	productService rest.ProductService
}

func NewProductController(log *slog.Logger, productService rest.ProductService) *ProductController {
	return &ProductController{
		log:            log,
		productService: productService,
	}
}

func (p *ProductController) Register(mux *http.ServeMux, tokenVerifier middleware.Verifier) {
	mux.Handle("POST /products", middleware.AuthMiddleware(http.HandlerFunc(p.addProductHandler), tokenVerifier, p.log, middleware.EmployeeOnly))
}

func (p *ProductController) addProductHandler(w http.ResponseWriter, req *http.Request) {
	var request rest.PostProductsJSONRequestBody
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		rest.WriteError(w, p.log, http.StatusBadRequest, "invalid request")
		return
	}

	var productType models.ProductType
	switch string(request.Type) {
	case string(rest.ProductTypeОбувь):
		productType = models.ProductTypeShoes
	case string(rest.ProductTypeОдежда):
		productType = models.ProductTypeClothing
	case string(rest.ProductTypeЭлектроника):
		productType = models.ProductTypeElectronics
	default:
		rest.WriteError(w, p.log, http.StatusBadRequest, "invalid product type")
		return
	}

	product, err := p.productService.AddProduct(req.Context(), productType, request.PvzId)
	if err != nil {
		rest.WriteError(w, p.log, http.StatusBadRequest, err.Error())
		return
	}

	resp := rest.Product{
		Id:          &product.ID,
		DateTime:    &product.DateTime,
		Type:        rest.ProductType(product.Type),
		ReceptionId: product.ReceptionID,
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		p.log.Error("failed to encode response", "error", err)
	}
}
