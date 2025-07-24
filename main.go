package main

import (
	"github.com/joho/godotenv"
	"github.com/kehl-gopher/logi/internal/utils"
)

func main() {
	log := utils.NewLogger()
	err := godotenv.Load()
	if err != nil {
		utils.PrintLog(log, err.Error(), utils.FatalLevel)
		return
	}


	
}
