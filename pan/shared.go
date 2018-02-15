// 百度网盘提取分享文件的下载链接
package pan

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/requester"
	"github.com/json-iterator/go"
	"net/url"
	"path"
	"strings"
)

// SharedInfo 百度网盘文件分享页信息
type SharedInfo struct {
	UK            uint64 `json:"uk"`            // 百度网盘用户id
	ShareID       uint64 `json:"shareid"`       // 分享id
	RootSharePath string `json:"rootSharePath"` // 分享的目录, 基于分享者的网盘根目录

	Timestamp int64  // unix 时间戳
	Sign      []byte // 签名

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

	si.Signature()

	return si, nil
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
	Dlink    string `json:"dlink"`           //下载直链
}

// fileDirectoryString 文件和目录的信息, 字段类型全为 string
type fileDirectoryString struct {
	FsID     string `json:"fs_id"`           // fs_id
	Path     string `json:"path"`            // 路径
	Filename string `json:"server_filename"` // 文件名 或 目录名
	Ctime    string `json:"server_ctime"`    // 创建日期
	Mtime    string `json:"server_mtime"`    // 修改日期
	MD5      string `json:"md5"`             // md5 值
	Size     string `json:"size"`            // 文件大小 (目录为0)
	Isdir    string `json:"isdir"`           // 是否为目录
	Dlink    string `json:"dlink"`           // 下载链接
}

func (fdss *fileDirectoryString) convert() *FileDirectory {
	return &FileDirectory{
		FsID:     MustParseInt64(fdss.FsID),
		Path:     fdss.Path,
		Filename: fdss.Filename,
		Ctime:    MustParseInt64(fdss.Ctime),
		Mtime:    MustParseInt64(fdss.Mtime),
		MD5:      fdss.MD5,
		Size:     MustParseInt64(fdss.Size),
		Isdir:    MustParseInt(fdss.Isdir),
		Dlink:    fdss.Dlink,
	}
}

// List 获取文件列表, subDir 为相对于分享目录的目录
func (si *SharedInfo) List(subDir string) (fds []*FileDirectory, err error) {
	var (
		isRoot     = 0
		escapedDir string
	)

	cleanedSubDir := path.Clean(subDir)
	if cleanedSubDir == "." || cleanedSubDir == "/" {
		isRoot = 1
	} else {
		dir := path.Clean(si.RootSharePath + "/" + subDir)
		escapedDir = url.PathEscape(dir)
	}

	listURL := fmt.Sprintf(
		"http://pan.baidu.com/share/list?shareid=%d&uk=%d&root=%d&dir=%s&sign=%x&timestamp=%d&devuid=&clienttype=1&channel=android_7.0_HUAWEI%%20NXT-AL10_bd-netdisk_1001540i&version=8.2.0",
		si.ShareID, si.UK,
		isRoot, escapedDir,
		si.Sign, si.Timestamp,
	)

	fmt.Println(listURL)

	body, err := si.client.Fetch("GET", listURL, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("获取文件列表网络错误, %s", err)
	}

	var errNo int
	if isRoot != 0 { // 根目录
		jsonData := &struct {
			ErrNo int
			List  []*fileDirectoryString `json:"list"`
		}{}

		err = jsoniter.Unmarshal(body, jsonData)
		if err == nil {
			fds = make([]*FileDirectory, len(jsonData.List))
			for k, info := range jsonData.List {
				fds[k] = info.convert()
			}
		}

		errNo = jsonData.ErrNo
	} else {
		jsonData := &struct {
			ErrNo int
			List  []*FileDirectory `json:"list"`
		}{}

		err = jsoniter.Unmarshal(body, jsonData)
		if err == nil {
			errNo = jsonData.ErrNo
			fds = jsonData.List
		}
	}

	if err != nil {
		return nil, fmt.Errorf("获取文件列表, json 数据解析失败, %s", err)
	}

	msgFmt := "获取文件列表, 远端服务器返回错误代码 " + fmt.Sprint(errNo) + ", 消息: %s"
	switch errNo {
	case 0:
	case -9:
		return nil, fmt.Errorf(msgFmt, "可能路径不存在或提取码错误")
	case -19:
		return nil, fmt.Errorf(msgFmt, "需要输入验证码")
	default:
		fmt.Printf("%s\n", body)
		return nil, fmt.Errorf(msgFmt, "未知错误")
	}

	return fds, nil
}

// GetDownloadLink 获取下载直链, filePath 为相对于分享目录的目录
func (si *SharedInfo) GetDownloadLink(filePath string) (dlink string, err error) {
	cleanedPath := path.Clean(filePath)
	if cleanedPath == "/" || cleanedPath == "." {
		return "", fmt.Errorf("不支持获取根目录下载直链")
	}

	dir, fileName := path.Split(cleanedPath)

	dirInfo, err := si.List(dir)
	if err != nil {
		return "", fmt.Errorf("获取目录信息出错, 路径: %s, %s", path.Clean(dir), err)
	}

	for k := range dirInfo {
		if strings.Compare(dirInfo[k].Filename, fileName) == 0 {
			if dirInfo[k].Isdir != 0 {
				return "", fmt.Errorf("不支持获取目录的下载直链, 路径: %s", cleanedPath)
			}
			return dirInfo[k].Dlink, nil
		}
	}

	return "", fmt.Errorf("未匹配到文件路径 %s", cleanedPath)
}
