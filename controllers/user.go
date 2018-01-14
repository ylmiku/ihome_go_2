package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ihome_go_2/models"
	"path"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

type Name struct {
	Name string `json:"name"`
}

func (this *UserController) Reg() {
	beego.Info("=== reg controller is called =====")

	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	//1 得到客户端传递的信息 json解析
	//request
	var regRequestMap = make(map[string]interface{})
	json.Unmarshal(this.Ctx.Input.RequestBody, &regRequestMap)

	beego.Info("client reg request = ", regRequestMap)

	//2 校验信息的合法性 (mobile  password  sms_code)
	if regRequestMap["mobile"] == "" || regRequestMap["password"] == "" || regRequestMap["sms_code"] == "" {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//3  入库 user
	user := models.User{}
	user.Mobile = regRequestMap["mobile"].(string)
	user.Password_hash = regRequestMap["password"].(string)
	user.Name = regRequestMap["mobile"].(string)

	o := orm.NewOrm()
	id, err := o.Insert(&user)
	if err != nil {
		beego.Info("insert error = ", err)
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	beego.Info("reg succ!!! user id = ", id)

	//4 将用的存储到session中
	this.SetSession("name", user.Mobile)
	this.SetSession("user_id", id)
	this.SetSession("mobile", user.Mobile)

	return
}

// /api/v1.0/sessions [post]
//登陆也
func (this *UserController) Login() {
	beego.Info("=== login controller is called =====")

	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	//1 得到客户端传递的信息 json解析
	//request
	var loginRequestMap = make(map[string]interface{})
	json.Unmarshal(this.Ctx.Input.RequestBody, &loginRequestMap)

	beego.Info("client login request = ", loginRequestMap)

	//2 校验信息的合法性 (mobile  password  sms_code)
	if loginRequestMap["mobile"] == "" || loginRequestMap["password"] == "" {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//3  根据mobile查询user表
	// 3.1 如果没有数据，--->错误
	// 3,1 如果有数据， 对于user.passwd 和 regRequestMap["password"] 是否相等
	var user models.User

	o := orm.NewOrm()
	if err := o.QueryTable("user").Filter("mobile", loginRequestMap["mobile"].(string)).One(&user); err != nil {
		//没有任何数据
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	//比对密码
	if user.Password_hash != loginRequestMap["password"].(string) {
		//没有任何数据
		resp.Errno = models.RECODE_PWDERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	beego.Info("login succ!!! user id = ", user.Id)

	//4 将用的存储到session中
	this.SetSession("name", user.Mobile)
	this.SetSession("user_id", user.Id)
	this.SetSession("mobile", user.Mobile)

	return
}

// /api/v1.0/session [delete]
//用户点击退出，应该将该用户的session中的name删除
func (this *UserController) DelSessionName() {
	beego.Info("=== DelSessionName controller is called =====")
	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	this.DelSession("name")
	this.DelSession("user_id")
	this.DelSession("mobile")
}

// /api/v1.0/session [get]
//通过session 请求当前已经成功登陆或者注册的用户名
func (this *UserController) GetSessionName() {
	beego.Info("=== GetSession controller is called =====")

	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	name := this.GetSession("name")
	if name == nil {
		resp.Errno = models.RECODE_USERERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//nameData := Name{Name: name.(string)}
	//resp.Data = nameData
	nameMap := make(map[string]interface{})
	nameMap["name"] = name.(string)
	resp.Data = nameMap

	return
}

// api/v1.0/user/avatar   [post]
func (this *UserController) UploadAvatar() {
	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	//拿到用户的文件二进制数据
	file, header, err := this.GetFile("avatar")
	if err != nil {
		resp.Errno = models.RECODE_SERVERERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	//创建一个file文件的buffer
	fileBuffer := make([]byte, header.Size)

	_, err = file.Read(fileBuffer)
	if err != nil {
		resp.Errno = models.RECODE_IOERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//user01.jpg
	suffix := path.Ext(header.Filename) // --> ".jpg"

	groupName, fileId, err1 := models.FDFSUploadByBuffer(fileBuffer, suffix[1:])
	if err1 != nil {
		resp.Errno = models.RECODE_IOERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Info("upload file error , name = ", header.Filename)
		return
	}

	beego.Info("fdfs upload file succ  gourpname ", groupName, " fileid = ", fileId)

	//通过session 得到当前用的user_id
	user_id := this.GetSession("user_id")

	user := models.User{Id: user_id.(int), Avatar_url: fileId}
	o := orm.NewOrm()

	if _, err := o.Update(&user, "avatar_url"); err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//将fileid 拼接一个完整的url路径 + ip + port 返回给前端

	avatar_url := "http://192.168.48.129:9091/" + fileId

	url_map := make(map[string]interface{})
	url_map["avatar_url"] = avatar_url
	resp.Data = url_map

	return
}
