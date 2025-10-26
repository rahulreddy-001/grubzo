package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"grubzo/internal/router/auth/oauth"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Provider struct {
	Config oauth2.Config
	CB     string
	Type   string
	Name   string
	Icon   string
}

const (
	IconEndpoint    = "https://github.githubassets.com/assets/GitHub-Mark-ea2971cee799.png"
	ProfileEndpoint = "https://api.github.com/user"
	EmailEndpoint   = "https://api.github.com/user/emails"
)

func Init(cid, cs, cb string) *Provider {
	return &Provider{
		Config: oauth2.Config{
			ClientID:     cid,
			ClientSecret: cs,
			RedirectURL:  cb,
			Scopes:       []string{"user", "email"},
			Endpoint:     github.Endpoint,
		},
		Type: "github",
		Name: "Github",
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
		return nil, fmt.Errorf("failed to fetch user profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub responded with %d", resp.StatusCode)
	}

	var u struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, fmt.Errorf("failed to decode user profile: %w", err)
	}

	if u.Email == "" {
		email, err := p.getEmail(ctx, token)
		if err != nil {
			return nil, err
		}
		u.Email = email
	}

	return &oauth.OAuthUser{
		ID:    strconv.Itoa(u.ID),
		Name:  u.Name,
		Email: u.Email,
	}, nil
}

func (p *Provider) getEmail(ctx context.Context, token string) (string, error) {
	oauthToken := &oauth2.Token{AccessToken: token}
	client := p.Config.Client(ctx, oauthToken)

	resp, err := client.Get(EmailEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to fetch emails: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub email endpoint responded with %d", resp.StatusCode)
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	return "", errors.New("no verified primary email found for GitHub user")
}

func (p *Provider) ValidateToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oauthToken := &oauth2.Token{AccessToken: token}
	client := p.Config.Client(ctx, oauthToken)

	resp, err := client.Get(ProfileEndpoint)
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid token: GitHub responded with %d", resp.StatusCode)
	}
	return nil
}
