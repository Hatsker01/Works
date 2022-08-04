package structs

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// UserStruct ...
type UserStruct struct {
	Id                string           `json:"id"`
	SpecId            string           `json:"specId"`
	FirstName         string           `json:"firstName"`
	LastName          string           `json:"lastName"`
	Email             string           `json:"email"`
	Password          string           `json:"-"`
	Cover             string           `json:"cover"`
	Branch            BranchStruct     `json:"branch"`
	Birthday          string           `json:"birthday"`
	Gender            string           `json:"gender"`
	AddedAt           time.Time        `json:"addedAt"`
	Role              RoleStruct       `json:"role"`
	Phone             string           `json:"phone"`
	Address           string           `json:"address"`
	WorkTime          string           `json:"workTime"`
	UserPlaceInTop    int              `json:"userPlaceInTop"`
	SocialMedias      SocialMedia      `json:"socialMedias"`
	AdditionalInforms AdditionalInform `json:"additionalInform"`
	ScoreInfo         UserScoreInfo    `json:"scoreInfo"`
	CreatedAt         time.Time        `json:"createdAt"`
	UpdatedAt         time.Time        `json:"updatedAt"`
}

type UserListResp struct {
	Id             string        `json:"id"`
	SpecId         string        `json:"specId"`
	FirstName      string        `json:"firstName"`
	LastName       string        `json:"lastName"`
	Cover          string        `json:"cover"`
	Role           RoleStruct    `json:"role"`
	UserPlaceInTop int           `json:"userPlaceInTop"`
	SocialMedias   SocialMedia   `json:"socialMedias"`
	ScoreInfo      UserScoreInfo `json:"scoreInfo"`
}

type UserLoginResp struct {
	Id        string `json:"id"`
	SpecId    string `json:"specId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"-"`
}

type UserScoreInfo struct {
	StaffAverage        float64
	ClientAverage       float64
	NumberOfStaffRated  int
	NumberOfClientRated int
}

type CreateUser struct {
	Id        string `json:"id"`
	SpecId    string `json:"specId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	BranchId  int64  `json:"branchId"`
	Gender    string `json:"gender"`
	RoleId    string `json:"roleId"`
}

type CreateUserReq struct {
	SpecId    string `json:"specId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	BranchId  int64  `json:"branchId"`
	Gender    string `json:"gender"`
	RoleId    string `json:"roleId"`
}

type UpdateUserFromUser struct {
	ID           string      `json:"id"`
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	Gender       string      `json:"gender"`
	Phone        string      `json:"phone"`
	Address      string      `json:"address"`
	Birthday     string      `json:"birthday"`
	SocialMedias SocialMedia `json:"socialMedias"`
}
type UpdateUserFromUserReq struct {
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	Gender       string      `json:"gender"`
	Phone        string      `json:"phone"`
	Address      string      `json:"address"`
	Birthday     string      `json:"birthday"`
	SocialMedias SocialMedia `json:"socialMedias"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserAuth struct {
	Id           string `json:"id"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type ByIdReq struct {
	Id string `json:"id"`
}

type SocialMedia struct {
	Instagram string `json:"instagram"`
	Telegram  string `json:"telegram"`
	Facebook  string `json:"facebook"`
	LindenIn  string `json:"lindenIn"`
	YouTube   string `json:"youTube"`
}

type AdditionalInform struct {
	Interested    string `json:"interested"`
	NotInterested string `json:"notInterested"`
	DoingStatus   string `json:"doingStatus"`
}

func (a SocialMedia) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *SocialMedia) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

const (
	LoadUserImage = `UPDATE users_client set cover = $1 where id = $2`
)

type UpdatePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
