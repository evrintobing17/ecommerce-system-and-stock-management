package http

import (
	"net/http"
	"strconv"

	usecase "github.com/evrintobing17/ecommerce-system/shop-service/app"
	"github.com/gin-gonic/gin"
)

type ShopHandler struct {
	shopUsecase usecase.ShopUsecase
}

func NewShopHandler(shopUsecase usecase.ShopUsecase) *ShopHandler {
	return &ShopHandler{shopUsecase: shopUsecase}
}

func (h *ShopHandler) CreateShop(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shop, err := h.shopUsecase.CreateShop(request.Name, request.Description, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"shop": shop,
	})
}

func (h *ShopHandler) GetShop(c *gin.Context) {
	shopID, _ := strconv.Atoi(c.Param("id"))

	shop, err := h.shopUsecase.GetShop(shopID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shop": shop,
	})
}

func (h *ShopHandler) GetMyShops(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	shops, total, err := h.shopUsecase.GetShops(userID.(int), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shops": shops,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *ShopHandler) UpdateShop(c *gin.Context) {
	shopID, _ := strconv.Atoi(c.Param("id"))
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	shop, err := h.shopUsecase.GetShop(shopID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}

	// Check if the user owns this shop
	if shop.OwnerID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Name != "" {
		shop.Name = request.Name
	}
	if request.Description != "" {
		shop.Description = request.Description
	}

	err = h.shopUsecase.UpdateShop(shop)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shop": shop,
	})
}

func (h *ShopHandler) DeleteShop(c *gin.Context) {
	shopID, _ := strconv.Atoi(c.Param("id"))
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	shop, err := h.shopUsecase.GetShop(shopID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}

	// Check if the user owns this shop
	if shop.OwnerID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	err = h.shopUsecase.DeleteShop(shopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shop deleted successfully",
	})
}
