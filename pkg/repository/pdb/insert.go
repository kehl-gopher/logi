package pdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kehl-gopher/logi/internal/utils"
)

func (p *postgresConn) Insert(ctx context.Context, model interface{}) error {
	err := p.bun.NewInsert().
		Model(model).
		Returning("*").
		Scan(ctx)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return utils.ErrorEmailAlreadyExists
		}
		return fmt.Errorf("unexpected db error: %w", err)
	}
	return nil
}

func (p *postgresConn) InsertMany(ctx context.Context, models ...interface{}) error {
	tx, err := p.bun.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.NewInsert().Model(&models).Exec(ctx); err != nil {
		if rerr := tx.Rollback(); rerr != nil && rerr != sql.ErrTxDone {
			return fmt.Errorf("insert error: %v, rollback error: %v", err, rerr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		if rerr := tx.Rollback(); rerr != nil && rerr != sql.ErrTxDone {
			return fmt.Errorf("commit error: %v, rollback error: %v", err, rerr)
		}
		return err
	}

	return nil
}
