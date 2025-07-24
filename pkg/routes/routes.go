package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/kehl-gopher/logi/pkg/repository/rdb"
)

func Setup(log *utils.Log, conf *config.Config, db pdb.Database, rdb rdb.RedisDB) *gin.Engine {
	r := gin.New()
	// gin.SetMode(gin.ReleaseMode)
	r.Use(gin.Recovery())

	health(r, log, conf, db, rdb)
	return r
}
