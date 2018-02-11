package pan

import (
	"regexp"
)

var (
	YunDataExp = regexp.MustCompile(`window\.yunData[\s]?=[\s]?(.*?);`)
)
