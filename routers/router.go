package routers

import (
	"github.com/astaxie/beego"
	"ihome_go_2/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/v1.0/areas", &controllers.AreaController{}, "get:GetAreas")
}
