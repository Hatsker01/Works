package structs

import (
	"time"
)

type SectionStruct struct {
	Id        int       `json:"id"`
	SpecId    int       `json:"specId"`
	Name      string    `json:"name"`
	Cover     string    `json:"cover"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateSectionStruct struct {
	Name   string `json:"name"`
	Cover  string `json:"cover"`
	SpecId int    `json:"specId"`
}

const (
	LoadSectionImage = `UPDATE sections set cover = $1 where id = $2`
)
