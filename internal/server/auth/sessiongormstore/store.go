package sessiongormstore

import (
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/wader/gormstore/v2"
	"gorm.io/gorm"
)

type GormStore interface {
	sessions.Store
}

type gormStore struct {
	*gormstore.Store
}

type GormCleanupFunc func(s *gormstore.Store)

type Options struct {
	gormStoreOptions gormstore.Options
	sessionOptions   sessions.Options
	KeyPairs         [][]byte
	cleanupEnabled   bool
	cleanupInterval  time.Duration
	cleanupFunc      *GormCleanupFunc
}

// New constructs a new store
func New(db *gorm.DB, options ...Option) GormStore {
	opts := &Options{gormStoreOptions: gormstore.Options{},
		KeyPairs:       make([][]byte, 0),
		cleanupEnabled: false,
	}

	if options != nil {
		for _, o := range options {
			o.Apply(opts)
		}
	}

	store := gormstore.NewOptions(db, opts.gormStoreOptions, opts.KeyPairs...)

	if opts.cleanupEnabled && opts.cleanupFunc == nil {
		quit := make(chan struct{})
		go store.PeriodicCleanup(opts.cleanupInterval, quit)
	}

	if opts.cleanupEnabled && opts.cleanupFunc != nil {
		run := *opts.cleanupFunc
		run(store)
	}

	g := &gormStore{store}
	g.Options(opts.sessionOptions)
	return g
}

func (s *gormStore) Options(options sessions.Options) {
	s.SessionOpts = options.ToGorillaOptions()
}

type Option interface {
	Apply(o *Options)
}

type OptionFunc func(o *Options)

func (f OptionFunc) Apply(o *Options) {
	f(o)
}

// WithGormStoreOpts exposes option to configure underlying gormstore
func WithGormStoreOpts(options gormstore.Options) Option {
	return OptionFunc(func(o *Options) {
		o.gormStoreOptions = options
	})
}

// WithSessionOpts exposes option to configure Gorilla session options
func WithSessionOpts(options sessions.Options) Option {
	return OptionFunc(func(o *Options) {
		o.sessionOptions = options
	})
}

// WithKeyPairs exposes option to configure key pairs
func WithKeyPairs(keyPairs ...[]byte) Option {
	return OptionFunc(func(o *Options) {
		o.KeyPairs = keyPairs
	})
}

// WithCleanup exposes option to configure periodic session cleanup with a simple go routine
func WithCleanup(d time.Duration) Option {
	return OptionFunc(func(o *Options) {
		o.cleanupEnabled = true
		o.cleanupInterval = d
	})
}

// WithCleanupFunc exposes option to configure cleanup with a function allowing fine-granular control
func WithCleanupFunc(f GormCleanupFunc) Option {
	return OptionFunc(func(o *Options) {
		o.cleanupEnabled = true
		o.cleanupFunc = &f
	})
}
