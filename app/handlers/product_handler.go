package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/guatom999/ecommerce-product-api/app/models"
	"github.com/guatom999/ecommerce-product-api/app/services"
	"github.com/guatom999/ecommerce-product-api/app/utils"
)

type ProductHandler struct{ svc services.ProductService }

func NewProductHandler(svc services.ProductService) *ProductHandler { return &ProductHandler{svc: svc} }

func (h *ProductHandler) Create(c echo.Context) error {
	var req models.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "bad request"})
	}
	p, err := h.svc.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, p)
}

func (h *ProductHandler) Get(c echo.Context) error {
	id := c.Param("id")
	p, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
	}
	return c.JSON(http.StatusOK, p)
}

func (h *ProductHandler) List(c echo.Context) error {
	pq := utils.ParsePageQuery(c.QueryParam("page"), c.QueryParam("limit"))
	items, total, err := h.svc.List(c.Request().Context(), pq.Limit, pq.Offset())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal"})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data":       items,
		"pagination": echo.Map{"page": pq.Page, "limit": pq.Limit, "total": total},
	})
}
