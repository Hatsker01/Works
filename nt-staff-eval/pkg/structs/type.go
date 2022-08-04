package structs

const (
	TypeUserImage    = 1
	TypeSectionImage = 2
)

type Data struct {
	Data    string `json:"data"`
	Message string `json:"message"`
}
