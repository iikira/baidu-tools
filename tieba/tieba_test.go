package tieba

import (
	"fmt"
	"testing"
)

func TestUserInfo(t *testing.T) {
	info, _ := NewUserInfoByUID(2)
	fmt.Println(info.Baidu)
}
