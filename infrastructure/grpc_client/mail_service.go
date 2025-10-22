package grpc_client

import (
	"github.com/anhvanhoa/service-core/domain/grpc_client"
	proto_mail_history "github.com/anhvanhoa/sf-proto/gen/mail_history/v1"
	proto_mail_provider "github.com/anhvanhoa/sf-proto/gen/mail_provider/v1"
	proto_mail_template "github.com/anhvanhoa/sf-proto/gen/mail_tmpl/v1"
	proto_status_history "github.com/anhvanhoa/sf-proto/gen/status_history/v1"
)

type MailService struct {
	Shc proto_status_history.StatusHistoryServiceClient
	Mtc proto_mail_template.MailTmplServiceClient
	Mpc proto_mail_provider.MailProviderServiceClient
	Mhc proto_mail_history.MailHistoryServiceClient
}

func NewMailService(client *grpc_client.Client) *MailService {
	if client == nil {
		return &MailService{}
	}
	return &MailService{
		Shc: proto_status_history.NewStatusHistoryServiceClient(client.GetConnection()),
		Mtc: proto_mail_template.NewMailTmplServiceClient(client.GetConnection()),
		Mpc: proto_mail_provider.NewMailProviderServiceClient(client.GetConnection()),
		Mhc: proto_mail_history.NewMailHistoryServiceClient(client.GetConnection()),
	}
}
