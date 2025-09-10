package jsonhttpresponse

//type helpers. Used as an example value in swagger documentation
type FailedResponse struct {
	Status  string      `json:"status" example:"failed"`
	Message interface{} `json:"message"`
}

type FailedUnauthorizedInvalidCredentialResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"invalid credential"`
}

type FailedUnauthorizedResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"invalid user token"`
}

type FailedBadRequestResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"bad request"`
}

type FailedInternalServerErrorResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"internal server error"`
}

//NewFailedResponse will return a json envelope (wrapper) to the
//HTTP Error response code
func NewFailedResponse(message interface{}) FailedResponse {
	return FailedResponse{Status: "failed", Message: message}
}
