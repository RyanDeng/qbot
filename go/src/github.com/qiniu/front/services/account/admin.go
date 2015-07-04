package account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"qbox.us/oauth"
)

type AdminService interface {
	Transport() *oauth.Transport
	Auth() error
	GetAccInfo(uid uint32) (*AccInfo, error)
}

type adminService struct {
	accountHost  string
	clientId     string
	clientSecret string
	username     string
	password     string
	transport    *oauth.Transport
	client       *http.Client
}

func NewAdminService(host, clientId, clientSecret, username, password string) AdminService {
	return &adminService{
		accountHost:  host,
		clientId:     clientId,
		clientSecret: clientSecret,
		username:     username,
		password:     password,
	}
}

func (s *adminService) Client() *http.Client {
	if s.client == nil {
		s.client = &http.Client{
			Transport: s.Transport(),
		}
	}
	return s.client
}

func (s *adminService) Transport() *oauth.Transport {
	if s.transport == nil {
		s.transport = &oauth.Transport{
			Config: &oauth.Config{
				ClientId:     s.clientId,
				ClientSecret: s.clientSecret,
				Scope:        "Scope",
				AuthURL:      "<AuthURL>",
				TokenURL:     s.accountHost + "/oauth2/token",
				RedirectURL:  "<RedirectURL>",
			},
			Transport: http.DefaultTransport, // it is default
		}
	}
	return s.transport
}

func (s *adminService) Auth() error {
	tr := s.Transport()
	_, _, err := tr.ExchangeByPassword(s.username, s.password)
	return err
}

func (s *adminService) GetAccInfo(uid uint32) (info *AccInfo, err error) {
	client := s.Client()
	resp, err := client.PostForm(s.accountHost+"/admin/user/info", url.Values{
		"uid": {strconv.FormatUint(uint64(uid), 10)},
	})
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("response status: %d", resp.StatusCode)
		return
	}
	info = &AccInfo{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(info)
	return
}
