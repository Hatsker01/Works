package postgres

import (
	"time"

	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/repo"
	"github.com/jmoiron/sqlx"
)

type newsRepository struct {
	db *sqlx.DB
}

func NewNewsRepo(db *sqlx.DB) repo.NewsRepoInterface {
	return &newsRepository{
		db: db,
	}
}

func (n newsRepository) CreateNews(news structs.CreateNews) (structs.News, error) {
	var id string
	err := n.db.QueryRow(`INSERT INTO news(id, title, body, cover, author, read_time, category_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7) returning id`,
		news.Id, news.Title, news.Body, news.Cover, news.AuthorId, news.ReadTime, news.CategoryId).Scan(&id)
	if err != nil {
		return structs.News{}, err
	}

	return n.GetNews(id)
}

func (n newsRepository) GetNews(id string) (structs.News, error) {
	var news structs.News
	err := n.db.QueryRow(`select id, title, body, cover, author, read_time, category_id, created_at, updated_at from news
		where deleted_at is null and id=$1`, id).
		Scan(&news.Id,
			&news.Title,
			&news.Body,
			&news.Cover,
			&news.Author.Id,
			&news.ReadTime,
			&news.Category.Id,
			&news.CreatedAt,
			&news.UpdatedAt)
	if err != nil {
		return structs.News{}, err
	}

	return news, nil
}

func (n newsRepository) GetListNews(filters map[string]string, page, limit int) ([]structs.News, int, error) {
	offset := (page - 1) * limit
	var filter string

	query := "SELECT id, title, body, cover, author, read_time, category_id, created_at, updated_at FROM news WHERE deleted_at IS NULL "
	if filters["title"] != "" { // ----------------------------------
		filter += `AND title ilike '%` + filters["title"] + `%' `
	}
	if filters["category_id"] != "" {
		filter += `AND category_id='` + filters["category_id"] + `' `
	}
	if filters["author_id"] != "" {
		filter += "AND author='" + filters["author_id"] + "' "
	}

	rows, err := n.db.Queryx(query+filter+"LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return []structs.News{}, 0, err
	}
	defer rows.Close()

	var newsAll []structs.News
	for rows.Next() {
		var news structs.News
		err = rows.Scan(&news.Id,
			&news.Title,
			&news.Body,
			&news.Cover,
			&news.Author.Id,
			&news.ReadTime,
			&news.Category.Id,
			&news.CreatedAt,
			&news.UpdatedAt)
		if err != nil {
			return []structs.News{}, 0, err
		}
		newsAll = append(newsAll, news)
	}

	var total int

	err = n.db.QueryRow(`select count(*) from news where deleted_at is null ` + filter).Scan(&total)
	if err != nil {
		return []structs.News{}, 0, err
	}

	return newsAll, total, nil
}

func (n newsRepository) UpdateNews(news structs.CreateNews) (string, error) {
	err := n.db.QueryRow(`UPDATE news SET title=$1, body=$2, cover=$3, author=$4, read_time=$5, category_id=$6, updated_at=$7
		WHERE id=$8 RETURNING id`,
		news.Title, news.Body, news.Cover, news.AuthorId, news.ReadTime, news.CategoryId, time.Now(), news.Id).
		Scan(&news.Id)
	if err != nil {
		return "", err
	}

	return news.Id, nil
}

func (n newsRepository) DeleteNews(id string) error {
	_, err := n.db.Exec(`UPDATE news SET deleted_at=$1 WHERE id=$2`, time.Now(), id)
	return err
}
