package grpcservice

import (
	"auth-service/constants"
	"auth-service/domain/entity"
	"auth-service/domain/usecase"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/anhvanhoa/service-core/domain/queue"
	"github.com/anhvanhoa/service-core/domain/saga"
	proto_mail_history "github.com/anhvanhoa/sf-proto/gen/mail_history/v1"
	proto_mail_template "github.com/anhvanhoa/sf-proto/gen/mail_tmpl/v1"
	proto_status_history "github.com/anhvanhoa/sf-proto/gen/status_history/v1"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	var taskId string
	var data map[string]any
	var result usecase.ForgotPasswordRes
	var tmpl *proto_mail_template.GetMailTmplResponse

	sagaId := fmt.Sprintf("forgot-password-%s-%s", req.GetEmail(), a.uuid.Gen())
	err := a.forgotPasswordUc.ForgotPasswordWithSaga(sagaId, func(ctx context.Context, sagaTx saga.SagaTransactionI) error {
		var err error
		sagaTx.AddStep(saga.NewSagaStep(
			"ForgotPassword",
			func(ctx context.Context) error {
				result, err = a.forgotPasswordUc.ForgotPassword(req.GetEmail(), req.GetOs(), method)
				if err != nil {
					return err
				}
				data = map[string]any{"user": result.User}
				if method == usecase.ForgotByToken {
					data["link"] = fmt.Sprintf("%s/auth/reset-password/%s", a.env.FrontendUrl, result.Token)
				} else {
					data["code"] = result.Code
				}
				return nil
			},
			func(ctx context.Context) error {
				return a.forgotPasswordUc.CompensateForgotPassword(ctx, usecase.CompensateForgotPassword{
					UserID: result.User.ID,
					Token:  result.Token,
					Code:   result.Code,
					Type:   method,
				})
			},
		))

		if tmpl, err = a.mailService.Mtc.GetMailTmpl(ctx, &proto_mail_template.GetMailTmplRequest{
			Id: constants.TPL_FORGOT_MAIL,
		}); err != nil {
			return err
		}

		sagaTx.AddStep(saga.NewSagaStep(
			"SendEmailForgotPassword",
			func(ctx context.Context) error {
				payload := queue.NewPayloadMail(data, []string{result.User.Email}, tmpl.MailTmpl.Id)
				if taskId, err = a.forgotPasswordUc.SendEmailForgotPassword(payload); err != nil {
					return err
				}
				return nil
			},
			func(ctx context.Context) error {
				return a.forgotPasswordUc.CompensateSendEmail(ctx, taskId)
			},
		))

		sagaTx.AddStep(saga.NewSagaStep(
			"CreateMailHistory",
			func(ctx context.Context) error {
				protoData, err := json.Marshal(&data)
				if err != nil {
					return err
				}
				if _, err := a.mailService.Mhc.CreateMailHistory(ctx, &proto_mail_history.CreateMailHistoryRequest{
					Id:            taskId,
					TemplateId:    tmpl.MailTmpl.Id,
					Subject:       tmpl.MailTmpl.Subject,
					Body:          tmpl.MailTmpl.Body,
					Tos:           []string{result.User.Email},
					Data:          string(protoData),
					EmailProvider: tmpl.MailTmpl.ProviderEmail,
				}); err != nil {
					return err
				}
				return nil
			},
			func(ctx context.Context) error {
				_, err := a.mailService.Mhc.DeleteMailHistory(ctx, &proto_mail_history.DeleteMailHistoryRequest{
					Id: taskId,
				})
				return err
			},
		))

		sagaTx.AddStep(saga.NewSagaStep(
			"CreateStatusHistory",
			func(ctx context.Context) error {
				if _, err := a.mailService.Shc.CreateStatusHistory(ctx, &proto_status_history.CreateStatusHistoryRequest{
					MailHistoryId: taskId,
					Status:        "pending",
					Message:       "Send email forgot password to " + result.User.Email,
					CreatedAt:     time.Now().Format(time.RFC3339),
				}); err != nil {
					return err
				}
				return nil
			},
			func(ctx context.Context) error {
				if _, err := a.mailService.Shc.DeleteStatusHistory(ctx, &proto_status_history.DeleteStatusHistoryRequest{
					Status:        "pending",
					MailHistoryId: taskId,
				}); err != nil {
					return err
				}
				return nil
			},
		))

		return err
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Đặt lại mật khẩu thất bại: "+err.Error())
	}

	return &proto_auth.ForgotPasswordResponse{
		User:    a.createUserInfo(result.User),
		Token:   result.Token,
		Code:    result.Code,
		Message: "Yêu cầu đặt lại mật khẩu đã được gửi. Vui lòng kiểm tra email.",
	}, nil
}

func (a *authService) createUserInfo(user entity.UserInfor) *proto_auth.UserInfo {
	userInfo := &proto_auth.UserInfo{
		Id:       user.ID,
		Email:    user.Email,
		Phone:    user.Phone,
		FullName: user.FullName,
		Avatar:   user.Avatar,
	}
	if user.Birthday != nil {
		userInfo.Birthday = timestamppb.New(*user.Birthday)
	}
	return userInfo
}
