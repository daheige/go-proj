package helper

import "regexp"

//根据ua获取设备名称
func GetDeviceByUa(ua string) string {
	plat := "web"
	regText := "(a|A)ndroid|dr"
	re := regexp.MustCompile(regText)
	if re.MatchString(ua) {
		plat = "android"
	} else {
		regText = "i(p|P)(hone|ad|od)|(m|M)ac"
		re = regexp.MustCompile(regText)
		if re.MatchString(ua) {
			plat = "ios"
		}
	}

	return plat
}
