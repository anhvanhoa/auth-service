package grpcservice

import (
	"context"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) ResetPasswordByToken(ctx context.Context, req *proto_auth.ResetPasswordByTokenRequest) (*proto_auth.ResetPasswordByTokenResponse, error) {
	// Business logic validation: check if passwords match
	if req.GetNewPassword() != req.GetConfirmPassword() {
		return nil, status.Error(codes.InvalidArgument, "Mật khẩu mới và xác nhận mật khẩu không khớp")
	}

	// Verify session
	userID, err := a.resetTokenUc.VerifySession(req.GetToken())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Reset password
	if err := a.resetTokenUc.ResetPass(userID, req.GetNewPassword(), req.GetConfirmPassword()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto_auth.ResetPasswordByTokenResponse{
		Message: "Đặt lại mật khẩu thành công",
	}, nil
}
