package grpcservice

import (
	"context"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) ResetPasswordByCode(ctx context.Context, req *proto_auth.ResetPasswordByCodeRequest) (*proto_auth.ResetPasswordByCodeResponse, error) {
	// Business logic validation: check if passwords match
	if req.GetNewPassword() != req.GetConfirmPassword() {
		return nil, status.Errorf(codes.InvalidArgument, "Mật khẩu mới và xác nhận mật khẩu không khớp")
	}

	// Verify session
	userID, err := a.resetCodeUc.VerifySession(req.GetCode(), req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Reset password
	if err := a.resetCodeUc.ResetPass(userID, req.GetNewPassword(), req.GetConfirmPassword()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto_auth.ResetPasswordByCodeResponse{
		Message: "Đặt lại mật khẩu thành công",
	}, nil
}
