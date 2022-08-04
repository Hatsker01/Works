package structs

import "time"

type BranchStruct struct {
	Id        int64     `json:"Id"`
	Name      string    `json:"name"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateBranch struct {
	Name string `json:"name"`
	City string `json:"city"`
}

type UpdateBranch struct {
	Name string `json:"name"`
	City string `json:"city"`
}
