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
	auth := models.Auth{
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

	fmt.Println(token)
	at := models.AuthToken{Token: token, AddedAt: time.Now().Add(time.Minute * 20), UserID: auth.Id}
	err = at.CreateToken(a.Db, a.Log)
	if err != nil {
		utils.PrintLog(a.Log, fmt.Sprintf("unable to create user token: %v", err), utils.ErrorLevel)
		return http.StatusInternalServerError, utils.ErrorResponse(500, "", err)
	}

	e := mailer.EmailJOB{
		To:   auth.Email,
		Type: mailer.VerificationEmail,
		Data: map[string]interface{}{
			"code": at.Token,
		},
	}
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

	return http.StatusCreated, utils.SuccessfulResponse(http.StatusCreated, "user created successfully", auth)
}

func (a *Auth) UserLogIn(email string, password string) (int, utils.Response) {
	auth := models.Auth{
		Email:    email,
		Password: password,
	}

	err := auth.GetUser(a.Db, a.Rdb, a.Conf, a.Log)

	if err != nil {
		if errors.Is(err, utils.ErrPasswordNotMatch) || errors.Is(err, utils.ErrorNotFound) {
			return http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "invalid email or password provided", "")
		}
		return http.StatusInternalServerError, utils.ErrorResponse(500, "", err)
	}
	return http.StatusOK, utils.SuccessfulResponse(http.StatusOK, "login successful", auth)
}
