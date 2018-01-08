package routers

import (
	"github.com/astaxie/beego"
	"ihome_go_2/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
