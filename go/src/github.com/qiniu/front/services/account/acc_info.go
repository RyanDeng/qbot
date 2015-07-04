package account

import (
	"fmt"
	"time"
)

type CustomerGroup int

const (
	CUSTOMER_GROUP_EXP     CustomerGroup = 0
	CUSTOMER_GROUP_NORMAL  CustomerGroup = 1
	CUSTOMER_GROUP_VIP     CustomerGroup = 2
	CUSTOMER_GROUP_INVALID CustomerGroup = 3
)

const (
	USER_TYPE_ADMIN      = 0x0001
	USER_TYPE_VIP        = 0x0002
	USER_TYPE_STDUSER    = 0x0004
	USER_TYPE_STDUSER2   = 0x0008
	USER_TYPE_EXPUSER    = 0x0010
	USER_TYPE_PARENTUSER = 0x0020
	USER_TYPE_OP         = 0x0040
	USER_TYPE_SUPPORT    = 0x0080
	USER_TYPE_CC         = 0x0100
	USER_TYPE_DISABLED   = 0x8000

	USER_TYPE_USERS   = USER_TYPE_STDUSER | USER_TYPE_STDUSER2 | USER_TYPE_EXPUSER
	USER_TYPE_SUDOERS = USER_TYPE_ADMIN | USER_TYPE_OP | USER_TYPE_SUPPORT
)

func (cg CustomerGroup) Humanize() string {
	switch cg {
	case CUSTOMER_GROUP_EXP:
		return "体验用户"
	case CUSTOMER_GROUP_NORMAL:
		return "标准用户"
	case CUSTOMER_GROUP_VIP:
		return "高级用户"
	case CUSTOMER_GROUP_INVALID:
		return "无效用户"
	default:
		return fmt.Sprintf("未知用户类型: %d", cg)
	}
}

type AccInfo struct {
	Id               string    `json:"id"`              // 用户名(UserName)。唯一。
	Email            string    `json:"email"`           // 电子邮箱。唯一。
	CreatedAt        int64     `json:"ctime"`           // 用户创建时间。
	UpdatedAt        int64     `json:"etime"`           // 最后一次修改时间。
	LastLoginAt      int64     `json:"lgtime"`          // 最后一次登录时间。
	Uid              uint32    `json:"uid"`             // 用户数字ID。唯一。
	Utype            uint32    `json:"utype"`           // 用户类型。
	ParentUid        uint32    `json:"parent_uid"`      // 父用户Uid
	Activated        bool      `json:"activated"`       // 用户是否已经激活。
	DisabledType     int       `json:"disabled_type"`   // 用户冻结类型
	DisabledReason   string    `json:"disabled_reason"` // 用户冻结原因
	DisabledAt       time.Time `json:"disabled_at"`     // 用户冻结时间
	ChildEmailDomain string    `json:"child_email_domain"`
	CanGetChildKey   bool      `json:"can_get_child_key"`
}

func (i *AccInfo) IsExpUser() bool {
	return i.GetCustomerGroup() == CUSTOMER_GROUP_EXP
}
func (i *AccInfo) IsStdUser() bool {
	return i.GetCustomerGroup() == CUSTOMER_GROUP_NORMAL
}

func (i *AccInfo) IsVipUser() bool {
	return i.GetCustomerGroup() == CUSTOMER_GROUP_VIP
}

func (i *AccInfo) IsInvalid() bool {
	return i.GetCustomerGroup() == CUSTOMER_GROUP_INVALID
}

func (i *AccInfo) GetCustomerGroup() CustomerGroup {
	if i.Utype&USER_TYPE_USERS == 0 {
		return CUSTOMER_GROUP_INVALID
	}
	if i.Utype&USER_TYPE_EXPUSER != 0 {
		return CUSTOMER_GROUP_EXP
	}
	if i.Utype&USER_TYPE_VIP != 0 {
		return CUSTOMER_GROUP_VIP
	}
	return CUSTOMER_GROUP_NORMAL
}
