package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kehl-gopher/logi/internal/utils"
	sauth "github.com/kehl-gopher/logi/service/auth"
)

func (a *AuthHandler) VerifyToken(c *gin.Context) {
	var token struct {
		Token  string `json:"token" binding:"required"`
		UserId string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBind(&token); err != nil {
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

	tt := sauth.Auth{
		Db:   a.Pdb,
		Conf: a.Conf,
		Log:  a.Log,
		RM:   a.RM,
	}
	statusCode, resp := tt.VerifyUserToken(token.Token, token.UserId)
	c.JSON(statusCode, resp)
}
