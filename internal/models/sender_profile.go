package models

import (
	"context"
	"time"

	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/uptrace/bun"
)

type SenderProfile struct {
	bun.BaseModel `bun:"table:sende"`
	ID            string `json:"-" bun:"id,pk"`
	FirstName     string `json:"first_name" bun:"first_name"`
	PhoneNumber   string `json:"phone_number" bun:"phone_number"`
	LastName      string `json:"last_name" bun:"last_name"`
	ProfileImage  string `json:"profile_image" bun:"profile_image,"`
	CreatedAt     string `json:"-" bun:"created_at,nullzero,default"`
	UpdatedAt     string `json:"updated_at" bun:"updated_at"`
}

// perform update on user profile
func (s *SenderProfile) CreateSenderProfile(pdb pdb.Database, log *utils.Log) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := pdb.Insert(ctx, s)

	if err != nil {
		utils.PrintLog(log, "failed to add data to sender table", utils.ErrorLevel)
		return err
	}
	utils.PrintLog(log, "successfully added data to sender table", utils.InfoLevel)
	return nil
}

func (s *SenderProfile) GetSenderProfile(pdb pdb.Database, log *utils.Log) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := `where id = ?`
	err := pdb.SelectSingle(ctx, s, query, s.ID)
	if err != nil {
		utils.PrintLog(log, "failed to select user profile from database successful", utils.ErrorLevel)
		return err
	}

	utils.PrintLog(log, "selected user profile from database successful", utils.InfoLevel)
	return nil
}

func (s *SenderProfile) DeleteSenderProfile() error {
	return nil
}
