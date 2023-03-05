package db

import (
	"context"
	"database/sql"
)

type SubjectRepo interface {
	GetByID(ctx context.Context, subjectID int) (subject string, err error)
	Insert(subjectID int, subject string) error
}

type subjectRepo struct {
	db *sql.DB
}

func (r *subjectRepo) GetByID(ctx context.Context, subjectID int) (subjectRaw string, err error) {
	row := r.db.QueryRowContext(ctx, "SELECT subject FROM subjects WHERE id = ?", subjectID)
	err = row.Scan(&subjectRaw)
	return
}

func (r *subjectRepo) Insert(subjectID int, subject string) (err error) {
	q := `insert into subjects (id, subject)
	values ($1, $2)`
	_, err = r.db.Exec(q, subjectID, subject)
	return
}

func NewSubjectRepo(db *sql.DB) SubjectRepo {
	return &subjectRepo{
		db: db,
	}
}
