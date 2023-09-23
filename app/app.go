package app

import (
	"WB_Tech_level_0/infrastructure"
	"WB_Tech_level_0/repository"
	"database/sql"
	"fmt"

	"github.com/nats-io/nats.go"
)

type App struct {
	OrderCache *infrastructure.OrderCache
	DB         *sql.DB
	NC         *nats.Conn
}

func (a *App) Init() {
	a.DB = infrastructure.InitDB()
	a.NC = infrastructure.InitNats()
	a.OrderCache = infrastructure.InitCache(a.DB)
	err := repository.LoadOrdersIntoCache(a.DB, a.OrderCache)
	if err != nil {
		fmt.Println(err)
	}
}
