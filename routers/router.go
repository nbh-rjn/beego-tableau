package routers

import (
	"beego-project/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

// fetch site ID from database
func init() {
	beego.Router("/", &controllers.TableauController{})

	beego.Router("/authenticate", &controllers.TableauController{}, "post:PostAuth")
	beego.Router("/sync", &controllers.TableauController{}, "post:PostSync")
	beego.Router("/download", &controllers.TableauController{}, "get:DownloadDataSource")

	// fetching projects, datasources, datalabels, etc . . .
	//
	beego.Router("/attribute/:param", &controllers.TableauController{}, "get:GetAttribute")
}
