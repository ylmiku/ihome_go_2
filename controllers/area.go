package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ihome_go_2/models"
)

type AreaController struct {
	beego.Controller
}

//针对area请求返回数据的格式
type AreaResp struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data,omitempty"`
}

func (this *AreaController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

// /api/v1.0/area [get]
func (this *AreaController) GetAreas() {
	beego.Info("=== area controller is called =====")

	resp := AreaResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)

	//1 应该从缓存中取得地域信息数据 直接返回给前端

	// 2 如果缓存没有数据， 从mysql中 查询area数据
	o := orm.NewOrm()

	var areas []models.Area

	qs := o.QueryTable("area")
	num, err := qs.All(&areas)
	if err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//3 将area数据 存入缓存

	// 4 将area数据 变成json  发送给前端
	if num == 0 {
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	beego.Info("areas = ", areas)

	resp.Data = areas

	return
}
