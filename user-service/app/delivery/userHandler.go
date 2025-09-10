package delivery

import (
	"shared"
	user "user-service/app"
	"user-service/app/models"

	jResp "shared/jsonhttpresponse"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	log         shared.Log
	userUseCase user.UserUsecase
}

func NewAuthHandler(r *gin.Engine, log shared.Log, userUseCase user.UserUsecase) {
	handler := &UserHandler{
		userUseCase: userUseCase,
		log:         log,
	}

	authorized := r.Group("/v1/user")
	{
		authorized.POST("/login", handler.Login)
	}
}

func (h *UserHandler) Login(c *gin.Context) {
	var loginReq models.LoginRequest
	var loginData string
	errBind := c.ShouldBind(&loginReq)
	if errBind != nil {
		c.Set("stackTrace", h.log.SetMessageLog(errBind))
		jResp.ErrBind(c, errBind)
		return
	}

	loginData = loginReq.Phone
	if loginReq.Email != "" {
		loginData = loginReq.Email
	}

	token, user, err := h.userUseCase.Login(c, loginData, loginReq.Password)
	if err != nil {
		c.Set("stackTrace", h.log.SetMessageLog(err))
		jResp.BadRequest(c, err.Error())
		return
	}

	response := map[string]interface{}{
		"token": token,
		"user":  user,
	}
	h.log.InfoLog("success")
	jResp.OK(c, response)
	return
}
