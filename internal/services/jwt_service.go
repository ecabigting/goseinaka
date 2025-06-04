// filename: internal/services/jwt_service.go

package services

import (
	"fmt"
	"time"

	"github.com/ecabigting/goseinaka/internal/domain"
	"github.com/ecabigting/goseinaka/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JWTService struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	db              *gorm.DB
}

type JWTCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type Tokens struct {
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

func NewJWTService(sKey string, aTokenTTL int, rTokenTTL int, dB *gorm.DB) (*JWTService, error) {
	return &JWTService{
		secretKey:       sKey,
		accessTokenTTL:  time.Duration(aTokenTTL) * time.Minute,
		refreshTokenTTL: time.Duration(rTokenTTL) * time.Minute,
		db:              dB,
	}, nil
}

func (s *JWTService) GenerateTokens(userID string) (tokens *Tokens, err error) {
	accessTokenExpiresAt := time.Now().Add(s.accessTokenTTL)
	accessClaims := JWTCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "this-api",
			Subject:   userID,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshTokenExpiresAt := time.Now().Add(s.refreshTokenTTL)
	opaqueRefreshToken := uuid.New().String()
	hashedRefreshToken := utils.HashString(opaqueRefreshToken)
	dbRefreshToken := &domain.RefreshToken{
		UserID:    userID,
		TokenHash: hashedRefreshToken,
		ExpiresAt: refreshTokenExpiresAt,
	}

	if result := s.db.Create(dbRefreshToken); result.Error != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", result.Error)
	}

	return &Tokens{
		AccessToken:           signedAccessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          opaqueRefreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}
