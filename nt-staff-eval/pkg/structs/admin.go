package structs

type AdminLogin struct {
	Login    string `json:"email"`
	Password string `json:"password"`
}

type AdminStruct struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	Password     string `json:"password"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
