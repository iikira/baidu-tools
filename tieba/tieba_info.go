package tieba

import (
	"bytes"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/iikira/baidu-tools/util"
	"log"
	"regexp"
	"strconv"
)

// GetBars 获取贴吧列表
func (user *Tieba) GetBars() error {
	bars, err := GetBars(user.Baidu.UID)
	if err != nil {
		return err
	}
	user.Bars = bars
	return nil
}

// GetBars 获取贴吧列表
func GetBars(uid string) ([]Bar, error) {
	if _, err := strconv.Atoi(uid); err != nil {
		return nil, fmt.Errorf("百度 UID 非法")
	}
	var (
		pageNo uint16
		bars   []Bar
	)
	bajsonRE := regexp.MustCompile("{\"id\":\".+?\"}")
	for {
		pageNo++
		rawQuery := fmt.Sprintf("_client_version=6.9.2.1&page_no=%d&page_size=200&uid=%s", pageNo, uid)
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
