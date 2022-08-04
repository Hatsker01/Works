package postgres

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/repo"
)

type ratedRepository struct {
	db *sqlx.DB
}

func NewRatedRepo(db *sqlx.DB) repo.RatedRepoInterface {
	return &ratedRepository{
		db: db,
	}
}

func (r ratedRepository) CreateRated(rated structs.CreateRated) error {
	var id string
	err := r.db.QueryRow(`
			insert into rated(id, additional, user_id, is_staff, created_at) 
			values ($1, $2, $3, $4, $5) returning id`, rated.Id, rated.Additional, rated.UserId, rated.IsStaff, time.Now().UTC()).Scan(&id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`
		update users_client set updated_evoluation=$1 where id=$2
	`, time.Now().UTC(), rated.UserId)
	if err != nil {
		return err
	}

	for _, j := range rated.EvaluationsId {
		_, err = r.db.Exec(`
			insert into rated_evaluations(evaluation_id, rated_id) 
			values ($1, $2) returning id`, j.Id, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r ratedRepository) GetListRateds(page, limit int) ([]structs.Rated, int, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Queryx(`
			select id, additional, user_id, is_staff, created_at from rated
			where deleted_at is null limit $1 offset $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var (
		count  int
		rateds []structs.Rated
	)

	for rows.Next() {
		var rated structs.Rated
		var additional sql.NullString

		err = rows.Scan(&rated.Id, &additional, &rated.User.Id, &rated.IsStaff, &rated.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		if additional.Valid {
			rated.Additional = additional.String
		}

		rowsEvaluation, err := r.db.Queryx(`select evaluation_id from rated_evaluations where rated_id = $1`, rated.Id)
		if err != nil {
			return nil, 0, err
		}

		var evaluations []structs.EvaluationStruct

		for rowsEvaluation.Next() {
			var evaluation structs.EvaluationStruct
			var id string

			err = rowsEvaluation.Scan(&id)
			if err != nil {
				return nil, 0, err
			}
			evaluation, err = NewEvaluationRepo(r.db).GetEvaluation(id)
			if err != nil {
				return nil, 0, err
			}

			evaluations = append(evaluations, evaluation)
		}
		rated.Evaluations = evaluations
		rateds = append(rateds, rated)
	}

	err = r.db.QueryRow("SELECT count(*) FROM rated WHERE deleted_at IS NULL").Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return rateds, count, nil
}

func (r ratedRepository) DeleteRated(id string) error {
	result, err := r.db.Exec(`UPDATE rated SET deleted_at = $1 WHERE id=$2`, time.Now().UTC(), id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// func (r ratedRepository)AvarageStar(id string)string{
// 	err:=r.db.QueryRow(`
// 	select posts.descripion, count(comments.post_id) as amount from posts,comments group by (posts.descripion);
// 	SELECT user_client,count(rated.id)
// 	`)
// }
