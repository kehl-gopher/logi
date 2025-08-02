package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/kehl-gopher/logi/pkg/repository/rabbitmq"
	"github.com/kehl-gopher/logi/pkg/repository/rdb"
	sauth "github.com/kehl-gopher/logi/service/auth"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
	Conf *config.Config
	Log  *utils.Log
	Pdb  pdb.Database
	Rdb  rdb.RedisDB
	RM   *rabbitmq.RabbitMQ
}

func (a *AuthHandler) CreateUser(c *gin.Context) {
	var auth struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&auth); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			resp := utils.ValidationErrorResponse(ve)
			c.JSON(http.StatusUnprocessableEntity, resp)
			return
		}
		resp := utils.ErrorResponse(http.StatusBadRequest, "bad error response", err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	au := sauth.Auth{
		Db:   a.Pdb,
		Rdb:  a.Rdb,
		Conf: a.Conf,
		Log:  a.Log,
		RM:   a.RM,
	}

	statusCode, resp := au.CreateUser(auth.Email, auth.Password)

	c.JSON(statusCode, resp)
}

func (a *AuthHandler) SignInUser(c *gin.Context) {
	var auth struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&auth); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			resp := utils.ValidationErrorResponse(ve)
			c.JSON(http.StatusUnprocessableEntity, resp)
			return
		}
		resp := utils.ErrorResponse(http.StatusBadRequest, "bad error response", err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	au := sauth.Auth{Db: a.Pdb, Rdb: a.Rdb, Conf: a.Conf, Log: a.Log}
	statusCode, resp := au.UserLogIn(auth.Email, auth.Password)
	c.JSON(statusCode, resp)
}

func (a *AuthHandler) SignInGoogleAuth(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", "google")

	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (a *AuthHandler) OAuthCallBack(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", "google")

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)

	if err != nil {
		fmt.Println(err.Error())
		resp := utils.ErrorResponse(http.StatusInternalServerError, "", err)
		c.JSON(http.StatusInternalServerError, resp)
	}
	au := sauth.Auth{
		Db:   a.Pdb,
		Conf: a.Conf,
		Log:  a.Log,
		Rdb:  a.Rdb,
		RM:   a.RM,
	}

	statusCode, resp := au.CreateGoogleUser(user)
	c.JSON(statusCode, resp)
}
