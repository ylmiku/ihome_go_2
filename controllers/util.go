package controllers

//针对请求返回数据的格式
type Resp struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}
