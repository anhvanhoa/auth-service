package grpcservice

import (
	"auth-service/domain/entity"
	"context"
	"strings"

	"github.com/anhvanhoa/service-core/constants"
	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *authService) Profile(ctx context.Context, req *emptypb.Empty) (*proto_auth.ProfileResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Context không chứa metadata")
	}
	token := a.getCookieFromMetadata(md, constants.KeyCookieAccessToken)
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "Cần đăng nhập để xem thông tin tài khoản")
	}
	user, err := a.profileUc.Execute(ctx, token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto_auth.ProfileResponse{
		User: a.convertProfile(user),
	}, nil
}

func (a *authService) getCookieFromMetadata(md metadata.MD, key string) string {
	cookies := a.getFirstValue(md, constants.Cookie)
	if cookies == "" {
		return ""
	}
	for _, c := range strings.Split(cookies, ";") {
		if strings.Contains(c, key+"=") {
			return strings.TrimSpace(strings.Split(c, "=")[1])
		}
	}
	return ""
}

func (a *authService) getFirstValue(md metadata.MD, key string) string {
	vals := md.Get(key)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func (a *authService) convertProfile(user *entity.UserInfor) *proto_auth.UserInfo {
	userInfo := &proto_auth.UserInfo{
		Id:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
		Address:  user.Address,
	}
	if user.Birthday != nil {
		userInfo.Birthday = timestamppb.New(*user.Birthday)
	}
	return userInfo
}
