package controllers

import (
	"github.com/astaxie/beego"
	"ihome_go_2/models"
)

type HousesIndexController struct {
	beego.Controller
}

func (this *HousesIndexController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

// /api/v1.0/houses/index [get]
func (this *HousesIndexController) HousesIndex() {
	beego.Info("=== HousesIndex controller is called =====")

	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	return
}
