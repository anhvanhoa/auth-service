package grpcservice

import (
	"context"
	"fmt"
	"time"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	proto_user_role "github.com/anhvanhoa/sf-proto/gen/user_role/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) RefreshToken(ctx context.Context, req *proto_auth.RefreshTokenRequest) (*proto_auth.RefreshTokenResponse, error) {
	if !a.refreshUc.CheckSessionByToken(req.GetRefreshToken()) {
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

	permissions, err := a.permissionClient.UserRoleService().GetUserPermissions(ctx, &proto_user_role.GetUserPermissionsRequest{
		UserId: claims.Data.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể lấy quyền")
	}
	uCtx := a.convertPermissions(permissions)
	bytes, err := uCtx.ToBytes()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể chuyển đổi quyền")
	}
	if err := a.cache.Set(accessToken, bytes, time.Until(accessExp)); err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể lưu quyền")
	}

	return &proto_auth.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "Làm mới token thành công",
	}, nil
}
