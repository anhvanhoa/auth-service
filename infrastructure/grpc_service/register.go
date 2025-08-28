package grpcservice

import (
	"auth-service/constants"
	"auth-service/domain/usecase"
	"context"
	"fmt"
	"time"

	"github.com/anhvanhoa/service-core/domain/queue"
	"github.com/anhvanhoa/service-core/domain/saga"
	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	proto_mail_history "github.com/anhvanhoa/sf-proto/gen/mail_history/v1"
	proto_mail_template "github.com/anhvanhoa/sf-proto/gen/mail_tmpl/v1"
	proto_status_history "github.com/anhvanhoa/sf-proto/gen/status_history/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *authService) Register(ctx context.Context, req *proto_auth.RegisterRequest) (*proto_auth.RegisterResponse, error) {
	if err := validatePasswordMatch(req.GetPassword(), req.GetConfirmPassword()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	existingUser, err := a.registerUc.CheckUserExist(req.GetEmail())
	if err == nil && existingUser.ID != "" {
		return nil, status.Error(codes.AlreadyExists, "Email đã được sử dụng")
	}

	var result usecase.ResRegister
	exp := time.Now().Add(15 * time.Minute)
	os := "web"
	sagaId := fmt.Sprintf("register-%s-%s", req.GetEmail(), a.uuid.Gen())
	err = a.registerUc.RegisterWithSaga(sagaId, func(ctx context.Context, sagaTx saga.SagaTransactionI) error {
		code := a.registerUc.GengerateCode(6)
		registerReq := usecase.RegisterReq{
			Email:           req.GetEmail(),
			FullName:        req.GetFullName(),
			Password:        req.GetPassword(),
			ConfirmPassword: req.GetConfirmPassword(),
			Code:            code,
		}
		var err error
		sagaTx.AddStep(
			saga.NewSagaStep(
				"Register",
				func(ctx context.Context) error {
					result, err = a.registerUc.Register(registerReq, os, exp)
					return err
				},
				func(ctx context.Context) error {
					return a.registerUc.CompensateRegister(ctx, result.UserInfor.ID, result.Token)
				},
			),
		)
		var tmpl *proto_mail_template.GetMailTmplResponse
		data := map[string]any{
			"link": a.env.FRONTEND_URL + "/auth/verify/" + result.Token,
			"user": result.UserInfor,
		}
		sagaTx.AddStep(saga.NewSagaStep(
			"GetMailTemplate",
			func(ctx context.Context) error {
				if tmpl, err = a.mailService.Mtc.GetMailTmpl(ctx, &proto_mail_template.GetMailTmplRequest{
					Id: constants.TPL_REGISTER_MAIL,
				}); err != nil {
					return err
				}
				return nil
			}, nil,
		))

		var taskId string
		sagaTx.AddStep(saga.NewSagaStep(
			"SendMail",
			func(ctx context.Context) error {
				payload := queue.NewPayload(data, []string{result.UserInfor.Email}, tmpl.MailTmpl.Id)
				if taskId, err = a.registerUc.SendMail(&payload); err != nil {
					return err
				}
				return nil
			},
			func(ctx context.Context) error {
				return a.registerUc.CompensateSendMail(ctx, taskId)
			},
		))

		protoData := make(map[string]*anypb.Any)
		for k := range data {
			if anyValue, err := anypb.New(timestamppb.New(time.Now())); err == nil {
				protoData[k] = anyValue
			}
		}

		sagaTx.AddStep(saga.NewSagaStep(
			"CreateMailHistory",
			func(ctx context.Context) error {
				if _, err := a.mailService.Mhc.CreateMailHistory(ctx, &proto_mail_history.CreateMailHistoryRequest{
					Id:            taskId,
					TemplateId:    constants.TPL_REGISTER_MAIL,
					Subject:       tmpl.MailTmpl.Subject,
					Body:          tmpl.MailTmpl.Body,
					Tos:           []string{result.UserInfor.Email},
					Data:          protoData,
					EmailProvider: tmpl.MailTmpl.ProviderEmail,
				}); err != nil {
					return err
				}
				return nil
			},
			func(ctx context.Context) error {
				a.mailService.Mhc.DeleteMailHistory(ctx, &proto_mail_history.DeleteMailHistoryRequest{
					Id: taskId,
				})
				return nil
			},
		))

		sagaTx.AddStep(saga.NewSagaStep(
			"CreateStatusHistory",
			func(ctx context.Context) error {
				if _, err := a.mailService.Shc.CreateStatusHistory(ctx, &proto_status_history.CreateStatusHistoryRequest{
					MailHistoryId: taskId,
					Status:        "pending",
					Message:       "Send mail to " + result.UserInfor.Email,
					CreatedAt:     time.Now().Format(time.RFC3339),
				}); err != nil {
					return err
				}
				return nil
			},
			func(ctx context.Context) error {
				a.mailService.Shc.DeleteStatusHistory(ctx, &proto_status_history.DeleteStatusHistoryRequest{
					Status:        "pending",
					MailHistoryId: taskId,
				})
				return nil
			},
		))

		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Đăng ký thất bại: "+err.Error())
	}

	userInfo := &proto_auth.UserInfo{
		Id:       result.UserInfor.ID,
		Email:    result.UserInfor.Email,
		Phone:    result.UserInfor.Phone,
		FullName: result.UserInfor.FullName,
		Avatar:   result.UserInfor.Avatar,
		Bio:      result.UserInfor.Bio,
		Address:  result.UserInfor.Address,
	}

	if result.UserInfor.Birthday != nil {
		userInfo.Birthday = timestamppb.New(*result.UserInfor.Birthday)
	}

	return &proto_auth.RegisterResponse{
		User:    userInfo,
		Token:   result.Token,
		Message: "Đăng ký thành công. Vui lòng kiểm tra email để xác thực tài khoản.",
	}, nil
}
