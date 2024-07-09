package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"fmt"
)

type TableauControllerDL struct {
	beego.Controller
}

func (c *TableauControllerDL) GetDataLabels() {
	fmt.Println("	in GetDataLabels . . . ")

}
