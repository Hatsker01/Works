package structs

type Suggest struct {
	Id        string     `json:"id"`
	User      UserStruct `json:"user"`
	Content   string     `json:"content"`
	Status    string     `json:"status"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt string     `json:"updatedAt"`
}

type CreateSuggest struct {
	Id      string `json:"id"`
	UserId  string `json:"userId"`
	Content string `json:"content"`
}

type CreateSuggestReq struct {
	UserId  string `json:"userId"`
	Content string `json:"content"`
}

type UpdateStatusSuggest struct {
	Status string `json:"status"`
}

type UpdateStatusSuggestReq struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}
