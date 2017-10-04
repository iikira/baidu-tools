package baidu

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/iikira/baidu-tools/util"
)

// NewUser 返回 Baidu 指针
func NewUser(uid, name string) (b *Baidu, err error) {
	b = &Baidu{
		UID:  uid,
		Name: name,
	}
	err = b.GetUserInfo()
	return b, err
}

// GetUserInfo 获取百度帐号详细信息, 优先级: name, uid
func (b *Baidu) GetUserInfo() error {
	switch {
	case b.UID == "" && b.Name == "":
		return fmt.Errorf("NewUser: name and uid is nil")
	case b.UID != "" && b.Name != "":
		fallthrough
	case b.UID == "" && b.Name != "":
		body, err := baiduUtil.HTTPGet("http://tieba.baidu.com/home/get/panel?un=" + b.Name)
		if err != nil {
			return err
		}
		json, err := simplejson.NewJson(body)
		if err != nil {
			return err
		}
		userJSON := json.Get("data")
		byteUID, err := userJSON.Get("id").MarshalJSON()
		if err != nil {
			return fmt.Errorf("Json 解析错误: %s", err)
		}
		b.UID = string(byteUID)
		b.Age = userJSON.Get("tb_age").MustString()
		b.NameShow = userJSON.Get("name_show").MustString()
		sex := userJSON.Get("sex").MustString()
		switch sex {
		case "male":
			b.Sex = "♂"
		case "female":
			b.Sex = "♀"
		default:
			b.Sex = "unknown"
		}
	case b.UID != "" && b.Name == "":
		rawQuery := "has_plist=0&need_post_count=1&rn=1&uid=" + b.UID
		urlStr := "http://c.tieba.baidu.com/c/u/user/profile?" + baiduUtil.TiebaClientRawQuerySignature(rawQuery)
		body, err := baiduUtil.HTTPGet(urlStr)
		if err != nil {
			return err
		}
		json, err := simplejson.NewJson(body)
		if err != nil {
			return err
		}
		userJSON := json.GetPath("user")
		b.Name = userJSON.Get("name").MustString()
		b.NameShow = userJSON.Get("name_show").MustString()
		b.Age = userJSON.Get("tb_age").MustString()
		sex := userJSON.Get("sex").MustInt()
		switch sex {
		case 1:
			b.Sex = "♂"
		case 2:
			b.Sex = "♀"
		default:
			b.Sex = "unknown"
		}
	}
	return nil
}

func NewAuth(bduss string) (*Auth, error) {
	a := &Auth{
		BDUSS: bduss,
	}
	err := a.getTbs()
	if err != nil {
		return nil, err
	}
	return a, nil
}

// getTbs 获取贴吧TBS
func (a *Auth) getTbs() error {
	if a.BDUSS == "" {
		return fmt.Errorf("获取TBS出错: BDUSS为空")
	}
	hd := map[string]string{
		"Cookie": "BDUSS=" + a.BDUSS,
	}
	body, err := baiduUtil.Fetch("GET", "http://tieba.baidu.com/dc/common/tbs", nil, nil, hd)
	if err != nil {
		return fmt.Errorf("获取TBS出错: %s", err)
	}
	json, err := simplejson.NewJson(body)
	if err != nil {
		return fmt.Errorf("json解析出错: %s", err)
	}
	a.Tbs = json.Get("tbs").MustString()
	return nil
}
