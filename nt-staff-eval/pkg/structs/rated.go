package structs

type Rated struct {
	Id          string             `json:"id"`
	User        UserStruct         `json:"user"`
	Additional  string             `json:"additional"`
	Evaluations []EvaluationStruct `json:"evaluations"`
	IsStaff     bool               `json:"isStaff"`
	CreatedAt   string             `json:"createdAt"`
}

type CreateRated struct {
	Id            string    `json:"id"`
	UserId        string    `json:"userId"`
	Additional    string    `json:"additional"`
	IsStaff       bool      `json:"isStaff"`
	EvaluationsId []ByIdReq `json:"evaluations"`
}

type CreateRatedReq struct {
	UserId        string    `json:"userId"`
	Additional    string    `json:"additional"`
	IsStaff       bool      `json:"isStaff"`
	EvaluationsId []ByIdReq `json:"evaluations"`
}
