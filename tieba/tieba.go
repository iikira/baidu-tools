package tieba

import (
	"github.com/iikira/baidu-tools"
)

// Tieba 百度贴吧账号详细情况
type Tieba struct {
	Baidu *baidu.Baidu
	Tbs   string
	Stat  *Stat
	Bars  []Bar //要执行任务的贴吧列表
}

//Bar 贴吧详情
type Bar struct {
	Fid, // 贴吧fid
	Name, // 名字
	Level string // 个人等级
	Exp int // 个人经验
}

// Stat 统计数据
type Stat struct {
	LikeForumNum int
	PostNum      int
}
