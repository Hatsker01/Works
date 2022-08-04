package postgres

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/repo"
)

type roleRepository struct {
	db *sqlx.DB
}

func NewRoleRepo(db *sqlx.DB) repo.RoleRepoInterface {
	return &roleRepository{
		db: db,
	}
}

func (r roleRepository) CreateRole(role structs.RoleStruct) (structs.RoleStruct, error) {
	err := r.db.QueryRow(`INSERT INTO roles(id, name, section_id)
	 VALUES ($1, $2, $3) returning id`, role.Id, role.Name, role.Section.Id).Scan(&role.Id)
	if err != nil {
		return structs.RoleStruct{}, err
	}

	role, err = r.GetRole(role.Id)
	if err != nil {
		return structs.RoleStruct{}, err
	}

	return role, nil
}

func (r roleRepository) GetRole(id string) (structs.RoleStruct, error) {
	var role structs.RoleStruct
	err := r.db.QueryRow(`select  id, name, section_id, created_at, updated_at from roles
	where deleted_at is null and id=$1`, id).
		Scan(&role.Id,
			&role.Name,
			&role.Section.Id,
			&role.CreatedAt,
			&role.UpdatedAt)
	if err != nil {
		return structs.RoleStruct{}, err
	}

	return role, nil
}

func (r roleRepository) GetListRoles(page, limit int) ([]structs.RoleStruct, int, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Queryx(`
		SELECT id, name, section_id, created_at, updated_at FROM roles WHERE deleted_at IS NULL 
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var (
		roles []structs.RoleStruct
		count int
	)

	for rows.Next() {
		var role structs.RoleStruct
		err = rows.Scan(
			&role.Id,
			&role.Name,
			&role.Section.Id,
			&role.CreatedAt,
			&role.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		roles = append(roles, role)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM roles WHERE deleted_at IS NULL`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return roles, count, nil
}

func (r roleRepository) UpdateRole(role structs.RoleStruct) (structs.RoleStruct, error) {
	result, err := r.db.Exec(`UPDATE roles SET name=$1, section_id=$2, updated_at=$3 WHERE id=$4`,
		&role.Name,
		&role.Section.Id,
		time.Now().UTC(),
		&role.Id)
	if err != nil {
		return structs.RoleStruct{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return structs.RoleStruct{}, sql.ErrNoRows
	}

	role, err = r.GetRole(role.Id)
	if err != nil {
		return structs.RoleStruct{}, err
	}

	return role, err
}

func (r roleRepository) DeleteRole(id string) error {
	result, err := r.db.Exec(`UPDATE roles SET deleted_at = $1 WHERE id = $2`, time.Now().UTC(), id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
