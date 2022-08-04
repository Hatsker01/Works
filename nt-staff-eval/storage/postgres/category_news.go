package postgres

import (
	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/repo"
	"github.com/jmoiron/sqlx"
)

type newsCategoryRepository struct {
	db *sqlx.DB
}

func NewNewsCategoryRepo(db *sqlx.DB) repo.NewsCategoryRepoInterface {
	return &newsCategoryRepository{
		db: db,
	}
}

func (n newsCategoryRepository) CreateNewsCategory(newsCategory structs.Category) (structs.Category, error) {
	var id string
	err := n.db.QueryRow(`INSERT INTO categories(id, name)
    		VALUES ($1, $2) returning id`, newsCategory.Id, newsCategory.Name).Scan(&id)
	if err != nil {
		return structs.Category{}, err
	}

	return n.GetNewsCategory(id)
}

func (n newsCategoryRepository) GetNewsCategory(id string) (structs.Category, error) {
	var newsCategory structs.Category
	err := n.db.QueryRow(`select id, name from categories
		where deleted_at is null and id=$1`, id).
		Scan(&newsCategory.Id,
			&newsCategory.Name)
	if err != nil {
		return structs.Category{}, err
	}

	return newsCategory, nil
}

func (n newsCategoryRepository) GetListNewsCategory() ([]structs.Category, int, error) {
	query := `SELECT id, name FROM categories WHERE deleted_at IS NULL`

	rows, err := n.db.Queryx(query)
	if err != nil {
		return nil, 0, err
	}

	var newsCategories []structs.Category
	for rows.Next() {
		var newsCategory structs.Category
		err := rows.Scan(
			&newsCategory.Id,
			&newsCategory.Name)
		if err != nil {
			return nil, 0, err
		}

		newsCategories = append(newsCategories, newsCategory)
	}

	return newsCategories, len(newsCategories), nil
}

func (n newsCategoryRepository) UpdateNewsCategory(newsCategory structs.Category) (structs.Category, error) {
	err := n.db.QueryRow(`UPDATE categories SET name=$1 WHERE id=$2 returning id`,
		newsCategory.Name, newsCategory.Id).Scan(&newsCategory.Id)
	if err != nil {
		return structs.Category{}, err
	}

	return n.GetNewsCategory(newsCategory.Id)
}

func (n newsCategoryRepository) DeleteNewsCategory(id string) error {
	_, err := n.db.Exec(`UPDATE categories SET deleted_at=now() WHERE id=$1`, id)
	if err != nil {
		return err
	}

	return nil
}
