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
			Stock:       product.Stock,
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
			Stock:       product.Stock,
			ShopId:      int32(product.ShopID),
			CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (s *productServer) UpdateStock(ctx context.Context, req *proto.UpdateStockRequest) (*proto.UpdateStockResponse, error) {
	var err error
	var newStock int32

	switch req.Operation {
	case "add":
		err = s.productUsecase.AddStock(int(req.ProductId), req.Quantity)
	case "subtract":
		err = s.productUsecase.SubtractStock(int(req.ProductId), req.Quantity)
	case "set":
		err = s.productUsecase.SetStock(int(req.ProductId), req.Quantity)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid operation: %s", req.Operation)
	}

	if err != nil {
		log.Printf("UpdateStock error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update stock: %v", err)
	}

	product, err := s.productUsecase.GetProduct(int(req.ProductId))
	if err != nil {
		log.Printf("GetProduct after update error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get updated product: %v", err)
	}

	newStock = product.Stock

	return &proto.UpdateStockResponse{
		Success:  true,
		NewStock: newStock,
	}, nil
}
