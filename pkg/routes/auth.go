package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/handlers/auth"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/kehl-gopher/logi/pkg/repository/rabbitmq"
	"github.com/kehl-gopher/logi/pkg/repository/rdb"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func authRoutes(r *gin.Engine, log *utils.Log, conf *config.Config, db pdb.Database, rd rdb.RedisDB, rm *rabbitmq.RabbitMQ) {

	// instantiate google auth credential

	c := conf.APP_CONFIG
	store := sessions.NewCookieStore([]byte(c.GOOGLE_API_KEY))
	gothic.Store = store
	goth.UseProviders(google.New(c.GOOGLE_CLIENT_ID, c.GOOGLE_CLIENT_SECRET, c.GOOGLE_CALLBACK, "openid", "email", "profile"))
	auth := auth.AuthHandler{Log: log, Conf: conf, Pdb: db, Rdb: rd, RM: rm}
	api := r.Group("/api/v1/auth")
	{
		api.POST("/signup", auth.CreateUser)
		api.POST("/signin", auth.SignInUser)
		api.POST("/verify-token", auth.VerifyToken)
		api.POST("/forgot-password", auth.ForgotPassword)
		api.POST("/:userId/:token/change-password", auth.ChangePassword)
		api.GET("/google", auth.SignInGoogleAuth)
		api.GET("/callback", auth.OAuthCallBack)
	}
}
