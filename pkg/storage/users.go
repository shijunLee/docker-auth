package storage

import "context"

type UserHandler interface {
	UserAuthHandler
	Create(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, userID string) error
	Update(ctx context.Context, user User) (User, error)
	List(ctx context.Context) ([]User, error)
}

type UserAuthHandler interface {
	GetUserByID(ctx context.Context, userID string) (User, error)
	AuthUser(ctx context.Context, userID, password string) bool
}

type User interface {
	GetUserID(ctx context.Context) string
	//GetUserName this is a disploy name for user
	GetUserName(ctx context.Context) string
	//UserGroups this is return user groups for namespace
	UserGroups(ctx context.Context) []string
}
