package fbauth

const UrlFBMe = "https://graph.facebook.com/me"

type Config struct {
	FBlientID   string `json:"FB_CLIENT_ID"`
	FBSecretKey string `json:"FB_SECRET_KEY"`
}

type ProfileData struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}
