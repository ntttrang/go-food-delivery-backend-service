package shareComponent

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOauth struct {
	oauthCfg *oauth2.Config
}

func NewGoogleOauth(cfg datatype.GoogleConfig) *GoogleOauth {
	return &GoogleOauth{
		oauthCfg: &oauth2.Config{
			RedirectURL:  cfg.RedirectUrl,
			ClientID:     cfg.ClientId,     //"YOUR_GOOGLE_CLIENT_ID",
			ClientSecret: cfg.ClientSecret, //"YOUR_GOOGLE_CLIENT_SECRET",
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

var state string

func (g *GoogleOauth) GenerateState(ctx context.Context) string {
	b := make([]byte, 16)
	rand.Read(b)
	state = base64.URLEncoding.EncodeToString(b)
	//state := "fooddelivery2025"
	return state
}

func (g *GoogleOauth) AuthCodeUrl(ctx context.Context, state string) string {
	url := g.oauthCfg.AuthCodeURL(state)
	return url
}

func (g *GoogleOauth) GetGGUserInfo(ctx context.Context, stateReq string, code string) (*datatype.GgUserInfo, error) {
	// if state != stateReq {
	// 	return nil, errors.New("State mismatch")
	// }

	token, err := g.oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, errors.New("Code exchange failed")
	}

	client := g.oauthCfg.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, errors.New("Failed to get user info")
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, errors.New("Failed to decode user info")
	}

	var ggUserInfo datatype.GgUserInfo
	ggUserInfo.Email = userInfo["email"].(string)
	ggUserInfo.GgId = userInfo["id"].(string)
	ggUserInfo.Name = userInfo["name"].(string)
	ggUserInfo.GivenName = userInfo["given_name"].(string)
	ggUserInfo.Avatar = userInfo["picture"].(string)
	ggUserInfo.VerifyEmail = userInfo["verified_email"].(bool)

	return &ggUserInfo, nil

}
