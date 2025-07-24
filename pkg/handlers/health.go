package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/kehl-gopher/logi/pkg/repository/rdb"
)

type HealthResponse struct {
	Healthy   bool            `json:"healthy"`
	Version   string          `json:"version"`
	Uptime    string          `json:"uptime"`
	Timestamp string          `json:"timestamp"`
	Checks    map[string]bool `json:"checks"`
}

var startTime = time.Now()

type Handler struct {
	Log  *utils.Log
	Conf *config.Config
	Db   pdb.Database
	Rdb  rdb.RedisDB
}

func (h *Handler) Health(c *gin.Context) {
	uptime := time.Since(startTime)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dHealth := h.Db.PingDB()
	rHealth := h.Rdb.PingCache(ctx)

	healthy := dHealth && rHealth

	hr := HealthResponse{
		Healthy:   healthy,
		Version:   h.Conf.APP_CONFIG.APP_VERSION,
		Uptime:    uptime.String(),
		Timestamp: time.Now().String(),
		Checks: map[string]bool{
			"postgres-healthy": dHealth,
			"redis-healthy":    rHealth,
		},
	}
	c.JSON(http.StatusOK, hr)
}
