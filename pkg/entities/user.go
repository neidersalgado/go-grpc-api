package entities

type User struct {
	UserID                int32
	pwdHash               string
	Email                 string
	Name                  string
	Age                   int32
	AdditionalInformation string
	Parent                []User
}
