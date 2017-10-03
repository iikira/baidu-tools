package tieba

import (
	"github.com/iikira/baidu-tools"
)

// Tieba 百度贴吧账号详细情况
type Tieba struct {
	baidu.Baidu
	Bars []Bar //要执行任务的贴吧列表
}

//Bar 贴吧详情
type Bar struct {
	Fid, // 贴吧fid
	Name, // 名字
	Level string // 等级
	Exp int // 经验
}
