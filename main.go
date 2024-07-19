package main

import (
	_ "beego-project/routers"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func main() {
	beego.Run()
}
