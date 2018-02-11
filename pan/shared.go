// 百度网盘提取分享文件的下载链接
package pan

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/requester"
	"github.com/json-iterator/go"
	"net/url"
	"path"
	"strconv"
	"strings"
)

// SharedInfo 百度网盘文件分享页信息
type SharedInfo struct {
	UK            uint64 `json:"uk"`
	ShareID       uint64 `json:"shareid"`
	RootSharePath string `json:"rootSharePath"` // 分享的目录, 基于分享者的网盘根目录
	DownloadSign  string `json:"downloadsign"`
	TimeStamp     uint64 `json:"timestamp"`

	BaseFileList []*FileDirectoryString `json:"file_list"`

	client *requester.HTTPClient
}

// NewSharedInfo 解析百度网盘文件分享页信息,
// 暂不支持带提取码的的分享
func NewSharedInfo(sharedURL string) (si *SharedInfo, err error) {
	h := requester.NewHTTPClient()
	h.SetKeepAlive(false)
	h.UserAgent = "Mozilla/5.0 (Linux; Android 7.0; HUAWEI NXT-AL10 Build/HUAWEINXT-AL10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.137 Mobile Safari/537.36"

	si = &SharedInfo{
		client: h,
	}

	body, err := si.client.Fetch("GET", sharedURL, nil, nil)
	if err != nil {
		return nil, err
	}

	rawYunData := YunDataExp.FindSubmatch(body)
	if len(rawYunData) < 2 {
		return nil, fmt.Errorf("分享页数据解析失败")
	}

	err = jsoniter.Unmarshal(rawYunData[1], si)
	if err != nil {
		return nil, fmt.Errorf("json 数据解析失败, %s", err)
	}

	return si, nil
}

// FileDirectoryString 文件和目录的信息, 字段类型全为 string
type FileDirectoryString struct {
	FsID     string `json:"fs_id"`           // fs_id
	Path     string `json:"path"`            // 路径
	Filename string `json:"server_filename"` // 文件名 或 目录名
	Ctime    string `json:"server_ctime"`    // 创建日期
	Mtime    string `json:"server_mtime"`    // 修改日期
	MD5      string `json:"md5"`             // md5 值
	Size     string `json:"size"`            // 文件大小 (目录为0)
	Isdir    string `json:"isdir"`           // 是否为目录
}

// FileDirectory 文件和目录的信息
type FileDirectory struct {
	FsID     int64  `json:"fs_id"`           // fs_id
	Path     string `json:"path"`            // 路径
	Filename string `json:"server_filename"` // 文件名 或 目录名
	Ctime    int64  `json:"server_ctime"`    // 创建日期
	Mtime    int64  `json:"server_mtime"`    // 修改日期
	MD5      string `json:"md5"`             // md5 值
	Size     int64  `json:"size"`            // 文件大小 (目录为0)
	Isdir    int    `json:"isdir"`           // 是否为目录
}

// List 获取文件列表
func (si *SharedInfo) List(subDir string) (fds []*FileDirectory, err error) {
	dir := path.Clean(si.RootSharePath + "/" + subDir)

	escapedDir := url.PathEscape(dir)

	listURL := "https://pan.baidu.com/share/list?app_id=250528&page=1&num=200&dir=" + escapedDir + "&uk=" + strconv.FormatUint(si.UK, 10) + "&shareid=" + strconv.FormatUint(si.ShareID, 10)

	body, err := si.client.Fetch("GET", listURL, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("获取文件列表网络错误, %s", err)
	}

	jsonData := &struct {
		ErrNo int
		List  []*FileDirectory `json:"list"`
	}{}

	err = jsoniter.Unmarshal(body, jsonData)
	if err != nil {
		return nil, fmt.Errorf("获取文件列表, json 数据解析失败, %s", err)
	}

	msgFmt := "获取文件列表, 远端服务器返回错误代码 " + fmt.Sprint(jsonData.ErrNo) + ", 消息: %s"
	switch jsonData.ErrNo {
	case 0:
	case -9:
		return nil, fmt.Errorf(msgFmt, "可能路径不存在")
	case -19:
		return nil, fmt.Errorf(msgFmt, "需要输入验证码")
	default:
		fmt.Printf("%s\n", body)
		return nil, fmt.Errorf(msgFmt, "未知错误")
	}

	return jsonData.List, err
}

// GetFSID 获取文件的 fs_id
func (si *SharedInfo) GetFSID(filePath string) (fsid int64, err error) {
	if filePath == "" {
		return 0, fmt.Errorf("文件路径为空")
	}

	cleanedPath := path.Clean(filePath)
	if cleanedPath == "/" {
		return 0, fmt.Errorf("不支持获取根目录fsid")
	}

	dir, fileName := path.Split(cleanedPath)

	if dir == "/" || dir == "" {
		for _, fileInfo := range si.BaseFileList {
			// 考虑使用通配符
			if strings.Compare(fileInfo.Filename, fileName) == 0 {
				return strconv.ParseInt(fileInfo.FsID, 10, 64)
			}
		}

		return 0, fmt.Errorf("未匹配到文件路径 %s", filePath)
	}

	dirInfo, err := si.List(dir)
	if err != nil {
		return 0, fmt.Errorf("获取文件的 fs_id 出错, %s", err)
	}

	for k := range dirInfo {
		if strings.Compare(dirInfo[k].Filename, fileName) == 0 {
			return dirInfo[k].FsID, nil
		}
	}

	return 0, fmt.Errorf("未匹配到文件路径 %s", filePath)
}

// GetDownloadLink 获取下载直链,
// 此接口较大概率会出现验证码
func (si *SharedInfo) GetDownloadLink(fsid int64) (dlink string, err error) {
	fidList := url.PathEscape(fmt.Sprintf("[%d]", fsid))
	getDlinkURL := fmt.Sprintf("http://pan.baidu.com/share/download?app_id=250528&uk=%d&shareid=%d&fid_list=%s&sign=%s&timestamp=%d", si.UK, si.ShareID, fidList, si.DownloadSign, si.TimeStamp)

	body, err := si.client.Fetch("GET", getDlinkURL, nil, nil)
	if err != nil {
		return "", fmt.Errorf("获取下载直链网络错误, %s", err)
	}

	jsonData := &struct {
		ErrNo int    `json:"errno"`
		Dlink string `json:"dlink"`
		Img   string `json:"img"`
		Vcode string `json:"vcode"`
	}{}

	err = jsoniter.Unmarshal(body, jsonData)
	if err != nil {
		return "", fmt.Errorf("获取下载直链 json 解析错误, %s", err)
	}

	msgFmt := "获取下载直链, 远端服务器返回错误代码 " + fmt.Sprint(jsonData.ErrNo) + ", 消息: %s"
	switch jsonData.ErrNo {
	case 0:
	case -19:
		return "", fmt.Errorf(msgFmt, "需要输入验证码")
	default:
		return "", fmt.Errorf(msgFmt, "未知错误")
	}

	return jsonData.Dlink, nil
}
