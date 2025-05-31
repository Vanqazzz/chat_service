package auth

import (
	"chat_service/pkg"
	"context"
	"errors"

	protos "github.com/Vanqazzz/protos/gen/go/chat_service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Auth
type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		app_id int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
}

type serverAPI struct {
	protos.UnimplementedAuthServer
	auth Auth
}

const (
	emptyValue = 0
)

// Register server
func Register(gRPC *grpc.Server, auth Auth) {
	protos.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

// Login
func (s *serverAPI) Login(
	ctx context.Context,
	req *protos.LoginRequest) (*protos.LoginResponse, error) {

	if err := ValidateLogin(req); err != nil {
		return nil, err
	}

	// Creating token
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, pkg.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "fail login")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &protos.LoginResponse{
		Token: token,
	}, nil

}

// Register
func (s *serverAPI) Register(ctx context.Context, req *protos.RegisterRequest) (*protos.RegisterResponse, error) {

	if err := ValidateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, pkg.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &protos.RegisterResponse{
		UserId: userID,
	}, nil

}

// Validate Login
func ValidateLogin(req *protos.LoginRequest) error {

	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

// Validate Register
func ValidateRegister(req *protos.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	return nil
}
