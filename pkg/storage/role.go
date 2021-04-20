package storage

import "context"

type RoleHandler interface {
	GetRoleByID(ctx context.Context,id string)(Role,error)
	ListRole(ctx context.Context)[]Role
	CreateRole(ctx context.Context,role Role)(Role,error)
	UpdateRole(ctx context.Context,role Role)(Role,error)
	DeleteRole(ctx context.Context,id string)error
}

type Role interface {
	GetRoleName(ctx context.Context) string
	GetRoleDetails(ctx context.Context) []RoleDetail
	GetRoleID(ctx context.Context) string
}
type Verbs string

var (
	VerbsPull Verbs = "Pull"
	VerbsPush Verbs = "Push"
	VerbsAll Verbs = "*"
	VerbsCatalog Verbs = "Catalog"
)

type RoleDetail interface {
	GetGroups(ctx context.Context) []string
	GetResources(ctx context.Context) []string
	GetVerbs(ctx context.Context) []Verbs
}