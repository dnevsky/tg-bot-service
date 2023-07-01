package models

type Service struct {
	Id            int    `json:"-"`
	UserId        int    `json:"userid"`
	NameOfService string `json:"nameofservice"`
	Login         string `json:"login"`
	Password      string `json:"password"`
}
