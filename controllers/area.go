package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"ihome_go_2/models"
	"time"
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

	resp := make(map[string]interface{})

	//resp := AreaResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)

	//1 应该从缓存中取得地域信息数据 直接返回给前端
	//(1) 创建redis链接
	cache_conn, err := cache.NewCache("redis", `{"key":"ihome_go_2","conn":"127.0.0.1:6379","dbNum":"0"}`)
	if err != nil {
		beego.Info("cache redis conn error , err =", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	//(2) 尝试从redis中 取出area的值
	areas_info_value := cache_conn.Get("area_info")
	if areas_info_value != nil {
		//(3) 如果有 直接将area的值返回个前端
		beego.Info(" ======= get area_info from cache !!!! ==== ")

		//将areas_info_value(字符串) -->  area 结构体 slice
		var area_info interface{}
		json.Unmarshal(areas_info_value.([]byte), &area_info)
		resp["data"] = area_info
		return
	}

	//    如果没有 则继续执行

	/*
		if err := cache_conn.Put("beegorediskey", "hahah", time.Second*30); err != nil {
			beego.Info("put cache error")
		}
	*/

	// 2 如果缓存没有数据， 从mysql中 查询area数据
	o := orm.NewOrm()

	var areas []models.Area

	qs := o.QueryTable("area")
	num, err := qs.All(&areas)
	if err != nil {
		/*
			resp.Errno = models.RECODE_DBERR
			resp.Errmsg = models.RecodeText(resp.Errno)
		*/
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	//3 将area数据 存入缓存
	//将area数据编程json
	areas_info_str, _ := json.Marshal(areas) //[{}, {}, {}]
	if err := cache_conn.Put("area_info", areas_info_str, time.Second*3600); err != nil {
		beego.Info("set area_info to cache error, err = ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	// 4 将area数据   发送给前端
	if num == 0 {
		/*
			resp.Errno = models.RECODE_NODATA
			resp.Errmsg = models.RecodeText(resp.Errno)
		*/
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(resp["errno"].(string))
		return
	}

	beego.Info("areas = ", areas)

	resp["data"] = areas

	return
}
