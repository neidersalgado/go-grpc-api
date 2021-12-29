package user

type CreateUserResponse struct {
	Err error `json:"err,omitempty"`
}

type UserResponse struct {
	PWDHash               string `json:"pwd_hash"`
	Email                 string `json:"email"`
	Name                  string `json:"name"`
	Age                   int32  `json:"age"`
	AdditionalInformation string `json:"additional_information"`
}

type DeleteUserResponse struct {
	Err error  `json:"err,omitempty"`
	Msg string `json:"msg,omitempty"`
}

type AuthResponse struct {
	Err error  `json:"err,omitempty"`
	Msg string `json:"msg,omitempty"`
}

type UpdateUserResponse struct {
	Err error `json:"err,omitempty"`
}
