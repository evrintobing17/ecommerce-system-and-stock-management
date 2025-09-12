package grpc

import (
	"context"
	"log"

	proto "github.com/evrintobing17/ecommerce-system/shared/proto/order"

	"github.com/evrintobing17/ecommerce-system/order-service/app"
	"github.com/evrintobing17/ecommerce-system/order-service/app/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type orderServer struct {
	proto.UnimplementedOrderServiceServer
	orderUsecase app.OrderUsecase
}

func NewOrderServer(orderUsecase app.OrderUsecase) *orderServer {
	return &orderServer{orderUsecase: orderUsecase}
}

func (s *orderServer) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	// Convert proto items to domain items
	var items []models.OrderItem
	for _, item := range req.Items {
		items = append(items, models.OrderItem{
			ProductID: int(item.ProductId),
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	order, err := s.orderUsecase.CreateOrder(int(req.UserId), items)
	if err != nil {
		log.Printf("CreateOrder error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	// Convert domain order to proto order
	var protoItems []*proto.OrderItem
	for _, item := range order.Items {
		protoItems = append(protoItems, &proto.OrderItem{
			ProductId: int32(item.ProductID),
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return &proto.CreateOrderResponse{
		Order: &proto.Order{
			Id:          int32(order.ID),
			UserId:      int32(order.UserID),
			Items:       protoItems,
			TotalAmount: order.TotalAmount,
			Status:      string(order.Status),
			CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   order.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (s *orderServer) GetOrder(ctx context.Context, req *proto.GetOrderRequest) (*proto.GetOrderResponse, error) {
	order, err := s.orderUsecase.GetOrder(int(req.OrderId))
	if err != nil {
		log.Printf("GetOrder error: %v", err)
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	// Convert domain order to proto order
	var protoItems []*proto.OrderItem
	for _, item := range order.Items {
		protoItems = append(protoItems, &proto.OrderItem{
			ProductId: int32(item.ProductID),
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return &proto.GetOrderResponse{
		Order: &proto.Order{
			Id:          int32(order.ID),
			UserId:      int32(order.UserID),
			Items:       protoItems,
			TotalAmount: order.TotalAmount,
			Status:      string(order.Status),
			CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   order.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (s *orderServer) ProcessPayment(ctx context.Context, req *proto.ProcessPaymentRequest) (*proto.ProcessPaymentResponse, error) {
	order, err := s.orderUsecase.ProcessPayment(int(req.OrderId), req.PaymentMethod, req.PaymentDetails)
	if err != nil {
		log.Printf("ProcessPayment error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to process payment: %v", err)
	}

	// Convert domain order to proto order
	var protoItems []*proto.OrderItem
	for _, item := range order.Items {
		protoItems = append(protoItems, &proto.OrderItem{
			ProductId: int32(item.ProductID),
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return &proto.ProcessPaymentResponse{
		Success: true,
		Message: "Payment processed successfully",
		Order: &proto.Order{
			Id:          int32(order.ID),
			UserId:      int32(order.UserID),
			Items:       protoItems,
			TotalAmount: order.TotalAmount,
			Status:      string(order.Status),
			CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   order.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (s *orderServer) CancelOrder(ctx context.Context, req *proto.CancelOrderRequest) (*proto.CancelOrderResponse, error) {
	err := s.orderUsecase.CancelOrder(int(req.OrderId))
	if err != nil {
		log.Printf("CancelOrder error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to cancel order: %v", err)
	}

	return &proto.CancelOrderResponse{
		Success: true,
		Message: "Order cancelled successfully",
	}, nil
}
