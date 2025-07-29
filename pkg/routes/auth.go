package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/handlers/auth"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/kehl-gopher/logi/pkg/repository/rabbitmq"
	"github.com/kehl-gopher/logi/pkg/repository/rdb"
)

func authRoutes(r *gin.Engine, log *utils.Log, conf *config.Config, db pdb.Database, rd rdb.RedisDB, rm *rabbitmq.RabbitMQ) {
	auth := auth.AuthHandler{Log: log, Conf: conf, Pdb: db, Rdb: rd, RM: rm}
	api := r.Group("/api/v1/auth")
	{
		api.POST("/signup", auth.CreateUser)
		api.POST("/signin", auth.SignInUser)
	}
}
