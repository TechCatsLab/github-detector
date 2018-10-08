package util

import "regexp"

func IsValidPhone(phone string) bool {
	reg := `^1(3[0-9]|4[579]|5[^4]|7[0135678]|8[0-9])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)
}
