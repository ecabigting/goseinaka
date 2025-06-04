// filename: linked_account.go
package domain

import (
	"time"

	"gorm.io/gorm"
)

func init() {
	registeredDomains = append(registeredDomains, &LinkedAccount{})
}

type LinkedAccount struct {
	gorm.Model             // ID, CreatedAt, UpdatedAt, DeletedAt
	UserID         string  `gorm:"type:uuid;not null;index:idx_user_provider,unique"` // Should be string
	Provider       string  `gorm:"size:50;not null;index:idx_user_provider,unique;index:idx_provider_userid"`
	ProviderUserID string  `gorm:"size:255;not null;index;index:idx_provider_userid,unique"`
	Email          *string `gorm:"size:255"`
	Name           *string `gorm:"size:255"`
	ImageURL       *string `gorm:"size:255"`
	AccessToken    *string `gorm:"type:text"`
	RefreshToken   *string `gorm:"type:text"`
	ExpiresAt      *time.Time
	Scope          *string `gorm:"size:500"`
	User           User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Consider CASCADE on delete
}
