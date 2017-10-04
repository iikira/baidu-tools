package tieba

import (
	"bytes"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/iikira/baidu-tools"
	"github.com/iikira/baidu-tools/util"
	"log"
	"regexp"
	"strconv"
)

// NewUserInfoByUID 提供 UID 获取百度帐号详细信息
func NewUserInfoByUID(uid uint64) (t *Tieba, err error) {
	b := baidu.NewUser(uid, "")
	rawQuery := "has_plist=0&need_post_count=1&rn=1&uid=" + fmt.Sprint(b.UID)
	urlStr := "http://c.tieba.baidu.com/c/u/user/profile?" + baiduUtil.TiebaClientRawQuerySignature(rawQuery)
	body, err := baiduUtil.HTTPGet(urlStr)
	if err != nil {
		return nil, err
	}
	json, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}
	userJSON := json.GetPath("user")
	b.Name = userJSON.Get("name").MustString()
	b.NameShow = userJSON.Get("name_show").MustString()
	b.Age = userJSON.Get("tb_age").MustFloat64()
	sex := userJSON.Get("sex").MustInt()
	switch sex {
	case 1:
		b.Sex = "♂"
	case 2:
		b.Sex = "♀"
	default:
		b.Sex = "unknown"
	}

	t = &Tieba{
		Baidu: b,
		Stat: &Stat{
			LikeForumNum: userJSON.Get("like_forum_num").MustInt(),
			PostNum:      userJSON.Get("post_num").MustInt(),
		},
	}
	return t, nil
}

// NewUserInfoByName 提供 name (百度用户名) 获取百度帐号详细信息
func NewUserInfoByName(name string) (t *Tieba, err error) {
	body, err := baiduUtil.HTTPGet("http://tieba.baidu.com/home/get/panel?un=" + name)
	if err != nil {
		return nil, err
	}
	json, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}
	return NewUserInfoByUID(json.GetPath("data", "id").MustUint64())
}

// FlushUserInfo 提供 name (百度用户名) 获取百度帐号详细信息
func (t *Tieba) FlushUserInfo(uids ...uint64) error {
	switch len(uids) {
	case 0:
	case 1:
		t.Baidu.UID = uids[0]
	}
	this, err := NewUserInfoByUID(t.Baidu.UID)
	if err != nil {
		return err
	}
	this.Baidu.Auth = t.Baidu.Auth
	t.Baidu = this.Baidu
	t.Stat = this.Stat
	return nil
}

// GetBars 获取贴吧列表
func (t *Tieba) GetBars() error {
	bars, err := GetBars(t.Baidu.UID)
	if err != nil {
		return err
	}
	t.Bars = bars
	return nil
}

// GetBars 获取贴吧列表
func GetBars(uid uint64) ([]Bar, error) {
	var (
		pageNo uint16
		bars   []Bar
	)
	bajsonRE := regexp.MustCompile("{\"id\":\".+?\"}")
	for {
		pageNo++
		rawQuery := fmt.Sprintf("_client_version=6.9.2.1&page_no=%d&page_size=200&uid=%d", pageNo, uid)
		//贴吧客户端签名
		body, err := baiduUtil.HTTPGet("http://c.tieba.baidu.com/c/f/forum/like?" + baiduUtil.TiebaClientRawQuerySignature(rawQuery))
		if err != nil {
			return nil, fmt.Errorf("获取贴吧列表网络错误, %s", err)
		}
		if !bytes.Contains(body, []byte("has_more")) { // 贴吧服务器响应有误, 再试一次
			pageNo--
			continue
		}
		jsonSlice := bajsonRE.FindAll(body, -1)
		if jsonSlice == nil { // 完成抓去贴吧列表
			break
		}
		for _, bajsonStr := range jsonSlice {
			bajson, err := simplejson.NewJson(bajsonStr)
			if err != nil {
				return nil, fmt.Errorf("获取贴吧列表json解析错误, %s", err)
			}
			if curScore, ok := bajson.CheckGet("cur_score"); ok {
				exp, _ := strconv.Atoi(curScore.MustString())
				bars = append(bars, Bar{
					Fid:   bajson.Get("id").MustString(),
					Name:  bajson.Get("name").MustString(),
					Level: bajson.Get("level_id").MustString(),
					Exp:   exp,
				})
			}
		}
	}
	return bars, nil
}

// GetTiebaFid 获取贴吧fid值
func GetTiebaFid(tiebaName string) (fid string, err error) {
	b, err := baiduUtil.HTTPGet("http://tieba.baidu.com/f/commit/share/fnameShareApi?ie=utf-8&fname=" + tiebaName)
	if err != nil {
		return "", fmt.Errorf("获取贴吧fid网络错误, %s", err)
	}
	json, err := simplejson.NewJson(b)
	if err != nil {
		return "", fmt.Errorf("获取贴吧fid json解析错误, %s", err)
	}
	intFid := json.GetPath("data", "fid").MustInt()
	return fmt.Sprint(intFid), nil
}

// IsTiebaExist 检测贴吧是否存在
func IsTiebaExist(tiebaName string) bool {
	b, err := baiduUtil.HTTPGet("http://c.tieba.baidu.com/mo/q/m?tn4=bdKSW&sub4=&word=" + tiebaName)
	if err != nil {
		log.Println(err)
	}
	return !bytes.Contains(b, []byte(`class="tip_text2">欢迎创建此吧，和朋友们在这里交流</p>`))
}
