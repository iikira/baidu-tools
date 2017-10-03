package baidu

import (
	"errors"
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
		return errors.New("NewUser: name and uid is nil")
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
			return errors.New("Json 解析错误: " + err.Error())
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
