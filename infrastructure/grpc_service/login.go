package grpcservice

import (
	"context"
	"regexp"
	"time"

	"github.com/anhvanhoa/service-core/domain/user_context"
	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	proto_user_role "github.com/anhvanhoa/sf-proto/gen/user_role/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *authService) Login(ctx context.Context, req *proto_auth.LoginRequest) (*proto_auth.LoginResponse, error) {
	client := a.permissionClient.UserRoleService()
	identifier := req.GetEmailOrPhone()
	if !isValidEmail(identifier) && !isValidPhone(identifier) {
		return nil, status.Errorf(codes.InvalidArgument, "Email hoặc số điện thoại không đúng định dạng")
	}

	user, err := a.loginUc.GetUserByEmailOrPhone(req.GetEmailOrPhone())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if !a.loginUc.CheckHashPassword(req.GetPassword(), user.Password) {
		return nil, status.Errorf(codes.InvalidArgument, "Mật khẩu không chính xác")
	}

	exp := time.Now().Add(15 * time.Minute)
	accessToken, err := a.loginUc.GengerateAccessToken(user.ID, user.FullName, user.Email, exp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể tạo access token")
	}

	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	refreshToken, err := a.loginUc.GengerateRefreshToken(user.ID, user.FullName, user.Email, refreshExp, req.GetOs())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể tạo refresh token")
	}

	permissions, err := client.GetUserPermissions(ctx, &proto_user_role.GetUserPermissionsRequest{
		UserId: user.ID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể lấy quyền")
	}

	uCtx := a.convertPermissions(permissions)
	bytes, err := uCtx.ToBytes()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể chuyển đổi quyền")
	}
	if err := a.cache.Set(user.ID, bytes, time.Until(exp)); err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể lưu quyền")
	}

	userInfo := &proto_auth.UserInfo{
		Id:       user.ID,
		Email:    user.Email,
		Phone:    user.Phone,
		FullName: user.FullName,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
		Address:  user.Address,
	}

	if user.Birthday != nil {
		userInfo.Birthday = timestamppb.New(*user.Birthday)
	}

	return &proto_auth.LoginResponse{
		User:         userInfo,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "Đăng nhập thành công",
	}, nil
}

func (a *authService) convertPermissions(data *proto_user_role.GetUserPermissionsResponse) *user_context.UserContext {
	uCtx := user_context.NewUserContext()
	uCtx.UserID = data.UserId
	scopes := make([]user_context.Scope, len(data.Scopes))
	for i, scope := range data.Scopes {
		scopes[i] = user_context.Scope{
			Resource:     scope.Resource,
			ResourceData: scope.ResourceData,
			Action:       scope.Action,
		}
	}
	uCtx.Scopes = scopes
	permissions := make([]user_context.Permission, len(data.Permissions))
	for i, permission := range data.Permissions {
		permissions[i] = user_context.Permission{
			Resource: permission.Resource,
			Action:   permission.Action,
		}
	}
	uCtx.Permissions = permissions
	return uCtx
}

func isValidPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^(0|\+84)(3[2-9]|5[689]|7[06-9]|8[1-689]|9[0-46-9])[0-9]{7}$`)
	return phoneRegex.MatchString(phone)
}
