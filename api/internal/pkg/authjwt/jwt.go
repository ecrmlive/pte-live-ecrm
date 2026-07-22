package authjwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	PortalPlatform = "platform"
	PortalMerchant = "merchant"
	PortalApp      = "app"
	PortalOpen     = "open"
	PortalService  = "service"
	PortalManager  = "manager"

	TokenAccess  = "access"
	TokenRefresh = "refresh"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrWrongPortal  = errors.New("token portal mismatch")
	ErrWrongKind    = errors.New("token kind mismatch")
)

type Claims struct {
	Portal  string `json:"portal"`
	Kind    string `json:"kind"`
	AdminID uint   `json:"admin_id,omitempty"`
	MerID   uint   `json:"mer_id,omitempty"`
	UID     uint   `json:"uid,omitempty"`
	Account string `json:"account,omitempty"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewManager(secret string, accessTTL, refreshTTL time.Duration) *Manager {
	return &Manager{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

type Pair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (m *Manager) Issue(portal string, adminID, merID, uid uint, account string) (*Pair, error) {
	access, err := m.sign(portal, TokenAccess, adminID, merID, uid, account, m.accessTTL)
	if err != nil {
		return nil, err
	}
	refresh, err := m.sign(portal, TokenRefresh, adminID, merID, uid, account, m.refreshTTL)
	if err != nil {
		return nil, err
	}
	return &Pair{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int64(m.accessTTL.Seconds()),
	}, nil
}

func (m *Manager) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (m *Manager) ParseExpect(tokenString, portal, kind string) (*Claims, error) {
	claims, err := m.Parse(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.Portal != portal {
		return nil, ErrWrongPortal
	}
	if claims.Kind != kind {
		return nil, ErrWrongKind
	}
	return claims, nil
}

func (m *Manager) sign(portal, kind string, adminID, merID, uid uint, account string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		Portal:  portal,
		Kind:    kind,
		AdminID: adminID,
		MerID:   merID,
		UID:     uid,
		Account: account,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Issuer:    "qixi-mergers",
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(m.secret)
}
