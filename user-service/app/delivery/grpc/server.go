package grpc

import (
	"context"
	"log"

	proto "github.com/evrintobing17/ecommerce-system/shared/proto/user"
	usecase "github.com/evrintobing17/ecommerce-system/user-service/app"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServer struct {
	proto.UnimplementedUserServiceServer
	userUsecase usecase.UserUsecase
}

func NewUserServer(userUsecase usecase.UserUsecase) *userServer {
	return &userServer{userUsecase: userUsecase}
}

func (s *userServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	user, token, err := s.userUsecase.Register(req.Email, req.Phone, req.Password, req.Name)
	if err != nil {
		log.Printf("Register error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &proto.RegisterResponse{
		User: &proto.User{
			Id:        int32(user.ID),
			Email:     user.Email,
			Phone:     user.Phone,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		Token: token,
	}, nil
}

func (s *userServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	user, token, err := s.userUsecase.Login(req.EmailOrPhone, req.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials: %v", err)
	}

	return &proto.LoginResponse{
		User: &proto.User{
			Id:        int32(user.ID),
			Email:     user.Email,
			Phone:     user.Phone,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		Token: token,
	}, nil
}

func (s *userServer) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	valid, user, err := s.userUsecase.ValidateToken(req.Token)
	if err != nil {
		log.Printf("ValidateToken error: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	if !valid {
		return &proto.ValidateTokenResponse{Valid: false}, nil
	}

	return &proto.ValidateTokenResponse{
		Valid: true,
		User: &proto.User{
			Id:        int32(user.ID),
			Email:     user.Email,
			Phone:     user.Phone,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (s *userServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user, err := s.userUsecase.GetUser(int(req.UserId))
	if err != nil {
		log.Printf("GetUser error: %v", err)
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &proto.GetUserResponse{
		User: &proto.User{
			Id:        int32(user.ID),
			Email:     user.Email,
			Phone:     user.Phone,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
