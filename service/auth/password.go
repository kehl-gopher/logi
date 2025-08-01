package sauth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kehl-gopher/logi/internal/jobs"
	"github.com/kehl-gopher/logi/internal/mailer"
	"github.com/kehl-gopher/logi/internal/models"
	"github.com/kehl-gopher/logi/internal/utils"
	semail "github.com/kehl-gopher/logi/service/email"
)

func (a *Auth) ForgotPassword(email string) (int, utils.Response) {
	fu := models.Auth{Email: email}

	err := fu.GetUserByEmail(a.Db, a.Log)
	if err != nil {
		if errors.Is(err, utils.ErrorNotFound) {
			return http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "user does not exsits", nil)
		}
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}

	lr, err := utils.NanoId()
	if err != nil {
		fmt.Println(err.Error())
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}
	fmt.Println(fu.Id)
	re := models.ResetPasswordLink{
		UserID:    fu.Id,
		ExpiresAt: time.Now().Add(time.Minute * 30),
		Token:     lr,
		Link:      fmt.Sprintf("%s/%s/change-password", fu.Id, lr),
	}
	err = re.CreateResetPasswordLink(a.Db, a.Log, a.Conf)
	if err != nil {
		fmt.Println(err.Error())
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}

	data := map[string]interface{}{
		"reset_link": fmt.Sprintf("%s/change-password", a.Conf.APP_CONFIG.FRONTEND_URL),
	}
	e := mailer.NewEmailJob(mailer.ForgotPasswordEmail, fu.Email, data)
	body, err := utils.MarshalJSON(e)

	if err != nil {
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}
	err = semail.PublishToEmailQUeue(a.RM, jobs.EMAIL_QUEUE, "email.forgot_password", "email_exchange", body, a.Log, &a.Conf.APP_CONFIG)

	if err != nil {
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}

	data = map[string]interface{}{
		"reset_link": fmt.Sprintf("%s/auth/%s", a.Conf.APP_CONFIG.APP_URL, re.Link),
		"user_id":    fu.Id,
	}

	return http.StatusOK, utils.SuccessfulResponse(http.StatusOK, "a password reset link has been sent to your email account", data)
}

func (a *Auth) ChangePassword(userId string, password string, token string) (int, utils.Response) {
	fmt.Println(userId, password, token)
	lr := models.ResetPasswordLink{UserID: userId, Token: token}
	err := lr.GetPasswordResetLink(a.Db, a.Log, a.Conf)

	if err != nil {
		return http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "password reset link expired", nil)
	}

	password, err = utils.HashPassword(password)
	if err != nil {
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}

	au := models.Auth{
		Id:       userId,
		Password: password,
	}

	err = au.UpdateUser(a.Db, a.Log, "password")
	if err != nil {
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}
	return http.StatusOK, utils.SuccessfulResponse(http.StatusOK, "user password updated successfully", nil)
}
