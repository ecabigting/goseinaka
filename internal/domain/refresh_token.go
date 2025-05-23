// filename: refresh_token.model.go
package domain

import (
	"time"
)

func init() {
	registeredDomains = append(registeredDomains, &RefreshToken{})
}

// RefreshToken stores API-specific refresh tokens.
// These are used by clients to obtain new access tokens.
type RefreshToken struct {
	ID            string     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID        string     `gorm:"type:uuid;not null;index"`               // Foreign key to the User model
	TokenHash     string     `gorm:"type:varchar(255);uniqueIndex;not null"` // Hashed version of the refresh token for security
	ExpiresAt     time.Time  `gorm:"not null"`                               // Expiration timestamp of this refresh token
	CreatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	RevokedAt     *time.Time `gorm:"index"`    // Optional: For explicit revocation tracking if needed later
	CreatedFromIP *string    `gorm:"size:50"`  // Optional: For security auditing
	UserAgent     *string    `gorm:"size:255"` // Optional: For security auditing

	// Foreign key relationship
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
