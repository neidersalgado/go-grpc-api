package user

type UserRequest struct {
	UserID                int32         `json:"user_id"`
	PWDHash               string        `json:"pwd_hash"`
	Email                 string        `json:"email"`
	Name                  string        `json:"name"`
	Age                   int32         `json:"age"`
	AdditionalInformation string        `json:"additional_information"`
	Parent                []UserRequest `json:"parent"`
}

type getUserRequest struct {
	Email string `json:"email,omitempty"`
}

type authRequest struct {
	Email string `json:"email,omitempty"`
	Hash  string `json:"hash,omitempty"`
}
