package structs

import "time"

type LikeStruct struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UsersLike struct {
	Id     string `json:"id"`
	IsLike bool   `json:"isLike"`
	UserID int    `json:"userId"`
	LikeId string `json:"likeId"`
}
