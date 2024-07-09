package routers

import (
	"beego-project/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/authenticate", &controllers.TableauController{}, "post:PostAuth")
	beego.Router("/datasources", &controllers.TableauControllerDS{}, "get:GetDataSources")
	beego.Router("/datalabels", &controllers.TableauControllerDL{}, "get:GetDataLabels")
}
