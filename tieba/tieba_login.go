package tieba

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/iikira/baidu-tools"
	"github.com/iikira/baidu-tools/util"
	"time"
)

// NewWithBDUSS 检测BDUSS有效性, 同时获取百度详细信息
func NewWithBDUSS(bduss string) (*Tieba, error) {
	post := map[string]string{
		"bdusstoken":  bduss + "|null",
		"channel_id":  "",
		"channel_uid": "",
		"stErrorNums": "0",
		"subapp_type": "mini",
		"timestamp":   baiduUtil.BeijingTimeOption("") + "922",
	}
	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Cookie":       "ka=open",
		"net":          "1",
		"User-Agent":   "bdtb for Android 6.9.2.1",
		"client_logid": baiduUtil.BeijingTimeOption("") + "416",
		"Connection":   "Keep-Alive",
	}
	baiduUtil.TiebaClientSignature(post)

	var (
		body []byte
		err  error
	)
	for errorTimes := 0; errorTimes <= 3; errorTimes++ { // 错误重试
		body, err = baiduUtil.Fetch("POST", "http://tieba.baidu.com/c/s/login", nil, post, header) // 获取百度ID的TBS，UID，BDUSS等
		if err == nil {
			break
		}
		if errorTimes >= 3 {
			return nil, fmt.Errorf("检测帐号状态失败, 错误次数超过3次: %s", err)
		}
		time.Sleep(1e9)
	}
	json, err := simplejson.NewJson(body)
	if err != nil {
		return nil, fmt.Errorf("json解析出错: %s", err)
	}

	errCode := json.Get("error_code").MustString()
	errMsg := json.Get("error_msg").MustString()
	if errCode != "0" {
		return nil, fmt.Errorf("错误代码: %s, 消息: %s", baiduUtil.ErrorColor(errCode), baiduUtil.ErrorColor(errMsg))
	}
	u := Tieba{
		Baidu: baidu.Baidu{
			UID:  json.GetPath("user", "id").MustString(),
			Name: json.GetPath("user", "name").MustString(),
			Auth: baidu.Auth{
				BDUSS: bduss,
				Tbs:   json.GetPath("anti", "tbs").MustString(),
			},
		},
	}
	err = u.GetUserInfo()
	return &u, err
}
