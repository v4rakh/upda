package server

import (
	"git.myservermanager.com/varakh/upda/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"os"
	"time"
)

func (u *Update) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

// Update entity holding information for updates
type Update struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;unique;not null"`
	Application string    `gorm:"uniqueIndex:idx_a_p_h;not null"`
	Provider    string    `gorm:"uniqueIndex:idx_a_p_h;not null"`
	Host        string    `gorm:"uniqueIndex:idx_a_p_h;not null"`
	Version     string    `gorm:"not null"`
	State       string    `gorm:"not null"`
	Metadata    JSONMap   `gorm:"jsonb"`
	CreatedAt   time.Time `gorm:"time;autoCreateTime;not null"`
	UpdatedAt   time.Time `gorm:"time;autoUpdateTime;not null"`
}

// BeforeCreate encrypts secret value before storing to database
func (wh *Webhook) BeforeCreate(tx *gorm.DB) (err error) {
	var er error
	var encryptedToken string

	if encryptedToken, er = util.EncryptAndEncode(wh.Token, os.Getenv(envSecret)); er != nil {
		return er
	}

	wh.ID = uuid.New()
	wh.Token = encryptedToken
	return
}

// AfterSave decrypt secret value after encrypted value has been retrieved from database
func (wh *Webhook) AfterSave(tx *gorm.DB) (err error) {
	var er error
	var decrypted string
	if decrypted, er = util.DecryptAndDecode(wh.Token, os.Getenv(envSecret)); er != nil {
		return er
	}

	wh.Token = decrypted
	return
}

// Webhook entity holding information for webhooks
type Webhook struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;unique;not null"`
	Type       string    `gorm:"not null"`
	Label      string    `gorm:"not null"`
	Token      string    `gorm:"not null"`
	IgnoreHost bool      `gorm:"default:false;not null"`
	CreatedAt  time.Time `gorm:"time;autoCreateTime;not null"`
	UpdatedAt  time.Time `gorm:"time;autoUpdateTime;not null"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	return
}

// Event entity holding information for events
type Event struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;unique;not null"`
	Name      string    `gorm:"not null"`
	State     string    `gorm:"not null"`
	Payload   JSONMap   `gorm:"jsonb"`
	CreatedAt time.Time `gorm:"time;autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"time;autoUpdateTime;not null"`
}

// BeforeCreate encrypts secret value before storing to database
func (e *Secret) BeforeCreate(tx *gorm.DB) (err error) {
	var er error
	var encryptedValue string

	if encryptedValue, er = util.EncryptAndEncode(e.Value, os.Getenv(envSecret)); er != nil {
		return er
	}

	e.ID = uuid.New()
	e.Value = encryptedValue
	return
}

// BeforeUpdate encrypts secret value before storing to database
func (e *Secret) BeforeUpdate(tx *gorm.DB) (err error) {
	var er error
	var encryptedValue string

	if encryptedValue, er = util.EncryptAndEncode(e.Value, os.Getenv(envSecret)); er != nil {
		return er
	}

	e.Value = encryptedValue
	return
}

// AfterSave decrypt secret value after encrypted value has been retrieved from database
func (e *Secret) AfterSave(tx *gorm.DB) (err error) {
	var er error
	var decrypted string
	if decrypted, er = util.DecryptAndDecode(e.Value, os.Getenv(envSecret)); er != nil {
		return er
	}

	e.Value = decrypted
	return
}

// AfterFind decrypt secret value after encrypted value has been retrieved from database
func (e *Secret) AfterFind(tx *gorm.DB) (err error) {
	var er error
	var decrypted string
	if decrypted, er = util.DecryptAndDecode(e.Value, os.Getenv(envSecret)); er != nil {
		return er
	}

	e.Value = decrypted
	return
}

// Secret entity holding information for secrets
type Secret struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;unique;not null"`
	Key       string    `gorm:"unique;not null"`
	Value     string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"time;autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"time;autoUpdateTime;not null"`
}

func (e *Action) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	return
}

// Action entity holding information for actions
type Action struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;unique;not null"`
	Label            string    `gorm:"not null"`
	Type             string    `gorm:"not null"`
	MatchEvent       *string   `gorm:""`
	MatchApplication *string   `gorm:""`
	MatchProvider    *string   `gorm:""`
	MatchHost        *string   `gorm:""`
	Payload          JSONMap   `gorm:"jsonb"`
	Enabled          bool      `gorm:"not null"`
	CreatedAt        time.Time `gorm:"time;autoCreateTime;not null"`
	UpdatedAt        time.Time `gorm:"time;autoUpdateTime;not null"`
}

func (e *ActionInvocation) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	return
}

// ActionInvocation entity holding information for invocations of actions
type ActionInvocation struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;unique;not null"`
	RetryCount int       `gorm:"not null;default:1"`
	State      string    `gorm:"not null"`
	Message    *string
	Event      Event     `gorm:"foreignKey:EventID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	EventID    string    `gorm:"not null"`
	Action     Action    `gorm:"foreignKey:ActionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ActionID   string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"time;autoCreateTime;not null"`
	UpdatedAt  time.Time `gorm:"time;autoUpdateTime;not null"`
}
