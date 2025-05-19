package repository

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/gihyodocker/taskapp/pkg/models"
)

type Task interface {
	Upsert(ctx context.Context, m *models.Task) error
	DeleteByID(ctx context.Context, id string) (int64, error)
	FindByID(ctx context.Context, id string) (*models.Task, error)
	FindAll(ctx context.Context) ([]*models.Task, error)
}

func NewTask(db *sql.DB) Task {
	return &task{
		db: db,
	}
}

type task struct {
	db *sql.DB
}

func (r task) Upsert(ctx context.Context, m *models.Task) error {
	// The Upsert signature for some SQLBoiler drivers (like PostgreSQL)
	// requires updateOnConflict (bool) and conflictColumns ([]string).
	return m.Upsert(
		ctx,
		r.db,
		true,                            // updateOnConflict: true to update if a conflict occurs
		[]string{models.TaskColumns.ID}, // conflictColumns: define conflict on the ID column
		boil.Whitelist("title", "content", "status", "updated"), // updateColumns
		boil.Infer(), // insertColumns
	)
}

func (r task) DeleteByID(ctx context.Context, id string) (int64, error) {
	return models.Tasks(
		qm.Where("id = ?", id),
	).DeleteAll(ctx, r.db)
}

func (r task) FindByID(ctx context.Context, id string) (*models.Task, error) {
	return models.FindTask(ctx, r.db, id)
}

func (r task) FindAll(ctx context.Context) ([]*models.Task, error) {
	return models.Tasks(qm.OrderBy("updated DESC")).All(ctx, r.db)
}
