package jsonhttpresponse

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

type Error struct {
	Key   string `json:"key"`
	Error string `json:"error"`
}

// OK - Function to return Status OK Response (200)
func OK(c *gin.Context, payloads interface{}) {
	status := http.StatusOK
	Resp := Response{
		Code:    status,
		Message: http.StatusText(status),
		Data:    payloads,
		Errors:  nil,
	}
	c.IndentedJSON(status, Resp)
	return
}

// BadRequest - Function to return Status Bad Request Response (400)
// use it if user request is wrong
func BadRequest(c *gin.Context, payloads interface{}) {
	status := http.StatusBadRequest
	Resp := Response{
		Code:    status,
		Message: http.StatusText(status),
		Data:    nil,
		Errors:  payloads,
	}
	c.IndentedJSON(status, Resp)
	return
}

func ErrBind(c *gin.Context, err error) {
	var errs []Error
	status := http.StatusBadRequest

	if verrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range verrs {
			raw := e.Error()
			parts := strings.SplitN(raw, " Error:", 2)
			errs = append(errs, Error{
				Key:   e.Field(), // or e.Namespace() for full path
				Error: parts[1],
			})
		}
	}
	Resp := Response{
		Code:    status,
		Message: http.StatusText(status),
		Data:    nil,
		Errors:  errs,
	}
	c.IndentedJSON(status, Resp)
	return
}

// InternalServerError - Function to return Internal Server Error Response (400)
// use it for any unhandled error that is not user's fault
func InternalServerError(c *gin.Context, payloads interface{}) {
	status := http.StatusInternalServerError
	Resp := Response{
		Code:    status,
		Message: http.StatusText(status),
		Data:    nil,
		Errors:  payloads,
	}
	c.IndentedJSON(status, Resp)
	return
}

// Unauthorized - Function to return Unauthorized Response (401)
// Use it only in authentication process
func Unauthorized(c *gin.Context, payloads interface{}) {
	status := http.StatusUnauthorized
	Resp := Response{
		Code:    status,
		Message: http.StatusText(status),
		Data:    nil,
		Errors:  payloads,
	}
	c.IndentedJSON(status, Resp)
	return
}

// NotFound - Function to return Not Found Response (404)
// Use it in case of any get operation that retrieve
// for resource and not exist
func NotFound(c *gin.Context, payloads interface{}) {
	c.JSON(http.StatusNotFound, payloads)
	return
}

// Conflict - Function to return Conflict Response (409)
// Use it in case if a process create a new resource,
// but somehow, another resource already exist
// (collision in unique identifier)
func Conflict(c *gin.Context, payloads interface{}) {
	c.JSON(http.StatusConflict, payloads)
	return
}

// Forbidden - Function to return Forbidden Response (403)
// Use it for any user attempting to access resource
// with lack of authorization
func Forbidden(c *gin.Context, payloads interface{}) {
	c.JSON(http.StatusForbidden, payloads)
	return
}

func StatusCreated(c *gin.Context, payload interface{}) {
	status := http.StatusCreated
	Resp := Response{
		Code:    status,
		Message: http.StatusText(status),
		Data:    payload,
		Errors:  nil,
	}
	c.IndentedJSON(status, Resp)
	return
}
