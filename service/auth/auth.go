package sauth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/jobs"
	"github.com/kehl-gopher/logi/internal/mailer"
	"github.com/kehl-gopher/logi/internal/models"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/kehl-gopher/logi/pkg/repository/rabbitmq"
	"github.com/kehl-gopher/logi/pkg/repository/rdb"
	semail "github.com/kehl-gopher/logi/service/email"
)

type Auth struct {
	Db   pdb.Database
	Rdb  rdb.RedisDB
	Conf *config.Config
	Log  *utils.Log
	RM   *rabbitmq.RabbitMQ
}

func (a *Auth) CreateUser(email string, password string) (int, utils.Response) {
	auth := models.User{
		Email:    email,
		Password: password,
	}

	err := auth.CreateUser(a.Db, a.Rdb, a.Conf, a.Log)
	if err != nil {
		if errors.Is(err, utils.ErrorEmailAlreadyExists) {
			message := "bad error response"
			return http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, message, err.Error())
		}
		return http.StatusInternalServerError, utils.ErrorResponse(500, "", err)
	}

	token, err := utils.Tokens()
	if err != nil {
		utils.PrintLog(a.Log, fmt.Sprintf("unable to generate token: %v", err), utils.ErrorLevel)
		return http.StatusInternalServerError, utils.ErrorResponse(500, "", err)
	}

	at := models.AuthToken{Token: token, AddedAt: time.Now().Add(time.Minute * 20), UserID: auth.Id}
	err = at.CreateToken(a.Db, a.Log)
	if err != nil {
		utils.PrintLog(a.Log, fmt.Sprintf("unable to create user token: %v", err), utils.ErrorLevel)
		return http.StatusInternalServerError, utils.ErrorResponse(500, "", err)
	}

	data := map[string]interface{}{
		"code": at.Token}
	e := mailer.NewEmailJob(mailer.VerificationEmail, auth.Email, data)
	body, err := utils.MarshalJSON(e)
	if err != nil {
		return http.StatusInternalServerError, utils.ErrorResponse(500, "", err)
	}

	err = semail.PublishToEmailQUeue(a.RM, jobs.EMAIL_QUEUE, "email.verify", "email_exchange", body, a.Log, &a.Conf.APP_CONFIG)
	if err != nil {
		utils.PrintLog(a.Log, fmt.Sprintf("unable to publish to email queue: %v", err), utils.ErrorLevel)
		return http.StatusInternalServerError, utils.ErrorResponse(500, "", err)
	}

	auth.IsActive = true

	err = auth.UpdateUser(a.Db, a.Log, "is_active")

	if err != nil {
		return http.StatusInternalServerError, utils.ErrorResponse(500, "", err)
	}

	urp := models.UserResponse{
		Email:        auth.Email,
		AuthProvider: auth.AuthProvider,
		Profile: models.ProfileResp{
			FirstName:  "",
			LastName:   "",
			ProfileUrl: "",
			FullName:   "",
		},
		IsVerified: auth.IsVerified,
		CreatedAt:  auth.CreatedAt.Unix(),
		Token:      auth.Token,
	}
	return http.StatusCreated, utils.SuccessfulResponse(http.StatusCreated, "user created successfully", urp)
}

func (a *Auth) UserLogIn(email string, password string) (int, utils.Response) {
	auth := models.User{
		Email:    email,
		Password: password,
	}
	err := auth.GetUserByEmailSignIn(a.Db, a.Rdb, a.Conf, a.Log)

	if err != nil {
		if errors.Is(err, utils.ErrPasswordNotMatch) || errors.Is(err, utils.ErrorNotFound) {
			return http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "invalid email or password provided", "")
		}
		return http.StatusInternalServerError, utils.ErrorResponse(500, "", err)
	}
	profile := models.Profile{Id: auth.Id}
	pr, err := profile.GetUserProfile(a.Rdb, a.Log)
	if err != nil && err.Error() != "cannot fetch data from redis" {
		return http.StatusInternalServerError, utils.ErrorResponse(500, "", err)
	}
	urp := UserResponse(pr, auth.Token, auth)
	return http.StatusOK, utils.SuccessfulResponse(http.StatusOK, "login successful", urp)
}

func UserResponse(p interface{}, at models.AccessToken, usr models.User) models.UserResponse {

	if p == nil {
		return models.UserResponse{
			Id:    usr.Id,
			Email: usr.Email,
			Profile: models.ProfileResp{
				FirstName:  "",
				LastName:   "",
				FullName:   "",
				ProfileUrl: "",
			},
			AuthProvider: "",
			IsVerified:   usr.IsVerified,
			CreatedAt:    usr.CreatedAt.Unix(),
			Token:        at,
		}
	}
	pr, ok := p.(*models.Profile)
	if !ok {
		return models.UserResponse{}
	}
	fmt.Println(pr.CreatedAt)
	return models.UserResponse{
		Id:    pr.Id,
		Email: pr.Email,
		Profile: models.ProfileResp{
			FirstName:  pr.FirstName,
			LastName:   pr.LastName,
			FullName:   pr.FullName,
			ProfileUrl: pr.ProfilePicUrl,
		},
		AuthProvider: pr.AuthProvider,
		IsVerified:   pr.IsVerified,
		CreatedAt:    pr.CreatedAt.Unix(),
		Token:        at,
	}
}
