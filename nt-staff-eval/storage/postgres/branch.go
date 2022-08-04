package postgres

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/repo"
)

type branchRepository struct {
	db *sqlx.DB
}

func NewBranchRepo(db *sqlx.DB) repo.BranchRepoInterface {
	return branchRepository{
		db: db,
	}
}

func (b branchRepository) CreateBranch(branch structs.CreateBranch) (structs.BranchStruct, error) {
	var id int64
	err := b.db.QueryRow(`INSERT INTO branches(name, city) VALUES($1, $2) returning id`,
		branch.Name, branch.City).Scan(&id)
	if err != nil {
		return structs.BranchStruct{}, err
	}

	NewBranch, err := b.GetBranch(id)
	if err != nil {
		return structs.BranchStruct{}, err
	}

	return NewBranch, nil
}

func (b branchRepository) GetBranch(id int64) (structs.BranchStruct, error) {
	var branch structs.BranchStruct
	err := b.db.QueryRow(`select id, name, city, created_at, updated_at from branches
	where deleted_at is null and id=$1`, id).
		Scan(&branch.Id,
			&branch.Name,
			&branch.City,
			&branch.CreatedAt,
			&branch.UpdatedAt)
	if err != nil {
		return structs.BranchStruct{}, err
	}

	return branch, nil
}

func (b branchRepository) GetListBranch() ([]structs.BranchStruct, int, error) {

	rows, err := b.db.Queryx(`
		select id, name, city, created_at, updated_at from branches
		where deleted_at is null`)
	if err != nil {
		return nil, 0, nil
	}

	var (
		count    int
		branches []structs.BranchStruct
	)

	for rows.Next() {
		var branch structs.BranchStruct
		err = rows.Scan(&branch.Id, &branch.Name, &branch.City, &branch.CreatedAt, &branch.UpdatedAt)
		if err != nil {
			return nil, 0, nil
		}

		branches = append(branches, branch)
	}

	err = b.db.QueryRow(`
		select count(*) from branches where deleted_at is null
	`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return branches, count, nil
}

func (b branchRepository) UpdateBranch(branch structs.BranchStruct) (structs.BranchStruct, error) {
	result, err := b.db.Exec(`UPDATE branches SET name=$1, city=$2, updated_at=$3 WHERE id=$4`,
		&branch.Name,
		&branch.City,
		time.Now().UTC(),
		&branch.Id)
	if err != nil {
		return structs.BranchStruct{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return structs.BranchStruct{}, sql.ErrNoRows
	}

	branch, err = b.GetBranch(branch.Id)
	if err != nil {
		return structs.BranchStruct{}, err
	}

	return branch, err
}

func (b branchRepository) DeleteBranch(id int64) error {
	result, err := b.db.Exec(`UPDATE branches SET deleted_at = $1 WHERE id = $2`, time.Now().UTC(), id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
