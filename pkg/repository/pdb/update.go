package pdb

import "context"

func (p *postgresConn) UpdateModel(ctx context.Context, model interface{}, column string, query string, args ...interface{}) error {
	_, err := p.bun.NewUpdate().Model(model).Column(column).Where(query, args...).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
