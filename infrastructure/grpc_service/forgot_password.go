package grpcservice

import (
	"auth-service/domain/usecase"
	"context"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *authService) ForgotPassword(ctx context.Context, req *proto_auth.ForgotPasswordRequest) (*proto_auth.ForgotPasswordResponse, error) {
	var method usecase.ForgotPasswordType
	switch req.GetMethod() {
	case proto_auth.ForgotPasswordType_FORGOT_PASSWORD_TYPE_UNSPECIFIED:
		method = usecase.ForgotByCode
	case proto_auth.ForgotPasswordType_FORGOT_PASSWORD_TYPE_TOKEN:
		method = usecase.ForgotByToken
	default:
		return nil, status.Errorf(codes.InvalidArgument, "Phương thức xác thực không hợp lệ")
	}

	result, err := a.forgotPasswordUc.ForgotPassword(req.GetEmail(), req.GetOs(), method)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	userInfo := &proto_auth.UserInfo{
		Id:       result.User.ID,
		Email:    result.User.Email,
		Phone:    result.User.Phone,
		FullName: result.User.FullName,
		Avatar:   result.User.Avatar,
		Bio:      result.User.Bio,
		Address:  result.User.Address,
	}

	if result.User.Birthday != nil {
		userInfo.Birthday = timestamppb.New(*result.User.Birthday)
	}

	return &proto_auth.ForgotPasswordResponse{
		User:    userInfo,
		Token:   result.Token,
		Code:    result.Code,
		Message: "Yêu cầu đặt lại mật khẩu đã được gửi. Vui lòng kiểm tra email.",
	}, nil
}
