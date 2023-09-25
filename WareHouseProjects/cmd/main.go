package main

import (
	"WareHouseProjects/api"
	"WareHouseProjects/api/handler"
	"WareHouseProjects/config"
	"WareHouseProjects/pkg/logger"
	postgres "WareHouseProjects/storage/postgress"
	"context"
	"fmt"
)

func main() {
	fmt.Println("start")
	cfg := config.Load()
	log := logger.NewLogger("warehouse_project", logger.LevelInfo)
	strg, err := postgres.NewStorage(context.Background(), cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	h := handler.NewHandler(strg, log)

	r := api.NewServer(h)
	r.Run(cfg.Port)
}
