package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"hd.qiniu.com/utils/object"
)

var (
	ErrTokenNotFound = errors.New("token is not found")
	ErrTokenExpired  = errors.New("token is expired")
	ErrTokenInvalid  = errors.New("token is invalid")
)

type OAuthToken struct {
	Error            int       `json:"error"`
	ErrorDescription string    `json:"error_description"`
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	ExpiresIn        int       `json:"expires_in"`
	ExpiresTime      time.Time `json:"-"`
}

func NewOAuthToken(r io.Reader) (token *OAuthToken, err error) {
	token = &OAuthToken{}
	decoder := json.NewDecoder(r)
	err = decoder.Decode(token)
	if err != nil {
		return
	}

	token.ExpiresTime = time.Now().Add(time.Second * time.Duration(token.ExpiresIn-60))
	return
}

func (t *OAuthToken) IsValid() bool {
	return t.AccessToken != "" &&
		t.RefreshToken != "" &&
		!t.IsExpired()
}

func (t *OAuthToken) IsExpired() bool {
	return time.Now().After(t.ExpiresTime)
}

func (t *OAuthToken) Serilize() (string, error) {
	return object.Serilize(t)
}

func (t *OAuthToken) Unserilize(str string) error {
	return object.Unserilize(str, t)
}

type OAuthService interface {
	AuthURL(redirectURL, state string) string
	Exchange(code string) (*OAuthToken, error)
	Refresh(refreshToken string) (*OAuthToken, error)
}

type oAuthService struct {
	authURL      string
	tokenURL     string
	clientId     string
	clientSecret string
}

func NewOAuthService(authURL, tokenURL, clientId, clientSecret string) OAuthService {
	return &oAuthService{
		authURL:      authURL,
		tokenURL:     tokenURL,
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func (s *oAuthService) AuthURL(redirectURL, state string) string {
	query := url.Values{}
	query.Set("client_id", s.clientId)
	query.Set("redirect_uri", redirectURL)
	query.Set("state", state)
	query.Set("response_type", "code")
	return s.authURL + "?" + query.Encode()
}

func (s *oAuthService) Exchange(code string) (token *OAuthToken, err error) {
	query := url.Values{}
	query.Set("grant_type", "authorization_code")
	query.Set("client_id", s.clientId)
	query.Set("client_secret", s.clientSecret)
	query.Set("code", code)
	tokenURL := s.tokenURL + "?" + query.Encode()

	res, err := http.Get(tokenURL)
	if err != nil {
		err = fmt.Errorf("exchange token get %s error: %s", tokenURL, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("response status: %d", res.StatusCode)
		return
	}
	token, err = NewOAuthToken(res.Body)
	return
}

func (s *oAuthService) Refresh(refreshToken string) (token *OAuthToken, err error) {
	query := &url.Values{}
	query.Set("grant_type", "refresh_token")
	query.Set("client_id", s.clientId)
	query.Set("client_secret", s.clientSecret)
	query.Set("refresh_token", refreshToken)
	refreshURL := s.tokenURL + "?" + query.Encode()

	res, err := http.Get(refreshURL)
	if err != nil {
		err = fmt.Errorf("refresh token get %s error: %s", refreshURL, err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = fmt.Errorf("response status: %d", res.StatusCode)
		return
	}
	token, err = NewOAuthToken(res.Body)
	return
}
