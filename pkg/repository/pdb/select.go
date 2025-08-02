package pdb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kehl-gopher/logi/internal/utils"
)

func (p *postgresConn) SelectSingle(ctx context.Context, model interface{}, query string, args ...interface{}) error {
	err := p.bun.NewSelect().Model(model).Where(query, args...).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrorNotFound
		}
		return err
	}

	return err
}

func (p *postgresConn) CheckExists(ctx context.Context, model interface{}, query string, args ...interface{}) (bool, error) {
	exists, err := p.bun.NewSelect().Model(model).Where(query, args...).Exists(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}
