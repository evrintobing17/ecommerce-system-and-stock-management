package http

import (
	"net/http"

	"github.com/evrintobing17/ecommerce-system/shared/jsonhttpresponse"
	"github.com/evrintobing17/ecommerce-system/user-service/app"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase app.UserUsecase
}

func NewUserHandler(userUsecase app.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) Register(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonhttpresponse.ErrBind(c, err)
		return
	}

	user, token, err := h.userUsecase.Register(request.Email, request.Phone, request.Password, request.Name)
	if err != nil {
		jsonhttpresponse.InternalServerError(c, err)
		return
	}
	jsonhttpresponse.StatusCreated(c, gin.H{"user": gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"phone":      user.Phone,
		"name":       user.Name,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	},
		"token": token})
}

func (h *UserHandler) Login(c *gin.Context) {
	var request struct {
		EmailOrPhone string `json:"email_or_phone" binding:"required"`
		Password     string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonhttpresponse.ErrBind(c, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.userUsecase.Login(request.EmailOrPhone, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"phone":      user.Phone,
			"name":       user.Name,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
		"token": token,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		jsonhttpresponse.Unauthorized(c, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userUsecase.GetUser(userID.(int))
	if err != nil {
		jsonhttpresponse.BadRequest(c, err)
		return
	}

	jsonhttpresponse.OK(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request struct {
		Name  string `json:"name"`
		Email string `json:"email" binding:"email"`
		Phone string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userUsecase.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if request.Name != "" {
		user.Name = request.Name
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.Phone != "" {
		user.Phone = request.Phone
	}

	err = h.userUsecase.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"phone":      user.Phone,
			"name":       user.Name,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}
