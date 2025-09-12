package grpc

import (
	"context"
	"log"

	usecase "github.com/evrintobing17/ecommerce-system/product-service/app"
	proto "github.com/evrintobing17/ecommerce-system/shared/proto/product"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type productServer struct {
	proto.UnimplementedProductServiceServer
	productUsecase usecase.ProductUsecase
}

func NewProductServer(productUsecase usecase.ProductUsecase) *productServer {
	return &productServer{productUsecase: productUsecase}
}

func (s *productServer) GetProducts(ctx context.Context, req *proto.GetProductsRequest) (*proto.GetProductsResponse, error) {
	products, total, err := s.productUsecase.GetProducts(int(req.ShopId), int(req.Page), int(req.Limit))
	if err != nil {
		log.Printf("GetProducts error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get products: %v", err)
	}

	var protoProducts []*proto.Product
	for _, product := range products {
		protoProducts = append(protoProducts, &proto.Product{
			Id:          int32(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			ShopId:      int32(product.ShopID),
			CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &proto.GetProductsResponse{
		Products: protoProducts,
		Total:    int32(total),
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

func (s *productServer) GetProduct(ctx context.Context, req *proto.GetProductRequest) (*proto.GetProductResponse, error) {
	product, err := s.productUsecase.GetProduct(int(req.ProductId))
	if err != nil {
		log.Printf("GetProduct error: %v", err)
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}

	return &proto.GetProductResponse{
		Product: &proto.Product{
			Id:          int32(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			ShopId:      int32(product.ShopID),
			CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
