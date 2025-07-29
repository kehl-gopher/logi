package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/kehl-gopher/logi/pkg/repository/rabbitmq"
	"github.com/kehl-gopher/logi/pkg/repository/rdb"
	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(log *utils.Log, conf *config.Config, db pdb.Database, rdb rdb.RedisDB, rm *rabbitmq.RabbitMQ) *gin.Engine {
	r := gin.New()
	// gin.SetMode(gin.ReleaseMode)
	r.Use(gin.Recovery())

	health(r, log, conf, db, rdb)
	authRoutes(r, log, conf, db, rdb, rm)

	r.StaticFile("/swagger.yml", "./static/swagger.yml")
	url := ginSwagger.URL("/swagger.yml")
	r.GET("/api/docs/*any", func(c *gin.Context) {
		ginSwagger.WrapHandler(swaggerFiles.Handler, url)(c)
	})

	return r
}
