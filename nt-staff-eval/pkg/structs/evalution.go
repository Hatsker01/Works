package structs

type EvaluationStruct struct {
	Id        string        `json:"id"`
	Content   string        `json:"content"`
	Section   SectionStruct `json:"section"`
	Type      string        `json:"type"`
	Star      int           `json:"star"`
	CreatedAt string        `json:"createdAt"`
	UpdatedAt string        `json:"updatedAt"`
}

type UpdateEvaluation struct {
	Id        string `json:"id"`
	Content   string `json:"content"`
	SectionId int    `json:"sectionId"`
	Type      string `json:"type"`
	Star      int    `json:"star"`
}

type UpdateEvaluationReq struct {
	Content   string `json:"content"`
	SectionId int    `json:"sectionId"`
	Type      string `json:"type"`
	Star      int    `json:"star"`
}
