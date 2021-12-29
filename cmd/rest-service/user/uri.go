package user

import "fmt"

var (
	UsersBaseUri = "/users"
	PostUser     = fmt.Sprintf("%s", UsersBaseUri)
	GetUser      = fmt.Sprintf("%s/{%s}", UsersBaseUri, Email)
	DeleteUser   = fmt.Sprintf("%s/{%s}", UsersBaseUri, Email)
	UpdateUser   = fmt.Sprintf("%s/{%s}", UsersBaseUri, Email)
	AuthUser     = UsersBaseUri
)
