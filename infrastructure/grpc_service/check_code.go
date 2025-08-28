package grpcservice

import (
	"context"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) CheckCode(ctx context.Context, req *proto_auth.CheckCodeRequest) (*proto_auth.CheckCodeResponse, error) {
	valid, err := a.checkCodeUc.CheckCode(req.GetCode(), req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if !valid {
		return &proto_auth.CheckCodeResponse{
			Valid:   false,
			Message: "Mã xác thực không hợp lệ hoặc đã hết hạn",
		}, nil
	}

	return &proto_auth.CheckCodeResponse{
		Valid:   true,
		Message: "Mã xác thực hợp lệ",
	}, nil
}
