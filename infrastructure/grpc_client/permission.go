package grpc_client

import (
	gc "github.com/anhvanhoa/service-core/domain/grpc_client"
	"github.com/anhvanhoa/service-core/domain/oops"
	proto_permission "github.com/anhvanhoa/sf-proto/gen/permission/v1"
	proto_user_role "github.com/anhvanhoa/sf-proto/gen/user_role/v1"
)

var (
	ErrPermissionClientNotAvailable = oops.New("Dịch vụ quyền không khả dụng")
)

type PermissionClient interface {
	PermissionService() proto_permission.PermissionServiceClient
	UserRoleService() proto_user_role.UserRoleServiceClient
}

type PermissionClientImpl struct {
	PermissionServiceClient proto_permission.PermissionServiceClient
	UserRoleServiceClient   proto_user_role.UserRoleServiceClient
}

func NewPermissionClient(client *gc.Client) (PermissionClient, error) {
	if client == nil {
		return nil, ErrPermissionClientNotAvailable
	}
	return &PermissionClientImpl{
		UserRoleServiceClient:   proto_user_role.NewUserRoleServiceClient(client.GetConnection()),
		PermissionServiceClient: proto_permission.NewPermissionServiceClient(client.GetConnection()),
	}, nil
}

func (p *PermissionClientImpl) PermissionService() proto_permission.PermissionServiceClient {
	return p.PermissionServiceClient
}

func (p *PermissionClientImpl) UserRoleService() proto_user_role.UserRoleServiceClient {
	return p.UserRoleServiceClient
}
