package helper

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jamaalul/proyek1pbl/app/model"
)

var (
	ErrInvalidToken = errors.New("token tidak valid")
	ErrExpiredToken = errors.New("token sudah kedaluwarsa")
)

type JWTClaims struct {
	UserID int64  `json:"user_id,omitempty"`
	Role   string `json:"role"`
	Email  string `json:"email,omitempty"`
	JTI    string `json:"jti,omitempty"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret    []byte
	issuer    string
	accessTTL time.Duration
}

func NewJWTManager(secret, issuer string, accessTTL time.Duration) *JWTManager {
	return &JWTManager{
		secret:    []byte(secret),
		issuer:    issuer,
		accessTTL: accessTTL,
	}
}

func (m *JWTManager) AccessTTL() time.Duration {
	return m.accessTTL
}

func (m *JWTManager) ExpiresIn() int {
	return int(m.accessTTL.Seconds())
}

func (m *JWTManager) GenerateAccess(u model.User) (string, error) {
	tokenStr, _, _, err := m.GenerateAccessWithJTI(u)
	return tokenStr, err
}

func (m *JWTManager) GenerateAccessWithJTI(u model.User) (string, time.Time, string, error) {
	now := time.Now()
	expiresAt := now.Add(m.accessTTL)
	jti := uuid.New().String()

	claims := JWTClaims{
		UserID: u.ID,
		Role:   string(u.Role),
		Email:  u.Email,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(u.ID, 10),
			ID:        jti,
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(m.secret)
	return tokenStr, expiresAt, jti, err
}

func (m *JWTManager) GenerateTokenFor(userID int64, email, role string) (string, error) {
	u := model.User{
		ID:    userID,
		Email: email,
		Role:  model.Role(role),
	}
	return m.GenerateAccess(u)
}

func (m *JWTManager) Parse(tokenString string) (model.AuthUser, error) {
	claims, err := m.ParseTokenClaims(tokenString)
	if err != nil {
		return model.AuthUser{}, err
	}

	userID := claims.UserID
	if userID == 0 && claims.Subject != "" {
		if id, parseErr := strconv.ParseInt(claims.Subject, 10, 64); parseErr == nil {
			userID = id
		}
	}

	jti := claims.JTI
	if jti == "" {
		jti = claims.ID
	}

	return model.AuthUser{
		UserID: userID,
		Email:  claims.Email,
		Role:   claims.Role,
		JTI:    jti,
	}, nil
}

func (m *JWTManager) ParseTokenClaims(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
