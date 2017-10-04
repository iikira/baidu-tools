package baidu

// NewUser 返回 Baidu 指针
func NewUser(uid uint64, name string) *Baidu {
	return &Baidu{
		UID:  uid,
		Name: name,
	}
}
