package delivery

import (
	"errors"
	"fmt"
	"net/http"
	user "user-service/app"
	"user-service/app/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userUseCase user.UserUsecase
}

func NewAuthHandler(r *gin.Engine, userUseCase user.UserUsecase) {
	handler := &UserHandler{
		userUseCase: userUseCase,
	}

	authorized := r.Group("/v1/user")
	{
		authorized.POST("/login", handler.Login)
	}
}

func (h *UserHandler) Login(c *gin.Context) {
	var loginReq models.LoginRequest
	var loginData string
	var errs []models.Error
	errBind := c.ShouldBind(&loginReq)
	if errBind != nil {
		if verrs, ok := errBind.(validator.ValidationErrors); ok {
			for _, e := range verrs {
				errs = append(errs, models.Error{
					Key:   e.Field(), // or e.Namespace() for full path
					Error: fmt.Sprintf("failed on '%s' rule", e.Tag()),
				})
			}
		}
		c.IndentedJSON(http.StatusBadRequest, map[string]interface{}{"error": errs})
		return
	}

	loginData = loginReq.Phone
	if loginReq.Email != "" {
		loginData = loginReq.Email
	}

	token, user, err := h.userUseCase.Login(c, loginData, loginReq.Password)
	if err != nil {
		responseErr := map[string]interface{}{
			"err": err.Error(),
		}
		if errors.Is(err, models.ErrUserNotFound) {
			c.IndentedJSON(http.StatusBadRequest, responseErr)
			return
		}
		c.IndentedJSON(http.StatusUnauthorized, responseErr)
		return
	}

	response := map[string]interface{}{
		"token": token,
		"user":  user,
	}

	resp := models.Response{
		Data: response,
	}

	c.IndentedJSON(http.StatusOK, resp)
}
