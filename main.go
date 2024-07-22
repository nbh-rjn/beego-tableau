package main

import (
	_ "beego-project/routers"
	"fmt"

	"github.com/beego/beego/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func init() {
	orm.RegisterDataBase("default", "postgres", "postgres://postgres:postgres@localhost/bg-db?sslmode=disable")
	orm.RegisterDriver("postgresql", orm.DRPostgres)
}

func main() {
	orm.RunSyncdb("default", false, false)
	db, err := orm.GetDB()
	if err != nil {
		fmt.Println("get default DataBase")
	}
	orm.AddAliasWthDB("default", "postgres", db)

	beego.Run()
}
