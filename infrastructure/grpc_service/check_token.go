package grpcservice

import (
	"context"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) CheckToken(ctx context.Context, req *proto_auth.CheckTokenRequest) (*proto_auth.CheckTokenResponse, error) {
	ok, err := a.checkTokenUc.CheckToken(req.GetToken())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &proto_auth.CheckTokenResponse{
		Data:    ok,
		Message: "Mã truyền lên hợp lệ",
	}, nil
}
