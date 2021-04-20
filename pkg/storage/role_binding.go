package storage

import "context"

type RoleBindingHandler interface {
	GetRoleBindingByID(ctx context.Context,roleBindingID string) RoleBinding
	GetRoleBindingByRoleID(ctx context.Context,roleID string) []RoleBinding
	GetRoleBindingByUserID(ctx context.Context,userID string) []RoleBinding
	CreateRoleBinding(ctx context.Context,roleBinding RoleBinding)(RoleBinding,error)
	UpdateRoleBinding(ctx context.Context,roleBinding RoleBinding)(RoleBinding,error)
	DeleteRoleBinding(ctx context.Context,roleBindingID string)(error)
	ListRoleBinding(ctx context.Context) ([]RoleBinding,error)
}

type RoleBinding interface {
	GetName(ctx context.Context)string
	GetID(ctx context.Context)string
	GetRoleIDs(ctx context.Context) []string
	GetUserIDs(ctx context.Context) []string
}
