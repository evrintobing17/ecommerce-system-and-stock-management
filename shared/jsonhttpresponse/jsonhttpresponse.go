package jsonhttpresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//OK - Function to return Status OK Response (200)
func OK(c *gin.Context, payloads interface{}) {
	c.JSON(http.StatusOK, payloads)
	return
}

//BadRequest - Function to return Status Bad Request Response (400)
//use it if user request is wrong
func BadRequest(c *gin.Context, payloads interface{}) {
	c.JSON(http.StatusBadRequest, payloads)
	return
}

//InternalServerError - Function to return Internal Server Error Response (400)
//use it for any unhandled error that is not user's fault
func InternalServerError(c *gin.Context, payloads interface{}) {
	c.JSON(http.StatusInternalServerError, payloads)
	return
}

//Unauthorized - Function to return Unauthorized Response (401)
//Use it only in authentication process
func Unauthorized(c *gin.Context, payloads interface{}) {
	c.JSON(http.StatusUnauthorized, payloads)
	return
}

//NotFound - Function to return Not Found Response (404)
//Use it in case of any get operation that retrieve
//for resource and not exist
func NotFound(c *gin.Context, payloads interface{}) {
	c.JSON(http.StatusNotFound, payloads)
	return
}

//Conflict - Function to return Conflict Response (409)
//Use it in case if a process create a new resource,
//but somehow, another resource already exist
//(collision in unique identifier)
func Conflict(c *gin.Context, payloads interface{}) {
	c.JSON(http.StatusConflict, payloads)
	return
}

//Forbidden - Function to return Forbidden Response (403)
//Use it for any user attempting to access resource
//with lack of authorization
func Forbidden(c *gin.Context, payloads interface{}) {
	c.JSON(http.StatusForbidden, payloads)
	return
}

func StatusCreated(c *gin.Context, payload interface{}) {
	c.JSON(http.StatusCreated, payload)
}
