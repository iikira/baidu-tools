package pan

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/requester"
	"testing"
)

func TestPan(t *testing.T) {
	//https://pan.baidu.com/s/1o9oDpdo
	//https://pan.baidu.com/s/1c08q9Tu
	//链接:https://pan.baidu.com/s/1djChHW 密码:ywsp
	si := NewSharedInfo("https://pan.baidu.com/s/1djChHW")
	err := si.Auth("ywsp")
	if err != nil {
		t.Log(err)
		return
	}

	err = si.InitInfo()
	if err != nil {
		t.Log(err)
		return
	}

	fileInfo, err := si.Meta("/567/23.txt")
	if err != nil {
		t.Log(err)
		return
	}

	fmt.Println(fileInfo.Dlink)
}

func TestPan2(t *testing.T) {
	si := NewSharedInfo("https://pan.baidu.com/s/1QC6obCSrR5_KoE3rvtRB2A")
	si.Client = requester.NewHTTPClient()
	si.Client.SetHTTPSecure(false)

	err := si.InitInfo()
	if err != nil {
		t.Log(err)
		return
	}

	fileInfo, err := si.Meta("randomdata1.txt")
	if err != nil {
		t.Log(err)
		return
	}

	fmt.Println(fileInfo.Dlink)
}
