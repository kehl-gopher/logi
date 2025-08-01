package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kehl-gopher/logi/internal/utils"
	sauth "github.com/kehl-gopher/logi/service/auth"
)

func (a *AuthHandler) ForgotPassword(c *gin.Context) {
	var info struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
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
		Conf: a.Conf,
		Db:   a.Pdb,
		Log:  a.Log,
		RM:   a.RM,
		Rdb:  a.Rdb,
	}

	code, resp := au.ForgotPassword(info.Email)
	c.JSON(code, resp)
}

func (a *AuthHandler) ChangePassword(c *gin.Context) {
	var user struct {
		Password string `json:"password" binding:"required"`
	}
	user_id := c.Param("userId")
	token := c.Param("token")

	if err := c.ShouldBindJSON(&user); err != nil {
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
		RM:   a.RM,
	}
	resp, statusCode := au.ChangePassword(user_id, user.Password, token)
	c.JSON(resp, statusCode)
}
