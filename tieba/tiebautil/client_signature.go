package tiebautil

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
)

// TiebaClientSignature 根据给定贴吧客户端的 post数据 进行签名, 以通过百度服务器验证。返回值为 sign 签名字符串值
func TiebaClientSignature(post map[string]string) {
	if post == nil {
		return
	}

	// 预设
	post["_client_type"] = "2"
	post["_client_version"] = "6.9.2.1"
	post["_phone_imei"] = "860983036542682"
	post["from"] = "mini_ad_wandoujia"
	post["model"] = "HUAWEI NXT-AL10"
	post["cuid"] = "61464018582906C485355A89D105ECFB|286245630389068"

	var keys []string
	for key := range post {
		keys = append(keys, key)
	}
	sort.Sort(sort.StringSlice(keys))

	var bb bytes.Buffer
	for _, key := range keys {
		bb.WriteString(key + "=" + post[key])
	}
	bb.WriteString("tiebaclient!!!")
	post["sign"] = fmt.Sprintf("%x", md5.Sum(bb.Bytes()))
}

// TiebaClientRawQuerySignature 给 rawQuery 进行贴吧客户端签名, 返回值为签名后的 rawQuery
func TiebaClientRawQuerySignature(rawQuery string) (signedRawQuery string) {
	return rawQuery + "&sign=" + fmt.Sprintf("%x", md5.Sum([]byte(strings.Replace(rawQuery, "&", "", -1)+"tiebaclient!!!")))
}
