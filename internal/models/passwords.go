package models

import (
	"context"
	"fmt"
	"time"

	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/uptrace/bun"
)

type ResetPasswordLink struct {
	bun.BaseModel `bun:"table:password_reset_link"`
	Id            string    `bun:"id,pk"`
	Link          string    `bun:"link"`
	Token         string    `bun:"token"`
	UserID        string    `bun:"user_id"`
	ExpiresAt     time.Time `bun:"expires_at"`
	AddedAt       time.Time `bun:"added_at,nullzero,default"`
}

func (r *ResetPasswordLink) CreateResetPasswordLink(pdb pdb.Database, lg *utils.Log, conf *config.Config) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r.Id = utils.GenerateUUID()
	err := pdb.Insert(ctx, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *ResetPasswordLink) GetPasswordResetLink(pdb pdb.Database, lg *utils.Log, conf *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	query := `user_id = ? AND expires_at >= NOW() AND token = ?`

	fmt.Println(r.Token, r.UserID)
	err := pdb.SelectSingle(ctx, r, query, r.UserID, r.Token)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
