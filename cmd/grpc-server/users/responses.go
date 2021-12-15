package grpc

type createUserResponse struct {
	Error error
	Id    string
}

type getUserResponse struct {
	UserResponse
}

type UserResponse struct {
	UserId                int32          `json:"user_id,omitempty"`
	PwdHash               string         `json:"pwd_hash,omitempty"`
	Email                 string         `json:"email,omitempty"`
	Name                  string         `json:"name,omitempty"`
	Age                   int32          `json:"age,omitempty"`
	AdditionalInformation string         `json:"additional_information,omitempty"`
	Parents               []*UserRequest `json:"parents,omitempty"`
}