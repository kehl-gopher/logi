package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/kehl-gopher/logi/pkg/repository/rdb"
	sauth "github.com/kehl-gopher/logi/pkg/service/auth"
)

type AuthHandler struct {
	Conf *config.Config
	Log  *utils.Log
	Pdb  pdb.Database
	Rdb  rdb.RedisDB
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
	}

	statusCode, resp := au.CreateUser(auth.Email, auth.Password)

	if statusCode != http.StatusCreated {
		c.JSON(statusCode, resp)
		return
	}
	c.JSON(statusCode, resp)
}
