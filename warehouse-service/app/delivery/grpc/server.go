package grpc

import (
	"context"
	"log"

	proto "github.com/evrintobing17/ecommerce-system/shared/proto/warehouse"
	usecase "github.com/evrintobing17/ecommerce-system/warehouse-service/app"
	"github.com/evrintobing17/ecommerce-system/warehouse-service/app/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type warehouseServer struct {
	proto.UnimplementedWarehouseServiceServer
	warehouseUsecase usecase.WarehouseUsecase
}

func NewWarehouseServer(warehouseUsecase usecase.WarehouseUsecase) *warehouseServer {
	return &warehouseServer{warehouseUsecase: warehouseUsecase}
}

func (s *warehouseServer) GetWarehouse(ctx context.Context, req *proto.GetWarehouseRequest) (*proto.GetWarehouseResponse, error) {
	warehouse, err := s.warehouseUsecase.GetWarehouse(int(req.WarehouseId))
	if err != nil {
		log.Printf("GetWarehouse error: %v", err)
		return nil, status.Errorf(codes.NotFound, "warehouse not found: %v", err)
	}

	return &proto.GetWarehouseResponse{
		Warehouse: &proto.Warehouse{
			Id:        int32(warehouse.ID),
			Name:      warehouse.Name,
			Location:  warehouse.Location,
			ShopId:    int32(warehouse.ShopID),
			Active:    warehouse.Active,
			CreatedAt: warehouse.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: warehouse.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (s *warehouseServer) GetWarehouses(ctx context.Context, req *proto.GetWarehousesRequest) (*proto.GetWarehousesResponse, error) {
	warehouses, err := s.warehouseUsecase.GetWarehouses(int(req.ShopId), req.ActiveOnly)
	if err != nil {
		log.Printf("GetWarehouses error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get warehouses: %v", err)
	}

	var protoWarehouses []*proto.Warehouse
	for _, warehouse := range warehouses {
		protoWarehouses = append(protoWarehouses, &proto.Warehouse{
			Id:        int32(warehouse.ID),
			Name:      warehouse.Name,
			Location:  warehouse.Location,
			ShopId:    int32(warehouse.ShopID),
			Active:    warehouse.Active,
			CreatedAt: warehouse.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: warehouse.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &proto.GetWarehousesResponse{
		Warehouses: protoWarehouses,
	}, nil
}

func (s *warehouseServer) CreateWarehouse(ctx context.Context, req *proto.CreateWarehouseRequest) (*proto.CreateWarehouseResponse, error) {
	warehouse, err := s.warehouseUsecase.CreateWarehouse(req.Name, req.Location, int(req.ShopId))
	if err != nil {
		log.Printf("CreateWarehouse error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create warehouse: %v", err)
	}

	return &proto.CreateWarehouseResponse{
		Warehouse: &proto.Warehouse{
			Id:        int32(warehouse.ID),
			Name:      warehouse.Name,
			Location:  warehouse.Location,
			ShopId:    int32(warehouse.ShopID),
			Active:    warehouse.Active,
			CreatedAt: warehouse.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: warehouse.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (s *warehouseServer) UpdateWarehouse(ctx context.Context, req *proto.UpdateWarehouseRequest) (*proto.UpdateWarehouseResponse, error) {
	warehouse, err := s.warehouseUsecase.UpdateWarehouse(int(req.WarehouseId), req.Name, req.Location, &req.Active)
	if err != nil {
		log.Printf("UpdateWarehouse error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update warehouse: %v", err)
	}

	return &proto.UpdateWarehouseResponse{
		Warehouse: &proto.Warehouse{
			Id:        int32(warehouse.ID),
			Name:      warehouse.Name,
			Location:  warehouse.Location,
			ShopId:    int32(warehouse.ShopID),
			Active:    warehouse.Active,
			CreatedAt: warehouse.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: warehouse.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (s *warehouseServer) TransferStock(ctx context.Context, req *proto.TransferStockRequest) (*proto.TransferStockResponse, error) {
	err := s.warehouseUsecase.TransferStock(int(req.ProductId), int(req.FromWarehouseId), int(req.ToWarehouseId), req.Quantity)
	if err != nil {
		log.Printf("TransferStock error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to transfer stock: %v", err)
	}

	return &proto.TransferStockResponse{
		Success: true,
		Message: "Stock transferred successfully",
	}, nil
}

func (s *warehouseServer) GetStock(ctx context.Context, req *proto.GetStockRequest) (*proto.GetStockResponse, error) {
	stock, err := s.warehouseUsecase.GetStock(int(req.ProductId), int(req.WarehouseId))
	if err != nil {
		log.Printf("GetStock error: %v", err)
		return nil, status.Errorf(codes.NotFound, "stock not found: %v", err)
	}

	return &proto.GetStockResponse{
		Stock: &proto.Stock{
			ProductId:   int32(stock.ProductID),
			WarehouseId: int32(stock.WarehouseID),
			Quantity:    stock.Quantity,
			Reserved:    stock.Reserved,
			CreatedAt:   stock.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   stock.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (s *warehouseServer) UpdateStock(ctx context.Context, req *proto.UpdateStockRequest) (*proto.UpdateStockResponse, error) {
	var err error
	var stock *models.Stock

	switch req.Operation {
	case "add":
		stock, err = s.warehouseUsecase.AddStock(int(req.ProductId), int(req.WarehouseId), req.Quantity, req.Reserved)
	case "subtract":
		stock, err = s.warehouseUsecase.SubtractStock(int(req.ProductId), int(req.WarehouseId), req.Quantity, req.Reserved)
	case "set":
		stock, err = s.warehouseUsecase.SetStock(int(req.ProductId), int(req.WarehouseId), req.Quantity, req.Reserved)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid operation: %s", req.Operation)
	}

	if err != nil {
		log.Printf("UpdateStock error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update stock: %v", err)
	}

	return &proto.UpdateStockResponse{
		Success: true,
		Stock: &proto.Stock{
			ProductId:   int32(stock.ProductID),
			WarehouseId: int32(stock.WarehouseID),
			Quantity:    stock.Quantity,
			Reserved:    stock.Reserved,
			CreatedAt:   stock.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   stock.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
