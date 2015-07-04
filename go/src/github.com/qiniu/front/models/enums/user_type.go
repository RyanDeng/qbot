package enums

const (
	// user type
	USER_TYPE_STDUSER  UserType = 0x0004
	USER_TYPE_STDUSER2          = 0x0008
)

type UserType uint32

// 是否标准用户
func (t UserType) IsStdUser() bool {
	return t&(USER_TYPE_STDUSER|USER_TYPE_STDUSER2) > 0
}
