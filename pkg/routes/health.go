package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/handlers"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/kehl-gopher/logi/pkg/repository/rdb"
)

func health(r *gin.Engine, log *utils.Log, conf *config.Config, db pdb.Database, rd rdb.RedisDB) {
	h := handlers.Handler{Log: log, Conf: conf, Db: db, Rdb: rd}
	r.GET("/api/v1/healthcheck", h.Health)
}
