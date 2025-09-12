package grpc

import (
	"context"
	"log"

	proto "github.com/evrintobing17/ecommerce-system/shared/proto/shop"
	usecase "github.com/evrintobing17/ecommerce-system/shop-service/app"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type shopServer struct {
	proto.UnimplementedShopServiceServer
	shopUsecase usecase.ShopUsecase
}

func NewShopServer(shopUsecase usecase.ShopUsecase) *shopServer {
	return &shopServer{shopUsecase: shopUsecase}
}

func (s *shopServer) CreateShop(ctx context.Context, req *proto.CreateShopRequest) (*proto.CreateShopResponse, error) {
	shop, err := s.shopUsecase.CreateShop(req.Name, req.Description, int(req.OwnerId))
	if err != nil {
		log.Printf("CreateShop error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create shop: %v", err)
	}

	return &proto.CreateShopResponse{
		Shop: &proto.Shop{
			Id:          int32(shop.ID),
			Name:        shop.Name,
			Description: shop.Description,
			OwnerId:     int32(shop.OwnerID),
			CreatedAt:   shop.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   shop.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (s *shopServer) GetShop(ctx context.Context, req *proto.GetShopRequest) (*proto.GetShopResponse, error) {
	shop, err := s.shopUsecase.GetShop(int(req.ShopId))
	if err != nil {
		log.Printf("GetShop error: %v", err)
		return nil, status.Errorf(codes.NotFound, "shop not found: %v", err)
	}

	return &proto.GetShopResponse{
		Shop: &proto.Shop{
			Id:          int32(shop.ID),
			Name:        shop.Name,
			Description: shop.Description,
			OwnerId:     int32(shop.OwnerID),
			CreatedAt:   shop.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   shop.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (s *shopServer) GetShops(ctx context.Context, req *proto.GetShopsRequest) (*proto.GetShopsResponse, error) {
	shops, total, err := s.shopUsecase.GetShops(int(req.OwnerId), int(req.Page), int(req.Limit))
	if err != nil {
		log.Printf("GetShops error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get shops: %v", err)
	}

	var protoShops []*proto.Shop
	for _, shop := range shops {
		protoShops = append(protoShops, &proto.Shop{
			Id:          int32(shop.ID),
			Name:        shop.Name,
			Description: shop.Description,
			OwnerId:     int32(shop.OwnerID),
			CreatedAt:   shop.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   shop.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &proto.GetShopsResponse{
		Shops: protoShops,
		Total: int32(total),
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}
