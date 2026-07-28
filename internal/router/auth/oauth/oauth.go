package oauth

import (
	"context"
	"errors"
	"fmt"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"grubzo/internal/repository"
	"grubzo/internal/router/ext"
	"grubzo/internal/router/session"
	"grubzo/internal/utils/random"
	"grubzo/internal/utils/tenantutils"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type ProviderLoginInfo struct {
	IsRedirect bool   `json:"is_redirect"`
	Icon       string `json:"icon"`
	URL        string `json:"url"`
	Name       string `json:"name"`
}

type Provider interface {
	GetConfig() *oauth2.Config
	GetType() string
	GetName() string
	GetIcon() string
	GetCallbackURL() string
	FetchUser(context.Context, string) (*OAuthUser, error)
	ValidateToken(context.Context, string) error
}

type OAuthUser struct {
	ID    string
	Email string
	Name  string
}

type Auth struct {
	providers    map[string]Provider
	domain       string
	env          string
	instance     string
	router       *gin.RouterGroup
	sessionStore session.Store
	repo         *repository.Repository
	logger       *zap.Logger
}

func New() *Auth {
	return &Auth{
		providers: make(map[string]Provider, 0),
	}
}
func (a *Auth) SetProviders(providers ...Provider) *Auth {
	for _, provider := range providers {
		a.providers[provider.GetType()] = provider
	}
	return a
}
func (a *Auth) UseRouter(r *gin.RouterGroup) *Auth {
	a.router = r
	return a
}
func (a *Auth) WithSessionStore(store session.Store) *Auth {
	a.sessionStore = store
	return a
}
func (a *Auth) WithRepository(repo *repository.Repository) *Auth {
	a.repo = repo
	return a
}
func (a *Auth) WithDomain(domain string) *Auth {
	a.domain = domain
	return a
}
func (a *Auth) WithEnv(env string) *Auth {
	a.env = env
	return a
}
func (a *Auth) WithInstance(instance string) *Auth {
	a.instance = instance
	return a
}
func (a *Auth) WithLogger(logger *zap.Logger) *Auth {
	a.logger = logger.Named("oauth")
	return a
}

func (a *Auth) Init() *Auth {
	for _, p := range a.providers {
		provider := p
		a.router.GET(fmt.Sprintf("/login/%s", provider.GetType()), func(ctx *gin.Context) {
			subDomain, ok := tenantutils.SubDomainFromHost(ctx.Request.Host, a.domain, a.env, a.instance)
			if !ok {
				ext.Ctx(ctx).RespondWithError(ext.Error("tenant subdomain is required"))
				return
			}
			state := fmt.Sprintf("%s.%s", random.SecureAlphaNumeric(50), subDomain)
			ctx.SetCookie("oauth_state", state, 300, "/", "", false, true)
			ctx.Redirect(http.StatusPermanentRedirect, provider.GetConfig().AuthCodeURL(state))
		})

		cbURL, _ := url.Parse(provider.GetCallbackURL())
		cleanCBPath := strings.TrimPrefix(cbURL.Path, a.router.BasePath())
		a.router.GET(cleanCBPath, func(c *gin.Context) {
			token, err := a.Exchange(provider, c)
			if err != nil {
				ext.Ctx(c).RespondWithError(fmt.Errorf("login failed: %w", err))
				return
			}
			if err := provider.ValidateToken(c.Request.Context(), token.AccessToken); err != nil {
				ext.Ctx(c).RespondWithError(fmt.Errorf("login failed: %w", err))
				return
			}
			user, err := provider.FetchUser(c.Request.Context(), token.AccessToken)
			if err != nil {
				ext.Ctx(c).RespondWithError(fmt.Errorf("login failed: %w", err))
				return
			}
			tenant, err := a.tenantFromOAuthState(c)
			if err != nil {
				ext.Ctx(c).RespondWithError(fmt.Errorf("login failed: %w", err))
				return
			}
			userEntity, err := a.repo.FindUser(c.Request.Context(), query.NewUserQuery(tenant.ID).WithEmail(user.Email))
			if err != nil {
				if err.Error() == repository.UserNotFound {
					userEntity, err = a.repo.CreateUser(c.Request.Context(), &dto.CreateUser{
						TenantID: tenant.ID,
						UserID:   user.ID,
						Email:    user.Email,
						Name:     user.Name,
					})
					if err != nil {
						ext.Ctx(c).RespondWithError(fmt.Errorf("login failed: %w", err))
						return
					}
				} else {
					ext.Ctx(c).RespondWithError(fmt.Errorf("login failed: %w", err))
					return
				}
			}
			userSession, err := a.sessionStore.RenewSession(c, userEntity.ID, userEntity.TenantID)
			if err != nil {
				ext.Ctx(c).RespondWithError(fmt.Errorf("login failed: %w", err))
				return
			}
			trueRef := true
			locationEntity, _ := a.repo.FindTenantLocation(c, &query.TenantLocationQuery{
				TenantID:  userEntity.TenantID,
				IsPrimary: &trueRef,
			})
			if err := userSession.Set("user", &session.UserSession{
				TenantID: userEntity.TenantID,
				UserID:   userEntity.ID,
				Email:    userEntity.Email,
				Type:     "user",
				Location: locationEntity.ID,
			}); err != nil {
				ext.Ctx(c).RespondWithError(fmt.Errorf("login failed: %w", err))
				return
			}
			a.RedirectToLoginSuccessPage(c)
		})
	}
	return a
}

func (a *Auth) GetLoginData() []ProviderLoginInfo {
	var loginInfo []ProviderLoginInfo
	for _, provider := range a.providers {
		providerLoginUrl := fmt.Sprintf("%s/login/%s", a.router.BasePath(), provider.GetType())
		loginInfo = append(loginInfo, ProviderLoginInfo{
			IsRedirect: true,
			Icon:       provider.GetIcon(),
			URL:        providerLoginUrl,
			Name:       provider.GetName(),
		})
	}
	return loginInfo
}

func (a *Auth) Exchange(provider Provider, ctx *gin.Context) (*oauth2.Token, error) {
	state, code := ctx.Query("state"), ctx.Query("code")
	if false { // skipping state verification for now
		cookieState, err := ctx.Cookie("oauth_state")
		if err != nil {
			return nil, errors.New("missing state cookie: " + err.Error())
		}
		if state != cookieState {
			return nil, errors.New("invalid state")
		}
	}
	ctx.SetCookie("oauth_state", state, -1, "/", "", false, true)

	providerConfig := provider.GetConfig()
	token, err := providerConfig.Exchange(ctx.Request.Context(), code)
	if err != nil {
		return nil, errors.New("exchange Failed")
	}
	return token, nil
}

func (a *Auth) RedirectToLoginPage(ctx *gin.Context) {
	_ = a.sessionStore.RevokeSession(ctx)
	ctx.Redirect(http.StatusTemporaryRedirect, "/login")
}

func (a *Auth) RedirectToLoginSuccessPage(ctx *gin.Context) {
	subDomain := "www"
	if tenant, err := a.tenantFromOAuthState(ctx); err == nil {
		subDomain = tenant.SubDomain
	}
	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("http://%s", tenantutils.HostForSubDomain(subDomain, a.domain, a.env)))
}

func (a *Auth) tenantFromOAuthState(ctx *gin.Context) (*entity.Tenant, error) {
	parts := strings.Split(ctx.Query("state"), ".")
	if len(parts) < 2 {
		return nil, errors.New("missing tenant state")
	}
	subDomain := parts[len(parts)-1]
	return a.repo.GetTenant(ctx.Request.Context(), query.NewTenantQuery().WithSubDomain(subDomain))
}
