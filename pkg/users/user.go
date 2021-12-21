package users

type User struct {
	UserId                int32   `json:"user_id,omitempty"`
	PwdHash               string  `json:"pwd_hash,omitempty"`
	Email                 string  `json:"email,omitempty"`
	Name                  string  `json:"name,omitempty"`
	Age                   int32   `json:"age,omitempty"`
	AdditionalInformation string  `json:"additional_information,omitempty"`
	Parents               []*User `json:"parents,omitempty"`
}

type Auth struct {
	Mail string `json:"name,omitempty"`
	Hash string `json:"hash,omitempty"`
}
