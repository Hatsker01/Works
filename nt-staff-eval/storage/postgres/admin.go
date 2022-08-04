package postgres

import (
	"database/sql"

	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/repo"
	"github.com/jmoiron/sqlx"
)

type adminRepository struct {
	db *sqlx.DB
}

func NewAdminRepo(db *sqlx.DB) repo.AdminRepoInterface {
	return &adminRepository{
		db: db,
	}
}

func (a adminRepository) Login(admin structs.AdminLogin) (structs.AdminStruct, error) {
	var resp structs.AdminStruct
	err := a.db.QueryRow(`select id, login, password, access_token, refresh_token from admin
	where login=$1 and password=$2`, admin.Login, admin.Password).
		Scan(&resp.ID,
			&resp.Login,
			&resp.Password,
			&resp.AccessToken,
			&resp.RefreshToken,
		)
	if err != nil {
		return structs.AdminStruct{}, err
	}

	return resp, nil
}

func (a adminRepository) Update(admin structs.AdminStruct) error {
	result, err := a.db.Exec(`update admin set access_token=$1, refresh_token=$2 where id=$3`,
		admin.AccessToken,
		admin.RefreshToken,
		admin.ID)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
