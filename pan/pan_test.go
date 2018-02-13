package pan

import (
	"fmt"
	"testing"
)

func TestPan(t *testing.T) {
	//https://pan.baidu.com/s/1o9oDpdo
	//https://pan.baidu.com/s/1c08q9Tu
	si, err := NewSharedInfo("https://pan.baidu.com/s/1o9oDpdo")
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
