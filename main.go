package main

import (
	_ "beego-project/routers"
	"fmt"
	"os"

	"github.com/beego/beego/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func init() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		panic("DATABASE_URL environment variable is not set")
	}
	orm.RegisterDriver("postgresql", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", dbURL)
}

func main() {
	orm.RunSyncdb("default", true, false)
	db, err := orm.GetDB()
	if err != nil {
		fmt.Println("get default DataBase")
	}
	orm.AddAliasWthDB("default", "postgres", db)

	beego.Run()
}
