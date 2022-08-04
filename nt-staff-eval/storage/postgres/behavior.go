package postgres

//
//import (
//	"database/sql"
//	"time"
//
//	"github.com/jmoiron/sqlx"
//
//	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
//	"github.com/Hatsker01/nt-staff-eval/storage/repo"
//)
//
//type behaviorRepository struct {
//	db *sqlx.DB
//}
//
//func NewBehaviorRepo(db *sqlx.DB) repo.BehaviorRepoInterface {
//	return &behaviorRepository{
//		db: db,
//	}
//}
//
//func (b behaviorRepository) CreateBehavior(behavior structs.BehaviorStruct) (structs.BehaviorStruct, error) {
//	var id string
//	err := b.db.QueryRow(`INSERT INTO behaviors(id, name, created_at, updated_at)
//	 VALUES ($1, $2, $3, $4) returning id`, behavior.Id, behavior.Name, time.Now().UTC(), time.Now().UTC()).Scan(&id)
//	if err != nil {
//		return structs.BehaviorStruct{}, err
//	}
//
//	NewBehavior, err := b.GetBehavior(id)
//	if err != nil {
//		return structs.BehaviorStruct{}, err
//	}
//
//	return NewBehavior, nil
//}
//
//func (b behaviorRepository) GetBehavior(id string) (structs.BehaviorStruct, error) {
//	var behavior structs.BehaviorStruct
//	err := b.db.QueryRow(`select  id, name, created_at, updated_at from behaviors
//	where deleted_at is null and id=$1`, id).
//		Scan(&behavior.Id,
//			&behavior.Name,
//			&behavior.CreatedAt,
//			&behavior.UpdatedAt)
//	if err != nil {
//		return structs.BehaviorStruct{}, err
//	}
//
//	return behavior, nil
//}
//
//func (b behaviorRepository) GetListBehaviors(page, limit int) ([]structs.BehaviorStruct, int, error) {
//	offset := (page - 1) * limit
//
//	rows, err := b.db.Queryx(`
//		select id, name, created_at, updated_at from behaviors
//		where deleted_at is null limit $1 offset $2`, limit, offset)
//	if err != nil {
//		return nil, 0, nil
//	}
//
//	var (
//		count     int
//		behaviors []structs.BehaviorStruct
//	)
//
//	for rows.Next() {
//		var behavior structs.BehaviorStruct
//		err = rows.Scan(&behavior.Id, &behavior.Name, &behavior.CreatedAt, &behavior.UpdatedAt)
//		if err != nil {
//			return nil, 0, nil
//		}
//
//		behaviors = append(behaviors, behavior)
//	}
//
//	err = b.db.QueryRow(`
//		select count(*) from behaviors where deleted_at is null
//	`).Scan(&count)
//	if err != nil {
//		return nil, 0, err
//	}
//
//	return behaviors, count, nil
//}
//
//func (b behaviorRepository) UpdateBehavior(behavior structs.BehaviorStruct) (structs.BehaviorStruct, error) {
//	result, err := b.db.Exec(`UPDATE behaviors SET name=$1, updated_at=$2 WHERE id=$3`,
//		&behavior.Name,
//		time.Now().UTC(),
//		&behavior.Id)
//	if err != nil {
//		return structs.BehaviorStruct{}, err
//	}
//
//	if i, _ := result.RowsAffected(); i == 0 {
//		return structs.BehaviorStruct{}, sql.ErrNoRows
//	}
//
//	behavior, err = b.GetBehavior(behavior.Id)
//	if err != nil {
//		return structs.BehaviorStruct{}, err
//	}
//
//	return behavior, err
//}
//
//func (b behaviorRepository) DeleteBehavior(id string) error {
//	result, err := b.db.Exec(`UPDATE behaviors SET deleted_at = $1 WHERE id = $2`, time.Now().UTC(), id)
//	if err != nil {
//		return err
//	}
//
//	if i, _ := result.RowsAffected(); i == 0 {
//		return sql.ErrNoRows
//	}
//
//	return nil
//}
