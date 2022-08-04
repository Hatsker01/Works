package repo

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"

	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
)

type UserRepoInterface interface {
	CreateUser(user structs.CreateUser) (structs.UserStruct, error)
	GetUser(id string) (structs.UserStruct, error)
	GetListUsers(filters map[string]string, page, limit int) ([]structs.UserListResp, int, error)
	GetTopUsers(page, limit int) ([]structs.UserStruct, int, error)
	UpdateUser(user structs.UpdateUserFromUser) (structs.UserStruct, error)
	StaffAverage(id string) (float64, error)
	ClientAverage(id string) (float64, error)
	LoginUser(user structs.LoginUser) (structs.UserLoginResp, error)
	UpdateToken(user structs.UserAuth) error
	ChangePassword(user structs.UserStruct, password structs.UpdatePassword) (changed bool, err error)
	HandleToken(token string) (user structs.UserStruct, err error)
}

type ImageRepoInterface interface {
	LoadImage(imageType int8, path, id string) (err error)
}

type AuthRepoInterface interface {
	CreateToken(userId string) (*structs.TokenDetails, error)
	CreateAuth(userId string, td *structs.TokenDetails) error
	ExtractToken(r *http.Request) string
	VerifyToken(r *http.Request) (*jwt.Token, error)
	TokenValid(r *http.Request) error
	ExtractTokenMetadata(r *http.Request) (*structs.AccessDetails, error)
	FetchAuth(authD *structs.AccessDetails) (int64, error)
	DeleteAuth(givenUUID string) (int64, error)
}

type SectionRepoInterface interface {
	CreateSection(section structs.SectionStruct) (structs.SectionStruct, error)
	GetSection(id int) (structs.SectionStruct, error)
	GetListSections() ([]structs.SectionStruct, int, error)
	UpdateSection(section structs.SectionStruct) (structs.SectionStruct, error)
	DeleteSection(id int) error
}

type RoleRepoInterface interface {
	CreateRole(role structs.RoleStruct) (structs.RoleStruct, error)
	GetRole(id string) (structs.RoleStruct, error)
	GetListRoles(page, limit int) ([]structs.RoleStruct, int, error)
	UpdateRole(role structs.RoleStruct) (structs.RoleStruct, error)
	DeleteRole(id string) error
}

//type BehaviorRepoInterface interface {
//	CreateBehavior(behavior structs.BehaviorStruct) (structs.BehaviorStruct, error)
//	GetBehavior(id string) (structs.BehaviorStruct, error)
//	GetListBehaviors(page, limit int) ([]structs.BehaviorStruct, int, error)
//	UpdateBehavior(behavior structs.BehaviorStruct) (structs.BehaviorStruct, error)
//	DeleteBehavior(id string) error
//}

type EvaluationRepoInterface interface {
	CreateEvaluation(evaluation structs.EvaluationStruct) (structs.EvaluationStruct, error)
	GetEvaluation(id string) (structs.EvaluationStruct, error)
	GetListEvaluations() ([]structs.EvaluationStruct, int, error)
	UpdateEvaluation(evaluation structs.UpdateEvaluation) (structs.EvaluationStruct, error)
	DeleteEvaluation(id string) error
}

type RatedRepoInterface interface {
	CreateRated(rated structs.CreateRated) error
	GetListRateds(page, limit int) ([]structs.Rated, int, error)
	DeleteRated(id string) error
}

type BranchRepoInterface interface {
	CreateBranch(branch structs.CreateBranch) (structs.BranchStruct, error)
	GetBranch(id int64) (structs.BranchStruct, error)
	GetListBranch() ([]structs.BranchStruct, int, error)
	UpdateBranch(branch structs.BranchStruct) (structs.BranchStruct, error)
	DeleteBranch(id int64) error
}

type AdminRepoInterface interface {
	Login(admin structs.AdminLogin) (structs.AdminStruct, error)
	Update(admin structs.AdminStruct) error
}

type SuggestRepoInterface interface {
	CreateSuggest(suggest structs.CreateSuggest) (structs.Suggest, error)
	GetSuggest(id string) (structs.Suggest, error)
	GetListSuggests(filters map[string]string, page, limit int) ([]structs.Suggest, int, error)
	UpdateStatusSuggest(req structs.UpdateStatusSuggestReq) (structs.Suggest, error)
	DeleteSuggest(id string) error
}

type NewsRepoInterface interface {
	CreateNews(news structs.CreateNews) (structs.News, error)
	GetNews(id string) (structs.News, error)
	GetListNews(filters map[string]string, page, limit int) ([]structs.News, int, error)
	UpdateNews(news structs.CreateNews) (string, error)
	DeleteNews(id string) error
}

type NewsCategoryRepoInterface interface {
	CreateNewsCategory(newsCategory structs.Category) (structs.Category, error)
	GetNewsCategory(id string) (structs.Category, error)
	GetListNewsCategory() ([]structs.Category, int, error)
	UpdateNewsCategory(newsCategory structs.Category) (structs.Category, error)
	DeleteNewsCategory(id string) error
}
