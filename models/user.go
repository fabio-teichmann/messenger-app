package models

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ProfilePic bool   `json:"profile_pic"`
}
