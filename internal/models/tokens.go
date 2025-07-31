package models

import (
	"context"
	"fmt"
	"time"

	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/uptrace/bun"
)

type AuthToken struct {
	bun.BaseModel `bun:"table:auth_token"`
	TokenID       string    `bun:"token_id,pk"`
	Token         string    `bun:"tokens"`
	UserID        string    `bun:"user_id"`
	AddedAt       time.Time `bun:"added_at"`
}

func (at *AuthToken) CreateToken(pdb pdb.Database, log *utils.Log) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	at.TokenID = utils.GenerateUUID()

	err := pdb.Insert(ctx, at)

	if err != nil {
		utils.PrintLog(log, fmt.Sprintf("could not add token table data: %v", err), utils.ErrorLevel)
		return err
	}

	return nil
}

func (at *AuthToken) GetToken(pdb pdb.Database, log *utils.Log) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := ` tokens = ? AND user_id = ? AND added_at >= NOW()`
	err := pdb.SelectSingle(ctx, at, query, at.Token, at.UserID)

	if err != nil {
		return err
	}

	return nil
}
