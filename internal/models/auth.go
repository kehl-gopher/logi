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

type User struct {
	bun.BaseModel `bun:"table:users"`
	Id            string      `bun:"id" json:"id"`
	Email         string      `bun:"email" json:"email"`
	Password      string      `bun:"password" json:"-"`
	IsActive      bool        `bun:"is_active" json:"-"`
	AuthProvider  string      `bun:"auth_provider,nullzero,default"`
	IsVerified    bool        `bun:"is_verified" json:"is_verified"`
	Deactivated   bool        `bun:"deactivated" json:"-"`
	CreatedAt     time.Time   `bun:"created_at,nullzero,default" json:"created_at"`
	UpdatedAt     time.Time   `bun:"updated_at,nullzero"`
	Token         AccessToken `bun:"-" json:"access_token"`
}

type AccessToken struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type ProfileResp struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ProfileUrl string `json:"profile_url"`
	FullName   string `json:"full_name"`
}
type UserResponse struct {
	Id           string      `json:"id"`
	Email        string      `json:"email"`
	AuthProvider string      `json:"auth_provider"`
	IsVerified   bool        `json:"is_verified"`
	CreatedAt    int64       `json:"created_at"`
	Profile      ProfileResp `json:"profile"`
	Token        AccessToken `json:"access_token"`
}

func (a *User) CreateUser(pdb pdb.Database, rdb rdb.RedisDB, conf *config.Config, log *utils.Log) error {
	a.Id = utils.GenerateUUID()
	password, err := utils.HashPassword(a.Password)
	if err != nil {
		utils.PrintLog(log, err.Error(), utils.ErrorLevel)
		return err
	}
	a.Password = password

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = pdb.Insert(ctx, a)

	if err != nil {
		utils.PrintLog(log, err.Error(), utils.ErrorLevel)
		return err
	}
	ttl, err := utils.ParseTime(conf.APP_CONFIG.JWT_DURATIONTIME)
	if err != nil {
		utils.PrintLog(log, "failed to convert unit time", utils.ErrorLevel)
		return err
	}
	sec := conf.APP_CONFIG.JWT_SECRETKEY
	acc, err := InsertAccessToken(ctx, rdb, sec, ttl, a.Id)

	if err != nil {
		utils.PrintLog(log, fmt.Sprintf("failed to generate access token: %s", err.Error()), utils.ErrorLevel)
		return err
	}
	a.Token = *acc
	fmt.Println(a.Token)
	return nil
}

func (a *User) GetUserByEmailSignIn(pdb pdb.Database, rdb rdb.RedisDB, conf *config.Config, log *utils.Log) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := `email = ?`
	password := a.Password
	err := pdb.SelectSingle(ctx, a, query, a.Email)
	if err != nil {
		return err
	}

	pCheck := utils.CompareHashedPassword(password, a.Password)
	if !pCheck {
		return utils.ErrPasswordNotMatch
	}

	ttl, err := utils.ParseTime(conf.APP_CONFIG.JWT_DURATIONTIME)
	if err != nil {
		return err
	}
	token, err := InsertAccessToken(ctx, rdb, conf.APP_CONFIG.JWT_SECRETKEY, ttl, a.Id)
	if err != nil {
		return err
	}
	a.Token = *token
	return nil
}

func (a *User) GetUserById(pdb pdb.Database, conf *config.Config, log *utils.Log) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := `id = ?`
	err := pdb.SelectSingle(ctx, a, query, a.Id)
	if err != nil {
		return err
	}
	return nil
}

func (a *User) UpdateUser(pdb pdb.Database, log *utils.Log, column string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := ` id = ?`
	err := pdb.UpdateModel(ctx, a, column, query, a.Id)

	if err != nil {
		utils.PrintLog(log, "failed to update user auth table", utils.ErrorLevel)
		return err
	}

	return nil
}

func (a *User) GetUserByEmail(pdb pdb.Database, log *utils.Log) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `email = ?`

	err := pdb.SelectSingle(ctx, a, query, a.Email)

	if err != nil {
		return err
	}
	return nil
}

func (u *User) CheckExists(db pdb.Database, query string, args ...interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	exist, err := db.CheckExists(ctx, u, query, args...)
	return exist, err
}

func InsertAccessToken(ctx context.Context, rdb rdb.RedisDB, secret string, ttl time.Duration, userId string) (*AccessToken, error) {
	key := fmt.Sprintf("user:%s", userId)
	token, exp, err := utils.GenerateJWT(userId, secret, ttl)
	if err != nil {
		return nil, err
	}

	err = rdb.Set(ctx, key, token, ttl)
	if err != nil {
		return nil, err
	}
	acc := &AccessToken{
		Token:     token,
		ExpiresAt: exp.UnixNano(),
	}

	return acc, nil
}
