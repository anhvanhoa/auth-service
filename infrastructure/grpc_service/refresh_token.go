package grpcservice

import (
	"context"
	"fmt"
	"time"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) RefreshToken(ctx context.Context, req *proto_auth.RefreshTokenRequest) (*proto_auth.RefreshTokenResponse, error) {
	if _, err := a.refreshUc.GetSessionByToken(req.GetRefreshToken()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Phiên làm việc không hợp lệ")
	}

	claims, err := a.refreshUc.VerifyToken(req.GetRefreshToken())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Token không hợp lệ")
	}

	if err := a.refreshUc.ClearSessionExpired(); err != nil {
		a.log.Info(fmt.Sprintf("Clear expired sessions: %v", err))
	}

	accessExp := time.Now().Add(15 * time.Minute)
	accessToken, err := a.refreshUc.GengerateAccessToken(claims.Data.Id, claims.Data.FullName, claims.Data.Email, accessExp)
	if err != nil {
		return nil, status.Error(codes.Internal, "Không thể tạo access token")
	}

	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	refreshToken, err := a.refreshUc.GengerateRefreshToken(claims.Data.Id, claims.Data.FullName, claims.Data.Email, refreshExp, req.GetOs())
	if err != nil {
		return nil, status.Error(codes.Internal, "Không thể tạo refresh token")
	}

	return &proto_auth.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "Làm mới token thành công",
	}, nil
}
