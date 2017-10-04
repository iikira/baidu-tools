package baidu

// Auth 百度验证
type Auth struct {
	BDUSS, // 百度BDUSS
	PTOKEN,
	STOKEN string
}

func NewAuth(bduss, ptoken, stoken string) *Auth {
	return &Auth{
		BDUSS:  bduss,
		PTOKEN: ptoken,
		STOKEN: stoken,
	}
}
