package baidu

// Baidu 百度帐号详细情况
type Baidu struct {
	UID, // 百度ID对应的uid
	Name, // 真实ID
	NameShow, // 显示的用户名(昵称)
	Sex, // 性别
	Age string // 帐号年龄
	Auth *Auth
}
