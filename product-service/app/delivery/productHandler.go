package http

import (
	"net/http"
	"strconv"

	usecase "github.com/evrintobing17/ecommerce-system/product-service/app"
	"github.com/evrintobing17/ecommerce-system/shared/jsonhttpresponse"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecase
}

func NewProductHandler(productUsecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: productUsecase}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	shopID, _ := strconv.Atoi(c.Query("shop_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, total, err := h.productUsecase.GetProducts(shopID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jsonhttpresponse.OK(c, gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Param("id"))

	product, err := h.productUsecase.GetProduct(productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var request struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Price       float64 `json:"price" binding:"required,min=0"`
		Stock       int32   `json:"stock" binding:"required,min=0"`
		ShopID      int     `json:"shop_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productUsecase.CreateProduct(request.Name, request.Description, request.Price, request.Stock, request.ShopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"product": product,
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Param("id"))
	var request struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price" min:"0"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productUsecase.GetProduct(productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if request.Name != "" {
		product.Name = request.Name
	}
	if request.Description != "" {
		product.Description = request.Description
	}
	if request.Price > 0 {
		product.Price = request.Price
	}

	err = h.productUsecase.UpdateProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func (h *ProductHandler) UpdateStock(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Param("id"))

	var request struct {
		Operation string `json:"operation" binding:"required,oneof=add subtract set"`
		Quantity  int32  `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var err error
	switch request.Operation {
	case "add":
		err = h.productUsecase.AddStock(productID, request.Quantity)
	case "subtract":
		err = h.productUsecase.SubtractStock(productID, request.Quantity)
	case "set":
		err = h.productUsecase.SetStock(productID, request.Quantity)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productUsecase.GetProduct(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Param("id"))

	err := h.productUsecase.DeleteProduct(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
