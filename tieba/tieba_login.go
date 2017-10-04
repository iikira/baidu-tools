package tieba

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/iikira/baidu-tools"
	"github.com/iikira/baidu-tools/util"
	"strconv"
	"time"
)

// NewWithBDUSS 检测BDUSS有效性, 同时获取百度详细信息
func NewWithBDUSS(bduss string) (*Tieba, error) {
	timestamp := baiduUtil.BeijingTimeOption("")
	post := map[string]string{
		"bdusstoken":  bduss + "|null",
		"channel_id":  "",
		"channel_uid": "",
		"stErrorNums": "0",
		"subapp_type": "mini",
		"timestamp":   timestamp + "922",
	}
	baiduUtil.TiebaClientSignature(post)

	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Cookie":       "ka=open",
		"net":          "1",
		"User-Agent":   "bdtb for Android 6.9.2.1",
		"client_logid": timestamp + "416",
		"Connection":   "Keep-Alive",
	}

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
			return nil, fmt.Errorf("检测BDUSS有效性失败, 错误次数超过3次: %s", err)
		}
		time.Sleep(1e9)
	}
	json, err := simplejson.NewJson(body)
	if err != nil {
		return nil, fmt.Errorf("检测BDUSS有效性json解析出错: %s", err)
	}
	errCode := json.Get("error_code").MustString()
	errMsg := json.Get("error_msg").MustString()
	if errCode != "0" {
		return nil, fmt.Errorf("检测BDUSS有效性错误代码: %s, 消息: %s", baiduUtil.ErrorColor(errCode), baiduUtil.ErrorColor(errMsg))
	}
	uidStr := json.GetPath("user", "id").MustString()
	uid, _ := strconv.ParseUint(uidStr, 10, 64)

	t := &Tieba{
		Baidu: &baidu.Baidu{
			UID:  uid,
			Name: json.GetPath("user", "name").MustString(),
			Auth: baidu.NewAuth(bduss, "", ""),
		},
		Tbs: json.GetPath("anti", "tbs").MustString(),
	}
	err = t.FlushUserInfo()
	if err != nil {
		return nil, err
	}
	return t, nil
}

// GetTbs 获取贴吧TBS
func (t *Tieba) GetTbs() error {
	bduss := t.Baidu.Auth.BDUSS
	if bduss == "" {
		return fmt.Errorf("获取贴吧TBS出错: BDUSS为空")
	}
	tbs, err := GetTbs(bduss)
	if err != nil {
		return err
	}
	t.Tbs = tbs
	return nil
}

// GetTbs 获取贴吧TBS
func GetTbs(bduss string) (tbs string, err error) {
	body, err := baiduUtil.Fetch("GET", "http://tieba.baidu.com/dc/common/tbs", nil, nil, map[string]string{
		"Cookie": "BDUSS=" + bduss,
	})
	if err != nil {
		return "", fmt.Errorf("获取贴吧TBS网络错误: %s", err)
	}
	json, err := simplejson.NewJson(body)
	if err != nil {
		return "", fmt.Errorf("获取贴吧TBS json解析出错: %s", err)
	}
	return json.Get("tbs").MustString(), nil
}
