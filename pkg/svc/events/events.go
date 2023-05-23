package events

import (
	"app/pkg/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Svc struct {
	pool *pgxpool.Pool
}

func NewSvc(pool *pgxpool.Pool) *Svc {
	return &Svc{pool: pool}
}

func (s *Svc) Register(ctx context.Context, model *models.Event) (err error) {
	tx, err := s.pool.Begin(ctx)
	defer func() {
		if err != nil {
			log.Printf("error: events register %s", err)
			_ = tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	_, err = tx.Exec(
		ctx,
		// language=PostgreSQL
		`INSERT INTO "events"("action", "product", "fingerprint") VALUES ($1, $2, $3)`,
		model.Action, model.Product, model.Fingerprint,
	)

	return
}
