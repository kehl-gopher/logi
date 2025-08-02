package models

import (
	"context"
	"fmt"
	"time"

	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/kehl-gopher/logi/pkg/repository/rdb"
	"github.com/uptrace/bun"
)

type Profile struct {
	bun.BaseModel `bun:"table:profiles"`
	Id            string      `bun:"id,pk" json:"id"`
	Email         string      `bun:"-" json:"email"`
	FirstName     string      `bun:"first_name" json:"first_name"`
	LastName      string      `bun:"last_name" json:"last_name"`
	FullName      string      `bun:"full_name" json:"full_name"`
	IsVerified    bool        `bun:"-" json:"is_verified"`
	AuthProvider  string      `bun:"-" json:"auth_provider"`
	ProfilePicUrl string      `bun:"profile_pic" json:"profile_pic"`
	CreatedAt     time.Time   `bun:"created_at,nullzero,default" json:"created_at"`
	UpdatedAt     time.Time   `bun:"updated_at,nullzero" json:"updated_at"`
	Token         AccessToken `bun:"token"`
}

func (p *Profile) CreateProfile(db pdb.Database, rdb rdb.RedisDB, lg *utils.Log, conf *config.Config) error {
	p.FullName = fmt.Sprintf("%s %s", p.FirstName, p.LastName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	key := "user_profile:" + p.Id
	pbyte, err := utils.MarshalJSON(p)

	if err != nil {
		return err
	}
	err = rdb.Queue(ctx, key, pbyte)
	if err != nil {
		return err
	}

	ttl, err := utils.ParseTime("7d")
	if err != nil {
		return err
	}
	err = rdb.Set(ctx, key, pbyte, ttl)

	if err != nil {
		return err
	}

	ttl, err = utils.ParseTime("2d")
	if err != nil {
		return err
	}
	sec := conf.APP_CONFIG.JWT_SECRETKEY
	acc, err := InsertAccessToken(ctx, rdb, sec, ttl, p.Id)
	if err != nil {
		return err
	}
	p.Token = *acc
	return nil
}

func (p *Profile) GetUserProfile(red rdb.RedisDB, log *utils.Log) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	key := fmt.Sprintf("user_profile:%s", p.Id)
	defer cancel()
	receiver, err := red.Get(ctx, key, p)
	if err != nil {
		return nil, err
	}
	return receiver, nil
}
