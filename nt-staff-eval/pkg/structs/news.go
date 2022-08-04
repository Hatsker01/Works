package structs

type News struct {
	Id        string     `json:"id"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Cover     string     `json:"cover"`
	Author    UserStruct `json:"author"`
	ReadTime  string     `json:"readTime"`
	Category  Category   `json:"category"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt string     `json:"updatedAt"`
}

type CreateNews struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Cover      string `json:"cover"`
	AuthorId   string `json:"authorId"`
	ReadTime   string `json:"readTime"`
	CategoryId string `json:"categoryId"`
}

type CreateNewsReq struct {
	Title      string `json:"title"`
	Body       string `json:"body"`
	Cover      string `json:"cover"`
	AuthorId   string `json:"authorId"`
	ReadTime   string `json:"readTime"`
	CategoryId string `json:"categoryId"`
}

type UpdateNewsReq struct {
	Title      string `json:"title"`
	Body       string `json:"body"`
	Cover      string `json:"cover"`
	AuthorId   string `json:"authorId"`
	ReadTime   string `json:"readTime"`
	CategoryId string `json:"categoryId"`
}

type Category struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
