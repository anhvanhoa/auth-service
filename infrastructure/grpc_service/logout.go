package grpcservice

import (
	"context"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) Logout(ctx context.Context, req *proto_auth.LogoutRequest) (*proto_auth.LogoutResponse, error) {
	if _, err := a.logoutUc.VerifyToken(req.GetRefreshToken()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Đăng xuất thất bại")
	}

	if err := a.logoutUc.Logout(req.GetRefreshToken()); err != nil {
		return nil, status.Error(codes.Internal, "Đăng xuất thất bại")
	}

	if err := a.cache.Delete(req.GetAccessToken()); err != nil {
		return nil, status.Errorf(codes.Internal, "Đăng xuất thất bại")
	}

	return &proto_auth.LogoutResponse{
		Message: "Đăng xuất thành công",
	}, nil
}
