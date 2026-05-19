package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const (
	RouteContext = "ROUTE_CONTEXT"
)

var (
	errForbidden       = errors.New("access is prohibited")
	errCredentialsBind = errors.New("could not bind credentials")
	errProfileEmpty    = errors.New("profile is empty")
)

type ProviderPassThroughFunc func(c *gin.Context) bool

type Provider interface {
	// Config hosts the provider's configuration used to apply provider to gin.Engine
	Config() *ProviderConfig

	// PublicMiddleware middleware which is added for this Provider for public paths
	PublicMiddleware() []gin.HandlerFunc

	// ProtectedMiddleware middleware which is added for this Provider for protected paths
	ProtectedMiddleware() []gin.HandlerFunc

	// IsAuthenticated defines a function which is added to protected paths, MUST evaluate and respect result of provided ProviderPassThroughFunc argument
	IsAuthenticated(c *gin.Context, passThroughFunc ...ProviderPassThroughFunc) bool

	// Login defines a function to log in and returns the HTTP status code
	Login(c *gin.Context) (error, int)

	// Logout defines a function to log out and returns the HTTP status code
	Logout(c *gin.Context) (error, int)

	// Profile defines a function to fetch profile and returns the HTTP status code
	Profile(c *gin.Context) (*Profile, error, int)

	// RouteContext defines a function to get profile to pass to subsequent handlers which is added for this Provider for protected paths
	RouteContext(c *gin.Context) (*Profile, error)
}

type ProviderConfig struct {
	HasLoginRoute   bool
	HasLogoutRoute  bool
	HasProfileRoute bool

	LoginPath   string
	LogoutPath  string
	ProfilePath string

	LoginMethod   string
	LogoutMethod  string
	ProfileMethod string
}

type Profile struct {
	PreferredUsername string `json:"preferredUsername"`
}
