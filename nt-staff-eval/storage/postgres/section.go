package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/repo"
)

type sectionRepository struct {
	db *sqlx.DB
}

func NewSectionRepo(db *sqlx.DB) repo.SectionRepoInterface {
	return &sectionRepository{
		db: db,
	}
}

func (s sectionRepository) CreateSection(section structs.SectionStruct) (structs.SectionStruct, error) {
	err := s.db.QueryRow(`INSERT INTO sections(name, cover, spec_id)
	 VALUES ($1, $2, $3) returning id`, section.Name, section.Cover, section.SpecId).Scan(&section.Id)
	if err != nil {
		return structs.SectionStruct{}, err
	}

	section, err = s.GetSection(section.Id)
	if err != nil {
		return structs.SectionStruct{}, err
	}

	return section, nil
}

func (s sectionRepository) GetSection(id int) (structs.SectionStruct, error) {
	var section structs.SectionStruct
	fmt.Println(id)
	err := s.db.QueryRow(`select  id, name, cover, spec_id, created_at, updated_at from sections
	where deleted_at is null and id=$1`, id).
		Scan(&section.Id,
			&section.Name,
			&section.Cover,
			&section.SpecId,
			&section.CreatedAt,
			&section.UpdatedAt)
		fmt.Println(section)
	if err != nil {
		return structs.SectionStruct{}, err
	}

	return section, nil
}

func (s sectionRepository) GetListSections() ([]structs.SectionStruct, int, error) {
	//offset := (page - 1) * limit
	rows, err := s.db.Queryx(`
		SELECT id, name, cover, spec_id, created_at, updated_at FROM sections WHERE deleted_at IS NULL order by spec_id
		`)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var (
		sections []structs.SectionStruct
		count    int
	)

	for rows.Next() {
		var section structs.SectionStruct
		err = rows.Scan(
			&section.Id,
			&section.Name,
			&section.Cover,
			&section.SpecId,
			&section.CreatedAt,
			&section.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		if err != nil {
			return nil, 0, err
		}

		sections = append(sections, section)
	}

	err = s.db.QueryRow(`SELECT count(*) FROM sections WHERE deleted_at IS NULL`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return sections, count, nil
}

func (s sectionRepository) UpdateSection(section structs.SectionStruct) (structs.SectionStruct, error) {
	result, err := s.db.Exec(`UPDATE sections SET name=$1, cover=$2, spec_id=$3, updated_at=$3 WHERE id=$4`,
		&section.Name,
		&section.Cover,
		&section.SpecId,
		time.Now().UTC(),
		&section.Id)
	if err != nil {
		return structs.SectionStruct{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return structs.SectionStruct{}, sql.ErrNoRows
	}

	section, err = s.GetSection(section.Id)
	if err != nil {
		return structs.SectionStruct{}, err
	}

	return section, err
}

func (s sectionRepository) DeleteSection(id int) error {
	result, err := s.db.Exec(`UPDATE sections SET deleted_at = $1 WHERE id = $2`, time.Now().UTC(), id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
