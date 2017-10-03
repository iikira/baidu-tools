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
	if _, err := strconv.Atoi(user.UID); err != nil {
		return err
	}
	var pageNo uint16
	bajsonRE := regexp.MustCompile("{\"id\":\".+?\"}")
	for {
		pageNo++
		rawQuery := fmt.Sprintf("_client_version=6.9.2.1&page_no=%d&page_size=200&uid=%s", pageNo, user.UID)
		//贴吧客户端签名
		body, err := baiduUtil.HTTPGet("http://c.tieba.baidu.com/c/f/forum/like?" + baiduUtil.TiebaClientRawQuerySignature(rawQuery))
		if err != nil {
			return err
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
				return err
			}
			if curScore, ok := bajson.CheckGet("cur_score"); ok {
				exp, _ := strconv.Atoi(curScore.MustString())
				user.Bars = append(user.Bars, Bar{
					Fid:   bajson.Get("id").MustString(),
					Name:  bajson.Get("name").MustString(),
					Level: bajson.Get("level_id").MustString(),
					Exp:   exp,
				})
			}
		}
	}
	return nil
}

// GetTiebaFid 获取贴吧fid值
func GetTiebaFid(tiebaName string) (fid string, err error) {
	b, err := baiduUtil.HTTPGet("http://tieba.baidu.com/f/commit/share/fnameShareApi?ie=utf-8&fname=" + tiebaName)
	if err != nil {
		return
	}
	c := regexp.MustCompile(`"data":{"fid":(.*?),"can_send_pics":`).FindSubmatch(b)
	if len(c) > 1 {
		fid = string(c[1])
	}
	return
}

// IsTiebaExist 检测贴吧是否存在
func IsTiebaExist(tiebaName string) bool {
	b, err := baiduUtil.HTTPGet("http://c.tieba.baidu.com/mo/q/m?tn4=bdKSW&sub4=&word=" + tiebaName)
	if err != nil {
		log.Println(err)
	}
	return !bytes.Contains(b, []byte(`class="tip_text2">欢迎创建此吧，和朋友们在这里交流</p>`))
}
