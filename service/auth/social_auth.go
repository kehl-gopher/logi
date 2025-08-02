package sauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kehl-gopher/logi/internal/jobs"
	"github.com/kehl-gopher/logi/internal/mailer"
	"github.com/kehl-gopher/logi/internal/models"
	"github.com/kehl-gopher/logi/internal/utils"
	semail "github.com/kehl-gopher/logi/service/email"
	"github.com/markbates/goth"
	"google.golang.org/api/idtoken"
)

func (a *Auth) CreateGoogleUser(user goth.User) (int, utils.Response) {
	var provider string
	payload, err := idtoken.Validate(context.Background(), user.IDToken, "")

	if err != nil {
		return http.StatusBadRequest,
			utils.ErrorResponse(http.StatusBadRequest, fmt.Sprintf("invalid token provided from: %s", payload.Audience), nil)
	}

	if user.Provider == "" {
		provider = "google"
	} else {
		provider = strings.ToLower(user.Provider)
	}

	appConf := a.Conf.APP_CONFIG
	password, err := utils.HashPassword(appConf.OAUTH_PASSWORD)
	if err != nil {
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}
	var usr = models.User{
		Email: user.Email,
	}

	err = usr.GetUserByEmail(a.Db, a.Log)
	if err != nil && !errors.Is(err, utils.ErrorNotFound) {
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	} else if err == nil && usr.AuthProvider == provider {
		fmt.Println("-------------------------------_>")
		profile := models.Profile{Id: usr.Id}
		pr, err := profile.GetUserProfile(a.Rdb, a.Log)
		if err != nil && err.Error() != "cannot fetch data from redis" {
			return http.StatusInternalServerError, utils.ErrorResponse(500, "", err)
		}

		ttl, err := utils.ParseTime("2d")
		if err != nil {
			return http.StatusInternalServerError, utils.ErrorResponse(500, "", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		ac := a.Conf.APP_CONFIG
		token, err := models.InsertAccessToken(ctx, a.Rdb, ac.JWT_SECRETKEY, ttl, usr.Id)
		if err != nil {
			return http.StatusInternalServerError, utils.ErrorResponse(500, "", err)
		}
		urp := UserResponse(pr, *token, usr)

		return http.StatusAccepted, utils.SuccessfulResponse(http.StatusAccepted, "login successful", urp)
	} else if usr.AuthProvider != provider && err == nil {
		return http.StatusConflict, utils.ErrorResponse(http.StatusConflict, "please login with email/password", err)
	}

	usr.Password = password
	usr.IsActive = true
	usr.IsVerified = true
	usr.AuthProvider = provider

	err = usr.CreateUser(a.Db, a.Rdb, a.Conf, a.Log)
	if err != nil {
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}

	p := models.Profile{
		Email:         usr.Email,
		Id:            usr.Id,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		FullName:      user.FirstName + " " + user.LastName,
		ProfilePicUrl: user.AvatarURL,
		IsVerified:    usr.IsVerified,
		CreatedAt:     time.Now(),
		AuthProvider:  usr.AuthProvider,
	}

	err = p.CreateProfile(a.Db, a.Rdb, a.Log, a.Conf)
	if err != nil {
		utils.PrintLog(a.Log, fmt.Sprintf("unable to cache user profile data: %v", err), utils.ErrorLevel)
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}

	e := mailer.NewEmailJob(mailer.WelcomeEmail, usr.Email, nil)
	body, err := utils.MarshalJSON(e)
	if err != nil {
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}
	err = semail.PublishToEmailQUeue(a.RM, jobs.EMAIL_QUEUE, "email.welcome", "email_exchange", body, a.Log, &a.Conf.APP_CONFIG)
	if err != nil {
		return http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "", err)
	}

	urp := models.UserResponse{
		Id:    usr.Id,
		Email: usr.Email,
		Profile: models.ProfileResp{
			ProfileUrl: p.ProfilePicUrl,
			LastName:   p.LastName,
			FirstName:  p.FirstName,
			FullName:   p.FullName,
		},
		IsVerified:   p.IsVerified,
		AuthProvider: p.AuthProvider,
		CreatedAt:    p.CreatedAt.Unix(),
		Token:        p.Token,
	}
	return http.StatusAccepted, utils.SuccessfulResponse(http.StatusAccepted, "user sign in successful", urp)
}
