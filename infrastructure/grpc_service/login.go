package grpcservice

import (
	"context"
	"regexp"
	"time"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *authService) Login(ctx context.Context, req *proto_auth.LoginRequest) (*proto_auth.LoginResponse, error) {
	identifier := req.GetEmailOrPhone()
	if !isValidEmail(identifier) && !isValidPhone(identifier) {
		return nil, status.Errorf(codes.InvalidArgument, "email hoặc số điện thoại không đúng định dạng")
	}

	user, err := a.loginUc.GetUserByEmailOrPhone(req.GetEmailOrPhone())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Không tìm thấy người dùng")
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

func isValidPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^(0|\+84)(3[2-9]|5[689]|7[06-9]|8[1-689]|9[0-46-9])[0-9]{7}$`)
	return phoneRegex.MatchString(phone)
}
