package auth

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	sessionKeyAuthenticated     = "isAuthenticated"
	sessionKeyPreferredUsername = "preferredUsername"
	routeLogin                  = "/login"
	routeLogout                 = "/logout"
	routeProfile                = "/profile"
)

type SessionProvider struct {
	config     *ProviderConfig
	cookieName string
	cookiePath string
	validator  Validator
	store      cookie.Store
}

func NewSessionProvider(cookieName string, cookiePath string, validator Validator, store sessions.Store) Provider {
	return &SessionProvider{config: &ProviderConfig{
		HasLoginRoute:   true,
		HasLogoutRoute:  true,
		HasProfileRoute: true,
		LoginPath:       routeLogin,
		LogoutPath:      routeLogout,
		ProfilePath:     routeProfile,
		LoginMethod:     http.MethodPost,
		LogoutMethod:    http.MethodGet,
		ProfileMethod:   http.MethodGet,
	}, cookieName: cookieName, cookiePath: cookiePath, store: store, validator: validator}
}

func (p *SessionProvider) Config() *ProviderConfig {
	return p.config
}

func (p *SessionProvider) PublicMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{p.sessionMiddleware()}
}

func (p *SessionProvider) ProtectedMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{p.sessionMiddleware()}
}

func (p *SessionProvider) IsAuthenticated(c *gin.Context, passThroughFunc ...ProviderPassThroughFunc) bool {
	for _, passFn := range passThroughFunc {
		if passFn(c) {
			return true
		}
	}

	session := sessions.Default(c)
	sessionIsAuthenticated := session.Get(sessionKeyAuthenticated)
	sessionPreferredUsername := session.Get(sessionKeyPreferredUsername)

	if sessionIsAuthenticated == nil || sessionPreferredUsername == nil {
		return false
	}

	return true
}

func (p *SessionProvider) Login(c *gin.Context) (error, int) {
	var req *UserCredentials

	if err := c.ShouldBindJSON(&req); err != nil {
		return errCredentialsBind, http.StatusBadRequest
	}
	if ok := p.validator.Validate(req); !ok {
		return errForbidden, http.StatusForbidden
	}

	session := sessions.Default(c)
	session.Set(sessionKeyAuthenticated, true)
	session.Set(sessionKeyPreferredUsername, req.Username)

	if err := session.Save(); err != nil {
		return fmt.Errorf("could not create session for '%s': %w", session, err), http.StatusInternalServerError
	}

	return nil, http.StatusNoContent
}

func (p *SessionProvider) Logout(c *gin.Context) (error, int) {
	session := sessions.Default(c)
	// marks the session as "written"
	session.Set(sessionKeyAuthenticated, false)
	session.Set(sessionKeyPreferredUsername, nil)
	session.Options(sessions.Options{Path: p.cookiePath, MaxAge: -1})

	if err := session.Save(); err != nil {
		return fmt.Errorf("could not logout session %w", err), http.StatusInternalServerError
	}
	session.Clear()

	return nil, http.StatusNoContent
}

func (p *SessionProvider) Profile(c *gin.Context) (*Profile, error, int) {
	session := sessions.Default(c)
	sessionPreferredUsername := session.Get(sessionKeyPreferredUsername)

	if sessionPreferredUsername == nil {
		return nil, errProfileEmpty, http.StatusInternalServerError
	}

	return &Profile{PreferredUsername: sessionPreferredUsername.(string)}, nil, http.StatusOK
}

func (p *SessionProvider) RouteContext(c *gin.Context) (*Profile, error) {
	session := sessions.Default(c)
	sessionPreferredUsername := session.Get(sessionKeyPreferredUsername)

	if sessionPreferredUsername == nil {
		return nil, errProfileEmpty
	}

	return &Profile{PreferredUsername: sessionPreferredUsername.(string)}, nil
}

func (p *SessionProvider) sessionMiddleware() gin.HandlerFunc {
	return sessions.Sessions(p.cookieName, p.store)
}
