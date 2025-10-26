package google

import (
	"context"
	"encoding/json"
	"fmt"
	"grubzo/internal/router/auth/oauth"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Provider struct {
	Config oauth2.Config
	CB     string
	Type   string
	Name   string
	Icon   string
}

const (
	IconEndpoint    = "https://www.gstatic.com/marketing-cms/assets/images/d5/dc/cfe9ce8b4425b410b49b7f2dd3f3/g.webp"
	ProfileEndpoint = "https://openidconnect.googleapis.com/v1/userinfo"
	TokenInfoURL    = "https://www.googleapis.com/oauth2/v3/tokeninfo?access_token="
)

func Init(cid, cs, cb string) *Provider {
	return &Provider{
		Config: oauth2.Config{
			ClientID:     cid,
			ClientSecret: cs,
			RedirectURL:  cb,
			Scopes:       []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		},
		Type: "google",
		Name: "Google",
		Icon: IconEndpoint,
		CB:   cb,
	}
}

func (p *Provider) GetType() string           { return p.Type }
func (p *Provider) GetName() string           { return p.Name }
func (p *Provider) GetIcon() string           { return p.Icon }
func (p *Provider) GetConfig() *oauth2.Config { return &p.Config }
func (p *Provider) GetCallbackURL() string    { return p.CB }

func (p *Provider) FetchUser(token string) (*oauth.OAuthUser, error) {
	ctx := context.Background()
	oauthToken := &oauth2.Token{AccessToken: token}
	client := p.Config.Client(ctx, oauthToken)

	resp, err := client.Get(ProfileEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s responded with %d", p.Name, resp.StatusCode)
	}

	var u struct {
		ID    string `json:"id"`
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	id := u.ID
	if id == "" {
		id = u.Sub
	}

	return &oauth.OAuthUser{
		ID:    id,
		Name:  u.Name,
		Email: u.Email,
	}, nil
}

func (p *Provider) ValidateToken(token string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(TokenInfoURL + token)
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid token: %s", resp.Status)
	}
	return nil
}
