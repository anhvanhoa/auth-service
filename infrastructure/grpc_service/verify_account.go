package grpcservice

import (
	"context"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) VerifyAccount(ctx context.Context, req *proto_auth.VerifyAccountRequest) (*proto_auth.VerifyAccountResponse, error) {
	// Verify register token
	claims, err := a.verifyAccountUc.VerifyRegister(req.GetToken())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Token không hợp lệ hoặc đã hết hạn")
	}

	// Get user by ID
	_, err = a.verifyAccountUc.GetUserById(claims.Data.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Không tìm thấy người dùng")
	}

	// Verify account
	if err := a.verifyAccountUc.VerifyAccount(claims.Data.Id); err != nil {
		return nil, status.Error(codes.Internal, "Không thể xác thực tài khoản")
	}

	return &proto_auth.VerifyAccountResponse{
		Message: "Xác thực tài khoản thành công",
	}, nil
}
