package http

import (
	"net/http"
	"strconv"

	usecase "github.com/evrintobing17/ecommerce-system/warehouse-service/app"
	"github.com/evrintobing17/ecommerce-system/warehouse-service/app/models"
	"github.com/gin-gonic/gin"
)

type WarehouseHandler struct {
	warehouseUsecase usecase.WarehouseUsecase
}

func NewWarehouseHandler(warehouseUsecase usecase.WarehouseUsecase) *WarehouseHandler {
	return &WarehouseHandler{warehouseUsecase: warehouseUsecase}
}

func (h *WarehouseHandler) GetWarehouse(c *gin.Context) {
	warehouseID, _ := strconv.Atoi(c.Param("id"))

	warehouse, err := h.warehouseUsecase.GetWarehouse(warehouseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "warehouse not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"warehouse": warehouse,
	})
}

func (h *WarehouseHandler) GetWarehouses(c *gin.Context) {
	shopID, _ := strconv.Atoi(c.Query("shop_id"))
	activeOnly := c.Query("active_only") == "true"

	warehouses, err := h.warehouseUsecase.GetWarehouses(shopID, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"warehouses": warehouses,
	})
}

func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var request struct {
		Name     string `json:"name" binding:"required"`
		Location string `json:"location" binding:"required"`
		ShopID   int    `json:"shop_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	warehouse, err := h.warehouseUsecase.CreateWarehouse(request.Name, request.Location, request.ShopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"warehouse": warehouse,
	})
}

func (h *WarehouseHandler) UpdateWarehouse(c *gin.Context) {
	warehouseID, _ := strconv.Atoi(c.Param("id"))

	var request struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Active   *bool  `json:"active"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	warehouse, err := h.warehouseUsecase.UpdateWarehouse(warehouseID, request.Name, request.Location, request.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"warehouse": warehouse,
	})
}

func (h *WarehouseHandler) TransferStock(c *gin.Context) {
	var request struct {
		ProductID       int   `json:"product_id" binding:"required"`
		FromWarehouseID int   `json:"from_warehouse_id" binding:"required"`
		ToWarehouseID   int   `json:"to_warehouse_id" binding:"required"`
		Quantity        int32 `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.warehouseUsecase.TransferStock(request.ProductID, request.FromWarehouseID, request.ToWarehouseID, request.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stock transferred successfully",
	})
}

func (h *WarehouseHandler) GetStock(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Query("product_id"))
	warehouseID, _ := strconv.Atoi(c.Query("warehouse_id"))

	if productID == 0 || warehouseID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id and warehouse_id are required"})
		return
	}

	stock, err := h.warehouseUsecase.GetStock(productID, warehouseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "stock not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stock": stock,
	})
}

func (h *WarehouseHandler) UpdateStock(c *gin.Context) {
	var request struct {
		ProductID   int    `json:"product_id" binding:"required"`
		WarehouseID int    `json:"warehouse_id" binding:"required"`
		Operation   string `json:"operation" binding:"required,oneof=add subtract set"`
		Quantity    int32  `json:"quantity" binding:"required,min=1"`
		Reserved    int32  `json:"reserved" binding:"min=0"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var err error
	var stock *models.Stock

	switch request.Operation {
	case "add":
		stock, err = h.warehouseUsecase.AddStock(request.ProductID, request.WarehouseID, request.Quantity, request.Reserved)
	case "subtract":
		stock, err = h.warehouseUsecase.SubtractStock(request.ProductID, request.WarehouseID, request.Quantity, request.Reserved)
	case "set":
		stock, err = h.warehouseUsecase.SetStock(request.ProductID, request.WarehouseID, request.Quantity, request.Reserved)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stock": stock,
	})
}
