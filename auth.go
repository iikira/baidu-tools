package baidu

// Auth 百度验证
type Auth struct {
	bduss, // 百度BDUSS
	ptoken,
	stoken string
}

func NewAuth(bduss, ptoken, stoken string) *Auth {
	return &Auth{
		bduss:  bduss,
		ptoken: ptoken,
		stoken: stoken,
	}
}

func (a *Auth) GetAuth() (bduss, ptoken, stoken string) {
	return a.bduss, a.ptoken, a.stoken
}

func (a *Auth) BDUSS() string {
	return a.bduss
}
