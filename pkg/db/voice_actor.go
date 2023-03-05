package db

import (
	"context"
	"database/sql"
)

type VoiceActorRepo interface {
	GetByID(ctx context.Context, voiceActorID int) (voiceActorRaw string, err error)
	GetAll(ctx context.Context) (voiceActorsRaw []string, err error)
	Insert(subjectID int, voiceActor string) error
}

type voiceActorRepo struct {
	db *sql.DB
}

func (r *voiceActorRepo) GetByID(ctx context.Context, voiceActorID int) (voiceActorRaw string, err error) {
	row := r.db.QueryRowContext(ctx, "SELECT voice_actor FROM voice_actors WHERE id = ?", voiceActorID)
	err = row.Scan(&voiceActorRaw)
	return
}

func (r *voiceActorRepo) GetAll(ctx context.Context) (voiceActorsRaw []string, err error) {
	rows, err := r.db.QueryContext(ctx, "SELECT voice_actor FROM voice_actors")
	if err != nil {
		return
	}
	err = rows.Scan(&voiceActorsRaw)
	return
}

func (r *voiceActorRepo) Insert(voiceActorID int, voiceActor string) (err error) {
	q := `insert into voice_actors (id, voice_actor)
	values ($1, $2)`
	_, err = r.db.Exec(q, voiceActorID, voiceActor)
	return
}

func NewVoiceActorRepo(db *sql.DB) VoiceActorRepo {
	return &voiceActorRepo{
		db: db,
	}
}
