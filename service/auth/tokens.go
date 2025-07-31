package sauth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kehl-gopher/logi/internal/jobs"
	"github.com/kehl-gopher/logi/internal/mailer"
	"github.com/kehl-gopher/logi/internal/models"
	"github.com/kehl-gopher/logi/internal/utils"
	semail "github.com/kehl-gopher/logi/service/email"
)

func (a *Auth) VerifyUserToken(token string, userId string) (int, utils.Response) {
	tken := models.AuthToken{Token: token, UserID: userId}

	err := tken.GetToken(a.Db, a.Log)

	if err != nil {
		utils.PrintLog(a.Log, fmt.Sprintf("failed to verify token: %v", err), utils.ErrorLevel)
		return http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid or expired verification code.", nil)
	}

	au := models.Auth{Id: tken.UserID}

	err = au.GetUserById(a.Db, a.Conf, a.Log)
	if err != nil {
		if errors.Is(err, utils.ErrPasswordNotMatch) || errors.Is(err, utils.ErrorNotFound) {
			return http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "user not found", nil)
		}
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}

	e := mailer.NewEmailJob(mailer.WelcomeEmail, au.Email, nil)
	body, err := utils.MarshalJSON(e)
	if err != nil {
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}

	err = semail.PublishToEmailQUeue(a.RM, jobs.EMAIL_QUEUE, "email.welcome", "email_exchange", body, a.Log, &a.Conf.APP_CONFIG)
	if err != nil {
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}

	au.IsVerified = true
	err = au.UpdateUser(a.Db, a.Log, "is_verified")
	if err != nil {
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}

	return http.StatusOK, utils.SuccessfulResponse(http.StatusOK, "user successfully verified", nil)
}
