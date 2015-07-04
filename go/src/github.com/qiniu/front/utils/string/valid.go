package string

import "regexp"

var eduSuffixReg = regexp.MustCompile(`^.+@(.+\.)?edu(\.[^.]+)?$`)

func IsEduEmail(email string) bool {
	return eduSuffixReg.MatchString(email)
}
