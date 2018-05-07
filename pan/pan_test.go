package pan

import (
	"fmt"
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

	dlink, err := si.GetDownloadLink("/567/23.txt")
	if err != nil {
		t.Log(err)
		return
	}

	fmt.Println(dlink)
}
