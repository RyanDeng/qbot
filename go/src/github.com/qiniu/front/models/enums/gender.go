package enums

import (
	"strconv"
)

type Gender int

const (
	GenderMale   Gender = 0
	GenderFemale Gender = 1
)

func (g Gender) String() string {
	return strconv.Itoa(int(g))
}

func (g Gender) Humanize() string {
	if g == GenderMale {
		return "男"
	}
	return "女"
}

func (g Gender) Salutation() string {
	if g == GenderMale {
		return "先生"
	}
	return "女士"
}
