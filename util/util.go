package baiduUtil

import (
	"compress/gzip"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
	"unsafe"
)

var (
	// PipeInput 命令中是否为管道输入
	PipeInput bool
)

func init() {
	fileInfo, _ := os.Stdin.Stat()
	PipeInput = (fileInfo.Mode() & os.ModeNamedPipe) == os.ModeNamedPipe
}

// ToString 将 []byte 转换为 string
func ToString(p []byte) string {
	return *(*string)(unsafe.Pointer(&p))
}

// ToBytes 将 string 转换为 []byte
func ToBytes(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
}

/*
BeijingTimeOption 根据给定的 get 返回时间格式.

	get:        时间格式

	"Refer":    2017-7-21 12:02:32.000
	"printLog": 2017-7-21_12:02:32
	"day":      21
	"ymd":      2017-7-21
	"hour":     12
	默认时间戳:   1500609752
*/
func BeijingTimeOption(get string) string {
	//获取北京（东八区）时间
	CSTLoc := time.FixedZone("CST", 8*3600) // 东8区
	now := time.Now().In(CSTLoc)
	year, mon, day := now.Date()
	hour, min, sec := now.Clock()
	millisecond := now.Nanosecond() / 1e6
	switch get {
	case "Refer":
		return fmt.Sprintf("%d-%d-%d %02d:%02d:%02d.%03d", year, mon, day, hour, min, sec, millisecond)
	case "printLog":
		return fmt.Sprintf("%d-%d-%d_%02dh%02dm%02ds", year, mon, day, hour, min, sec)
	case "day":
		return fmt.Sprint(day)
	case "ymd":
		return fmt.Sprintf("%d-%d-%d", year, mon, day)
	case "hour":
		return fmt.Sprint(hour)
	default:
		return fmt.Sprint(time.Now().Unix())
	}
}

// GetURLCookieString 返回cookie字串
func GetURLCookieString(urlString string, jar *cookiejar.Jar) string {
	url, _ := url.Parse(urlString)
	cookies := jar.Cookies(url)
	cookieString := ""
	for _, v := range cookies {
		cookieString += v.String() + "; "
	}
	cookieString = strings.TrimRight(cookieString, "; ")
	return cookieString
}

// Md5Encrypt 对 str 进行md5加密, 返回值为 str 加密后的密文
func Md5Encrypt(str interface{}) string {
	md5Ctx := md5.New()
	switch value := str.(type) {
	case string:
		md5Ctx.Write(ToBytes(str.(string)))
	case *string:
		md5Ctx.Write(ToBytes(*str.(*string)))
	case []byte:
		md5Ctx.Write(str.([]byte))
	case *[]byte:
		md5Ctx.Write(*str.(*[]byte))
	default:
		fmt.Println("MD5Encrypt: undefined type:", value)
		return ""
	}
	return fmt.Sprintf("%X", md5Ctx.Sum(nil))
}

// DecompressGZIP 对 io.Reader 数据, 进行 gzip 解压
func DecompressGZIP(r io.Reader) ([]byte, error) {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	gzipReader.Close()
	return ioutil.ReadAll(gzipReader)
}

// FlagProvided 检测命令行是否提供名为 name 的 flag, 支持多个name(names)
func FlagProvided(names ...string) bool {
	if len(names) == 0 {
		return false
	}
	var targetFlag *flag.Flag
	for _, name := range names {
		targetFlag = flag.Lookup(name)
		if targetFlag == nil {
			return false
		}
		if targetFlag.DefValue == targetFlag.Value.String() {
			return false
		}
	}
	return true
}
