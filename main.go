package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	lg "log"

	"github.com/joho/godotenv"
	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/kehl-gopher/logi/pkg/repository/pdb"
	"github.com/kehl-gopher/logi/pkg/repository/rabbitmq"
	"github.com/kehl-gopher/logi/pkg/repository/rdb"
	"github.com/kehl-gopher/logi/pkg/routes"
)

func main() {
	var drop bool
	flag.BoolVar(&drop, "drop", false, "force drop tables on dirty migrations") // not to be used in production
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

	if drop && strings.ToLower(conf.APP_CONFIG.APP_ENV) == "dev" {
		utils.Droptables(db.DB())
		utils.PrintLog(log, "all tables are dropped.", utils.DebugLevel)
		return
	}
	// rabbitmq connection
	rq := rabbitmq.NewMQManager(&conf.RabbitMQ, log)
	defer rq.Close()

	r := routes.Setup(log, conf, db, red, rq)
	lg.Fatal(r.Run(fmt.Sprintf(":%d", 8080)))
}
