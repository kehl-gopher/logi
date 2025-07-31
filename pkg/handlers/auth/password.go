package auth

import "github.com/gin-gonic/gin"

func (a *AuthHandler) ForgotPassword(c *gin.Context) {
	var info struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {

	}
}
