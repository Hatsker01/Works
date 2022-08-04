package structs

import "time"

type RoleStruct struct {
	Id        string        `json:"id"`
	Name      string        `json:"name"`
	Section   SectionStruct `json:"section"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

type CreateRoleStruct struct {
	Name      string `json:"name"`
	SectionId int    `json:"sectionId"`
}
