// filename: verification_token.go

package domain

import (
	"time"

	"gorm.io/gorm"
)

func init() {
	registeredDomains = append(registeredDomains, &VerificationToken{})
}

// VerificationToken is used for one-time tokens, e.g., email verification, password reset.
type VerificationToken struct {
	gorm.Model // ID, CreatedAt, UpdatedAt, DeletedAt (ID is numeric by default with gorm.Model)

	// Identifier can be an email address (for email verification) or a user ID.
	// This makes it flexible. For email verification, it would be the email.
	// For password reset tied to a user, it could be user ID or email.
	Identifier string `gorm:"type:varchar(255);not null;uniqueIndex:idx_identifier_token"`

	// Token is the secure, unguessable token string.
	Token string `gorm:"type:varchar(255);not null;uniqueIndex;uniqueIndex:idx_identifier_token"`

	// Expires is the timestamp when this token is no longer valid.
	ExpiresAt time.Time `gorm:"not null"`

	// Optional: If you want to directly link a verification token to a user record via UserID.
	// UserID    uint `gorm:"null"` // Nullable if token is for an email not yet in users table
	// User      User  // Optional belongs to relationship
}
