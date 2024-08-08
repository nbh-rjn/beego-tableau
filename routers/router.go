package routers

import (
	"beego-project/controllers"
	"beego-project/logger" // Adjusted import path

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Initialize the logger
	log, _ := logger.NewZapLogger()

	// Create controller instances with logger
	projectController := &controllers.TableauController{Logger: log}

	// Set up routes
	beego.Router("/", projectController)

	beego.Router("/authenticate", projectController, "post:PostAuth")
	beego.Router("/sync", projectController, "post:PostSync")
	beego.Router("/download", projectController, "get:DownloadDataSource")

	// fetching projects, datasources, datalabels, etc.
	beego.Router("/attribute/:param", projectController, "get:GetAttribute")
}
