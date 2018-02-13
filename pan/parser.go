package pan

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const (
	bdKey = "B8ec24caf34ef7227c66767d29ffd3fb"
)

var (
	YunDataExp = regexp.MustCompile(`window\.yunData[\s]?=[\s]?(.*?);`)
)

func HmacSha1(key, origData []byte) (cipherText []byte) {
	mac := hmac.New(sha1.New, key)
	mac.Write(origData)
	return mac.Sum(nil)
}

func MustParseInt64(s string) (i int64) {
	i, _ = strconv.ParseInt(s, 10, 64)
	return
}

func MustParseInt(s string) (i int) {
	i, _ = strconv.Atoi(s)
	return
}

func (si *SharedInfo) auth() {
	// url := fmt.Sprintf(
	// 	"http://pan.baidu.com/share/list?shareid=%d&uk=%d&%s&sign=375dc1c3b4eeb35bb6e458f17e7f9c37e613ce76&timestamp=1618527500&bdstoken=76bf1ff30d7cc98550ce1a618ed2bc7e&devuid=&clienttype=1&channel=android_7.0_HUAWEI%%20NXT-AL10_bd-netdisk_1001540i&version=8.2.0",

	// )
	si.Timestamp = time.Now().Unix()
	orig := fmt.Sprintf("%d_%d__%d", si.ShareID, si.UK, si.Timestamp)

	si.Sign = HmacSha1([]byte(bdKey), []byte(orig))
}
