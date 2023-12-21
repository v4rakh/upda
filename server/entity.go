package server

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func (u *Update) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

func (wh *Webhook) BeforeCreate(tx *gorm.DB) (err error) {
	wh.ID = uuid.New()
	return
}

func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
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

// Event entity holding information for events
type Event struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;unique;not null"`
	Name      string    `gorm:"not null"`
	State     string    `gorm:"not null"`
	Payload   JSONMap   `gorm:"jsonb"`
	CreatedAt time.Time `gorm:"time;autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"time;autoUpdateTime;not null"`
}
