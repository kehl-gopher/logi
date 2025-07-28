package sauth

import (
	"errors"
	"net/http"

	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/models"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/kehl-gopher/logi/pkg/repository/rdb"
)

type Auth struct {
	Db   pdb.Database
	Rdb  rdb.RedisDB
	Conf *config.Config
	Log  *utils.Log
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
