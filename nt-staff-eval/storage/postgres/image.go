package postgres

import (
	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/repo"
	"github.com/jmoiron/sqlx"
)

type imageRepository struct {
	db *sqlx.DB
}

func NewImageRepo(db *sqlx.DB) repo.ImageRepoInterface {
	return &imageRepository{
		db: db,
	}
}

func (r imageRepository) LoadImage(imageType int8, path, id string) (err error) {
	var query string
	switch imageType {
	case structs.TypeSectionImage:
		query = structs.LoadSectionImage
	case structs.TypeUserImage:
		query = structs.LoadUserImage
	}
	_, err = r.db.Exec(query, path, id)
	if err != nil {
		return err
	}
	return
}
