package pdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kehl-gopher/logi/internal/utils"
)

func (p *postgresConn) Insert(ctx context.Context, model interface{}) error {
	r, err := p.bun.NewInsert().
		Model(model).
		Exec(ctx)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return utils.ErrorEmailAlreadyExists
		}
		return fmt.Errorf("unexpected db error: %w", err)
	}

	if r, err := r.RowsAffected(); err != nil || r != 1 {
		return utils.ErrorTableInsertFailed
	}
	return nil
}

// bulk insert...
func (p *postgresConn) InsertMany(ctx context.Context, models ...interface{}) error {
	tx, err := p.bun.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.NewInsert().Model(&models).Exec(ctx)
	if err != nil {
		return err // rollback happens here
	}
	return tx.Commit()
}
