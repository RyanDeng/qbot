package account

import "hd.qiniu.com/services/account"

func AdminService(host, clientId, clientSecret, username, password string) account.AdminService {
	return account.NewAdminService(host, clientId, clientSecret, username, password)
}

func AccountService(accountURL string) interface{} {
	return func() account.AccountService {
		return account.NewAccountService(accountURL)
	}
}

func OAuthService(authURL, tokenURL, clientId, clientSecret string) interface{} {
	return func() account.OAuthService {
		return account.NewOAuthService(authURL, tokenURL, clientId, clientSecret)
	}
}
