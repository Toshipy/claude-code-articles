package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/toshipy/claude-code-articles/backend/internal/domain/model"
	"github.com/toshipy/claude-code-articles/backend/internal/repository/postgres"
)

type AuthService struct {
	userRepo  *postgres.UserRepository
	jwtSecret []byte
}

func NewAuthService(ur *postgres.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo: ur, jwtSecret: []byte(jwtSecret)}
}

type JWTClaims struct {
	UserID int64  `json:"sub"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type AuthResponse struct {
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	TokenType    string             `json:"token_type"`
	ExpiresIn    int                `json:"expires_in"`
	User         model.UserResponse `json:"user"`
}

func (s *AuthService) AuthenticateWithGoogle(email, displayName, avatarURL string) (*AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if user == nil {
		var av *string
		if avatarURL != "" {
			av = &avatarURL
		}
		user = &model.User{
			Email:       email,
			DisplayName: displayName,
			AvatarURL:   av,
			Role:        "user",
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User: model.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.DisplayName,
			AvatarURL: user.AvatarURL,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *AuthService) ValidateToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (s *AuthService) GetUserByID(id int64) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *AuthService) GetBookmarkCount(userID int64) (int, error) {
	return s.userRepo.GetBookmarkCount(userID)
}

func (s *AuthService) generateAccessToken(user *model.User) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
