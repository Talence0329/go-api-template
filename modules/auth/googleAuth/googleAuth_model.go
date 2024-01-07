package googleauth

// https://developers.google.com/identity/protocols/oauth2/scopes?hl=zh-tw
const ScopeEmail = "userinfo.email"
const ScopeProfile = "userinfo.profile"

const UrlLogin = "https://accounts.google.com/o/oauth2/v2/auth"
const UrlToken = "https://www.googleapis.com/oauth2/v4/token"
const UrlScope = "https://www.googleapis.com/auth"
const UrlUserInfo = "https://www.googleapis.com/oauth2/v1/userinfo"

const ResponseTypeCode = "code"

const GrantTypeAuthCode = "authorization_code"

type Config struct {
	ClientID  string `json:"GOOLE_CLIENT_ID"`
	SecretKey string `json:"GOOGLE_SECRET_KEY"`
}
type ScopeProfileData struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type ScopeEmailData struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}
