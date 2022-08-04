package postgres

import (
	"time"

	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/repo"
	"github.com/jmoiron/sqlx"
)

type suggestRepository struct {
	db *sqlx.DB
}

func NewSuggestRepo(db *sqlx.DB) repo.SuggestRepoInterface {
	return &suggestRepository{
		db: db,
	}
}

func (s suggestRepository) CreateSuggest(suggest structs.CreateSuggest) (structs.Suggest, error) {
	err := s.db.QueryRow(`INSERT INTO suggests(id, user_id, content) 
			VALUES ($1, $2, $3) returning id`, suggest.Id, suggest.UserId, suggest.Content).Scan(&suggest.Id)
	if err != nil {
		return structs.Suggest{}, err
	}

	suggestNew, err := s.GetSuggest(suggest.Id)
	if err != nil {
		return structs.Suggest{}, err
	}

	return suggestNew, nil
}

func (s suggestRepository) GetSuggest(id string) (structs.Suggest, error) {
	var suggest structs.Suggest
	err := s.db.QueryRow(`select  id, user_id, content, status, created_at, updated_at from suggests
		where deleted_at is null and id=$1`, id).
		Scan(&suggest.Id,
			&suggest.User.Id,
			&suggest.Content,
			&suggest.Status,
			&suggest.CreatedAt,
			&suggest.UpdatedAt)
	if err != nil {
		return structs.Suggest{}, err
	}

	return suggest, nil
}

func (s suggestRepository) GetListSuggests(filters map[string]string, page, limit int) ([]structs.Suggest, int, error) {
	offset := (page - 1) * limit

	query := `SELECT id, user_id, content, status, created_at, updated_at FROM suggests WHERE deleted_at IS NULL`
	if filters["status"] == "new" || filters["status"] == "active" || filters["status"] == "inactive" {
		query += ` AND status = '` + filters["status"] + "' "
	}
	if filters["user_id"] != "" {
		query += ` AND user_id = '` + filters["user_id"] + `' `
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := s.db.Queryx(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var suggests []structs.Suggest
	for rows.Next() {
		var suggest structs.Suggest
		err := rows.Scan(&suggest.Id,
			&suggest.User.Id,
			&suggest.Content,
			&suggest.Status,
			&suggest.CreatedAt,
			&suggest.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}

		suggests = append(suggests, suggest)
	}
	var count int
	err = s.db.QueryRow(`SELECT count(*) FROM suggests WHERE deleted_at IS NULL`).Scan(&count)

	return suggests, count, nil
}

func (s suggestRepository) UpdateStatusSuggest(req structs.UpdateStatusSuggestReq) (structs.Suggest, error) {
	err := s.db.QueryRow(`UPDATE suggests SET status = $1, updated_at = $2 WHERE id = $3 returning id`,
		req.Status, time.Now(), req.Id).Scan(&req.Id)
	if err != nil {
		return structs.Suggest{}, err
	}

	suggestNew, err := s.GetSuggest(req.Id)
	if err != nil {
		return structs.Suggest{}, err
	}

	return suggestNew, nil
}

func (s suggestRepository) DeleteSuggest(id string) error {
	_, err := s.db.Exec(`UPDATE suggests SET deleted_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}
