package tieba

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"my/Baidu-Tools/util"
	"strconv"
)

// TiebaSign 执行贴吧签到
func (user *Tieba) TiebaSign(fid, name string) (status int, bonusExp int, err error) {
	timestamp := baiduUtil.BeijingTimeOption("")
	post := map[string]string{
		"BDUSS":       user.BDUSS,
		"_client_id":  "wappc_" + timestamp + "150_607",
		"fid":         fid,
		"kw":          name,
		"stErrorNums": "1",
		"stMethod":    "1",
		"stMode":      "1",
		"stSize":      "229",
		"stTime":      "185",
		"stTimesNum":  "1",
		"subapp_type": "mini",
		"tbs":         user.Tbs,
		"timestamp":   timestamp + "083",
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

	body, err := baiduUtil.Fetch("POST", "http://c.tieba.baidu.com/c/c/forum/sign", nil, post, header)
	if err != nil {
		return 1, 0, fmt.Errorf("POST 错误: %s", err)
	}
	json, err := simplejson.NewJson(body)
	if err != nil {
		return 1, 0, fmt.Errorf("json解析错误: %s", err)
	}
	if signBonusPoint, ok := json.Get("user_info").CheckGet("sign_bonus_point"); ok { // 签到成功, 获取经验
		bonusExp, _ = strconv.Atoi(signBonusPoint.MustString())
		return 0, bonusExp, nil
	}
	errorCode := json.Get("error_code").MustString()
	errorMsg := json.Get("error_msg").MustString()
	err = fmt.Errorf("贴吧签到时发生错误, 错误代码: %s, 消息: %s", baiduUtil.ErrorColor(errorCode), baiduUtil.ErrorColor(errorMsg))
	switch errorCode {
	case "340010", "160002", "3": // 已签到
		return 0, 0, nil
	case "110001": // 签名错误
		return 1, 0, nil
	case "340011": // 操作太快
		return 2, 0, err
	case "340008", "340006", "3250002": // 340008黑名单, 340006封吧, 3250002永久封号
		return 3, 0, err
	case "1", "1990055": // 1掉线, 1990055未实名
		return 4, 0, err
	default:
		if errorMsg == "" {
			errorMsg = "未找到错误原因, 请检查：" + baiduUtil.ToString(body)
		}
		return 1, 0, err
	}
}
