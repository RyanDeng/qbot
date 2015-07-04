package biz

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"labix.org/v2/mgo/bson"
)

type BizService interface {
	GetIdentity(uid uint32) (*IdentityInfo, error)
}

type bizService struct {
	bizHost string
	client  *http.Client
}

func NewBizService(host string, tr http.RoundTripper) BizService {
	return &bizService{
		bizHost: host,
		client: &http.Client{
			Transport: tr,
		},
	}
}

type IdentityInfo struct {
	Id                    string        `json:"id"`
	Uid                   uint32        `json:"uid"`
	EnterpriseName        string        `json:"enterprise_name"`          // 名称
	BusinessLicenseNo     string        `json:"business_license_no"`      // 营业执照注册号
	OrganizationNo        string        `json:"organization_no"`          // 组织机构代码
	BusinessLicenseCopy   string        `json:"business_license_copy"`    // 营业执照副本扫描件
	ContactName           string        `json:"contact_name"`             // 联系人姓名
	ContactIdentityNo     string        `json:"contact_identity_no"`      // 联系人身份证号码
	ContactIdentityPhoto  string        `json:"contact_identity_photo"`   // 联系人身份证持证照片
	ContactIdentityPhotoB string        `json:"contact_identity_photo_b"` // 联系人身份证持证背面照片
	ContactAddress        string        `json:"contact_address"`          // 联系地址
	ContactProvince       string        `json:"contact_province"`         // 所在省
	ContactCity           string        `json:"contact_city"`             // 所在市
	ContactRegion         string        `json:"contact_region"`           // 所在区
	Status                int           `json:"status"`                   // 状态
	StatusNote            string        `json:"status_note"`              // 状态信息
	IsEnterprise          bool          `json:"is_enterprise"`            // 是企业认证
	Memo                  string        `json:"memo"`                     // 备忘
	CreatorId             bson.ObjectId `json:"creator_id,omitempty"`     // 创建者
	OperatorId            bson.ObjectId `json:"operator_id,omitempty"`    // 创建者
	CreatedAt             time.Time     `json:"created_at"`               // 创建时间
	UpdatedAt             time.Time     `json:"updated_at"`               // 更新时间
}

func (s *bizService) GetIdentity(uid uint32) (info *IdentityInfo, err error) {
	resp, err := s.client.PostForm(s.bizHost+"/admin/developer/identity", url.Values{
		"uid": {strconv.FormatUint(uint64(uid), 10)},
	})
	if err != nil {
		return
	}
	defer resp.Body.Close()

	info = &IdentityInfo{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(info)
	if err != nil {
		return
	}

	return
}
