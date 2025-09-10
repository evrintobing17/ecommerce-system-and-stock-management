package delivery

import (
	product "product-service/app"
	"shared"

	jResp "shared/jsonhttpresponse"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	log            shared.Log
	productUseCase product.ProductUsecase
}

func NewAuthHandler(r *gin.Engine, log shared.Log, productUseCase product.ProductUsecase) {
	handler := &ProductHandler{
		productUseCase: productUseCase,
		log:            log,
	}

	authorized := r.Group("/v1/product")
	{
		authorized.GET("", handler.ListProduct)
	}
}

func (h *ProductHandler) ListProduct(c *gin.Context) {

	product, err := h.productUseCase.ProductList(c)
	if err != nil {
		c.Set("stackTrace", h.log.SetMessageLog(err))
		jResp.BadRequest(c, err.Error())
		return
	}

	h.log.InfoLog("success")
	jResp.OK(c, product)
	return
}
