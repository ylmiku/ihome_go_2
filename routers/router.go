package routers

import (
	"ihome_go_2/ihome_go_2/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
