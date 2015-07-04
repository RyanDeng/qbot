package account

import (
	"encoding/json"
	"fmt"
	"net/http"

	"hd.qiniu.com/utils/object"
)

type Account struct {
	Uid      uint32 `json:"uid"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Gender   int    `json:"gender"`
}

type AccountRes struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

type AccountInfo struct {
	AccountRes
	Account Account `json:"data"`
}

func (i *AccountInfo) Serilize() (string, error) {
	return object.Serilize(i)
}

func (i *AccountInfo) Unserilize(str string) error {
	return object.Unserilize(str, i)
}

type AccountService interface {
	GetAccountInfo(accessToken string) (*AccountInfo, error)
	Signout(accessToken string) (*AccountRes, error)
}

type accountService struct {
	accountURL string
}

func NewAccountService(accountURL string) AccountService {
	return &accountService{
		accountURL: accountURL,
	}
}

func (i *accountService) GetAccountInfo(accessToken string) (info *AccountInfo, err error) {
	infoURL := i.accountURL + "/info?access_token=" + accessToken

	res, err := http.Get(infoURL)
	if err != nil {
		err = fmt.Errorf("get account info %s error: %s", infoURL, err)
		return
	}
	defer res.Body.Close()

	info = &AccountInfo{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(info)
	if err != nil {
		err = fmt.Errorf("decode account json %s error: %s", infoURL, err)
		return
	}

	return
}

func (i *accountService) Signout(accessToken string) (res *AccountRes, err error) {
	signoutURL := i.accountURL + "/signout?access_token=" + accessToken

	resp, err := http.Get(signoutURL)
	if err != nil {
		err = fmt.Errorf("signout %s error: %s", signoutURL, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("response code: %d", resp.StatusCode)
		return
	}
	res = &AccountRes{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(res)
	if err != nil {
		err = fmt.Errorf("decode signout json %s error: %s", signoutURL, err)
		return
	}

	return
}
