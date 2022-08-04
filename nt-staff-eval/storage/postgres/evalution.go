package postgres

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/repo"
)

type evaluationRepository struct {
	db *sqlx.DB
}

func NewEvaluationRepo(db *sqlx.DB) repo.EvaluationRepoInterface {
	return &evaluationRepository{
		db: db,
	}
}

func (e evaluationRepository) CreateEvaluation(evaluation structs.EvaluationStruct) (structs.EvaluationStruct, error) {
	var id string

	if evaluation.Star > 5 {
		return structs.EvaluationStruct{}, nil
	}

	err := e.db.QueryRow(`
	insert into evaluations(id, content, star, section_id, eval_type, created_at, updated_at)
	values($1, $2, $3, $4, $5) returning id`,
		evaluation.Id,
		evaluation.Content,
		evaluation.Star,
		evaluation.Section.Id,
		evaluation.Type,
		time.Now().UTC(),
		time.Now().UTC(),
	).Scan(&id)
	if err != nil {
		return structs.EvaluationStruct{}, nil
	}

	NewEvaluation, err := e.GetEvaluation(id)
	if err != nil {
		return structs.EvaluationStruct{}, nil
	}

	return NewEvaluation, nil
}

func (e evaluationRepository) GetEvaluation(id string) (structs.EvaluationStruct, error) {
	var (
		evaluation structs.EvaluationStruct
		sectionId  sql.NullInt64
	)
	err := e.db.QueryRow(`
		select id, content, star, section_id, eval_type, created_at, updated_at from evaluations
		where id = $1 and deleted_at is null`, id).Scan(
		&evaluation.Id,
		&evaluation.Content,
		&evaluation.Star,
		&sectionId,
		&evaluation.Type,
		&evaluation.CreatedAt,
		&evaluation.UpdatedAt)
	if err != nil {
		return structs.EvaluationStruct{}, err
	}
	if sectionId.Valid {
		evaluation.Section.Id = int(sectionId.Int64)
	}
	return evaluation, nil
}

func (e evaluationRepository) GetListEvaluations() ([]structs.EvaluationStruct, int, error) {
	//offset := (page - 1) * limit

	rows, err := e.db.Queryx(`
		select id, content, star, section_id, eval_type, created_at, updated_at from evaluations
		where deleted_at is null`)
	if err != nil {
		return nil, 0, nil
	}

	var (
		count       int
		evaluations []structs.EvaluationStruct
		sectionId   sql.NullInt64
	)
	for rows.Next() {
		var evaluation structs.EvaluationStruct
		err = rows.Scan(
			&evaluation.Id,
			&evaluation.Content,
			&evaluation.Star,
			&sectionId,
			&evaluation.Type,
			&evaluation.CreatedAt,
			&evaluation.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		if sectionId.Valid {
			evaluation.Section.Id = int(sectionId.Int64)
		}
		evaluations = append(evaluations, evaluation)
	}

	err = e.db.QueryRow(`
		select count(*) from evaluations where deleted_at is null
	`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return evaluations, count, nil
}

func (e evaluationRepository) UpdateEvaluation(evaluation structs.UpdateEvaluation) (structs.EvaluationStruct, error) {
	result, err := e.db.Exec(`update evaluations set content = $1, star=$2, section_id=$3, eval_type=$4, updated_at=$5
		where id=$6 and deleted_at is null`,
		evaluation.Content,
		evaluation.Star,
		evaluation.SectionId,
		evaluation.Type,
		time.Now().UTC(),
		evaluation.Id,
	)
	if err != nil {
		return structs.EvaluationStruct{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return structs.EvaluationStruct{}, sql.ErrNoRows
	}

	var NewEvaluation structs.EvaluationStruct

	NewEvaluation, err = e.GetEvaluation(evaluation.Id)

	if err != nil {
		return structs.EvaluationStruct{}, err
	}

	return NewEvaluation, nil
}

func (e evaluationRepository) DeleteEvaluation(id string) error {
	result, err := e.db.Exec("update evaluations set deleted_at = $1 where id=$2", time.Now().UTC(), id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
