// filename: user.go

package domain

import (
	"time" // Package for time-related functions

	"gorm.io/gorm" // GORM library for ORM functionalities
)

func init() {
	registeredDomains = append(registeredDomains, &User{})
}

type User struct {
	gorm.Model         // ID, CreatedAt, UpdatedAt, DeletedAt
	ID         string  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"` // Should be string
	Name       *string `gorm:"size:255"`                                         // Should be *string (pointer)
	Email      string  `gorm:"type:varchar(100);uniqueIndex;not null"`
	// EmailVerified stores the timestamp when the user's email was verified. Null if not verified.
	EmailVerified *time.Time `gorm:"null"`
	FullName      string     `gorm:"type:varchar(255);null"` // Can be 'name' from provider
	ImageURL      *string    `gorm:"size:255"`
	LastLoginAt   *time.Time // New field	PasswordHash  string     `gorm:"null"`                   // Nullable for OAuth-only users

	IsActive    bool    `gorm:"default:true"`
	Bio         string  `gorm:"type:text;null"`
	DisplayName *string `gorm:"type:varchar(255);default:null"` // User's display name, nullable
	// Relationships
	LinkedAccounts []LinkedAccount `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// VerificationTokens []VerificationToken `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // If linking tokens to user ID
}
