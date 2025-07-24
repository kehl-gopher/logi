package main

import (
	"context"
	"fmt"
	"time"

	lg "log"

	"github.com/joho/godotenv"
	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/kehl-gopher/logi/pkg/repository/rdb"
	"github.com/kehl-gopher/logi/pkg/routes"
)

func main() {
	log := utils.NewLogger()
	err := godotenv.Load()
	if err != nil {
		utils.PrintLog(log, err.Error(), utils.FatalLevel)
		return
	}
	conf := config.LoadConfig(log)

	// postgres connection
	db := pdb.NewPostgresConn(conf, log)
	err = db.ConnectPostgres()
	if err != nil {
		utils.PrintLog(log, fmt.Sprintf("failed to connect to database: %s", err.Error()), utils.ErrorLevel)
		return
	}
	defer db.Close()

	// redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	red := rdb.NewRedisConn(log, conf)
	err = red.ConnectRedis(ctx)
	if err != nil {
		utils.PrintLog(log, fmt.Sprintf("failed to connect to redis %s", err.Error()), utils.ErrorLevel)
		return
	}
	defer red.Close()

	r := routes.Setup(log, conf, db, red)
	lg.Fatal(r.Run(fmt.Sprintf(":%d", 8080)))
}
