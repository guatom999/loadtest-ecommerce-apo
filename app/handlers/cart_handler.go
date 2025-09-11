package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/guatom999/ecommerce-product-api/app/models"
	"github.com/guatom999/ecommerce-product-api/app/services"
)

type CartHandler struct{ svc services.CartService }

func NewCartHandler(svc services.CartService) *CartHandler { return &CartHandler{svc: svc} }

func (h *CartHandler) AddToCart(c echo.Context) error {
	var req models.AddToCartRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "bad request"})
	}
	userID := services.UserIDFromCtx(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	if err := h.svc.AddToCart(c.Request().Context(), userID, req.ProductID, req.Quantity); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *CartHandler) GetCart(c echo.Context) error {
	userID := services.UserIDFromCtx(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	resp, err := h.svc.GetCart(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal"})
	}
	return c.JSON(http.StatusOK, resp)
}
