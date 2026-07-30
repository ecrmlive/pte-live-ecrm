package authjwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	PortalPlatform = "platform"
	PortalMerchant = "merchant"
	PortalApp      = "app"
	PortalOpen     = "open"
	PortalManager  = "manager"

	TokenAccess  = "access"
	TokenRefresh = "refresh"

	ScopeCUser        = "c_user"
	ScopeAdminConsole = "admin_console"
	ScopeStoreConsole = "store_console"
	ScopeMerchantAPI  = "merchant_api"
	ScopeOpenClient   = "open_client"

	PrincipalCUser        = "c_user"
	PrincipalAdminUser    = "admin_user"
	PrincipalStoreAccount = "store_account"
	PrincipalMerchantUser = "merchant_user"
	PrincipalOpenClient   = "open_client"

	ContextPlatform = "platform"
	ContextAdmin    = "admin"
	ContextMerchant = "merchant"
	ContextStore    = "store"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrWrongPortal  = errors.New("token portal mismatch")
	ErrWrongKind    = errors.New("token kind mismatch")
)

type Claims struct {
	Portal           string   `json:"portal"`
	Kind             string   `json:"kind"`
	Scope            string   `json:"scope"`
	PrincipalType    string   `json:"principal_type"`
	PrincipalID      uint     `json:"principal_id"`
	Roles            []string `json:"roles"`
	ClientPlatform   string   `json:"client_platform"`
	AuthContext      string   `json:"auth_context"`
	IdentityVersion  uint64   `json:"identity_version"`
	DataScopeVersion uint64   `json:"data_scope_version"`
	MerchantID       uint     `json:"merchant_id"`
	StoreID          uint     `json:"store_id"`
	MerchantAppID    string   `json:"merchant_app_id"`
	IMSDKAppID       string   `json:"im_sdk_app_id"`
	SessionID        string   `json:"session_id"`

	// 下列字段是既有服务的兼容字段；新代码以 principal_*、client_platform、
	// merchant_app_id 为准。StoreAppID 必须与 MerchantAppID 一致。
	Channel    string `json:"channel,omitempty"`
	AdminID    uint   `json:"admin_id,omitempty"`
	MerID      uint   `json:"mer_id,omitempty"`
	StoreAppID string `json:"store_app_id,omitempty"`
	UID        uint   `json:"uid,omitempty"`
	Account    string `json:"account,omitempty"`
	jwt.RegisteredClaims
}

// Identity 是签发 JWT 时唯一允许传入的身份快照。身份、租户和当前 IM SDK
// 必须由服务端数据库解析，客户端请求体不得直接提供这些字段。
type Identity struct {
	Portal           string
	Scope            string
	PrincipalType    string
	PrincipalID      uint
	Roles            []string
	ClientPlatform   string
	AuthContext      string
	IdentityVersion  uint64
	DataScopeVersion uint64
	MerchantID       uint
	StoreID          uint
	MerchantAppID    string
	IMSDKAppID       string
	Account          string
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
	identity := Identity{
		Portal: portal, PrincipalID: firstNonZero(uid, adminID, merID), Account: account,
		MerchantID: merID, IdentityVersion: 1,
	}
	switch portal {
	case PortalApp:
		identity.Scope, identity.PrincipalType, identity.Roles, identity.ClientPlatform, identity.AuthContext = ScopeCUser, PrincipalCUser, []string{"customer"}, "legacy_app", ContextPlatform
	case PortalOpen:
		identity.Scope, identity.PrincipalType, identity.Roles, identity.ClientPlatform, identity.AuthContext = ScopeOpenClient, PrincipalOpenClient, []string{"open_client"}, "server_to_server", ContextMerchant
	case PortalMerchant, PortalManager:
		identity.Scope, identity.PrincipalType, identity.Roles, identity.ClientPlatform, identity.AuthContext = ScopeMerchantAPI, PrincipalMerchantUser, []string{"merchant"}, "merchant_web", ContextMerchant
	default:
		identity.Scope, identity.PrincipalType, identity.Roles, identity.ClientPlatform, identity.AuthContext = ScopeAdminConsole, PrincipalAdminUser, []string{"platform"}, "admin_web", ContextAdmin
	}
	return m.issueIdentity(identity)
}

// IssueCUser 签发六端共享的 C 端 JWT；JWT subject 始终是统一 user ID，channel 只用于审计与风控。
func (m *Manager) IssueCUser(uid uint, account, channel string) (*Pair, error) {
	return m.IssueCUserWithIdentityVersion(uid, account, channel, 1)
}

// IssueCUserWithIdentityVersion 将账号身份版本写入令牌；密码修改、账号禁用、
// 强制退出或登录方式解绑时必须递增该版本。
func (m *Manager) IssueCUserWithIdentityVersion(uid uint, account, channel string, identityVersion uint64) (*Pair, error) {
	if uid == 0 || channel == "" {
		return nil, ErrInvalidToken
	}
	return m.issueIdentity(Identity{
		Portal: PortalApp, Scope: ScopeCUser, PrincipalType: PrincipalCUser,
		PrincipalID: uid, Roles: []string{"customer"}, ClientPlatform: channel,
		AuthContext: ContextPlatform, IdentityVersion: identityVersion, Account: account,
	})
}

// IssueCUserStoreContext 仅在 C 端已通过 X-AppId 解析到店铺后签发。
// 同一个用户可以在不同店铺之间切换，因此商户/AppId/IM 不能在全局登录令牌中伪造。
func (m *Manager) IssueCUserStoreContext(uid uint, account, channel string, identityVersion uint64, merchantID, storeID uint, merchantAppID, imSDKAppID string) (*Pair, error) {
	if merchantID == 0 || storeID == 0 || merchantAppID == "" {
		return nil, ErrInvalidToken
	}
	return m.issueIdentity(Identity{
		Portal: PortalApp, Scope: ScopeCUser, PrincipalType: PrincipalCUser,
		PrincipalID: uid, Roles: []string{"customer"}, ClientPlatform: channel,
		AuthContext: ContextStore, IdentityVersion: identityVersion, MerchantID: merchantID,
		StoreID: storeID, MerchantAppID: merchantAppID, IMSDKAppID: imSDKAppID, Account: account,
	})
}

// IssueAdminConsole 签发平台、商户、区域、客服、运营共用的后台 JWT。
func (m *Manager) IssueAdminConsole(adminID uint, account string, roles []string, dataScopeVersion uint64) (*Pair, error) {
	return m.IssueAdminConsoleWithIdentityVersion(adminID, account, roles, dataScopeVersion, 1)
}

func (m *Manager) IssueAdminConsoleWithIdentityVersion(adminID uint, account string, roles []string, dataScopeVersion, identityVersion uint64) (*Pair, error) {
	if adminID == 0 || len(roles) == 0 {
		return nil, ErrInvalidToken
	}
	return m.issueIdentity(Identity{
		Portal: PortalPlatform, Scope: ScopeAdminConsole, PrincipalType: PrincipalAdminUser,
		PrincipalID: adminID, Roles: roles, ClientPlatform: "admin_web", AuthContext: ContextAdmin,
		IdentityVersion: identityVersion, DataScopeVersion: dataScopeVersion, Account: account,
	})
}

// IssueStoreConsole 签发独立店铺后台令牌。merchant_id 与 store_id 同时写入，
// 后续每个店铺接口都必须以 store_id 做数据隔离。
func (m *Manager) IssueStoreConsole(accountID, merchantID, storeID uint, storeAppID, imSDKAppID, account, role string) (*Pair, error) {
	return m.IssueStoreConsoleWithIdentityVersion(accountID, merchantID, storeID, storeAppID, imSDKAppID, account, role, 1)
}

func (m *Manager) IssueStoreConsoleWithIdentityVersion(accountID, merchantID, storeID uint, storeAppID, imSDKAppID, account, role string, identityVersion uint64) (*Pair, error) {
	if accountID == 0 || merchantID == 0 || storeID == 0 || storeAppID == "" || role == "" {
		return nil, ErrInvalidToken
	}
	return m.issueIdentity(Identity{
		Portal: PortalMerchant, Scope: ScopeStoreConsole, PrincipalType: PrincipalStoreAccount,
		PrincipalID: accountID, Roles: []string{role}, ClientPlatform: "merchant_web", AuthContext: ContextStore,
		IdentityVersion: identityVersion, MerchantID: merchantID, StoreID: storeID, MerchantAppID: storeAppID,
		IMSDKAppID: imSDKAppID, Account: account,
	})
}

func (m *Manager) issueIdentity(identity Identity) (*Pair, error) {
	if !validIdentity(identity) {
		return nil, ErrInvalidToken
	}
	sessionID, err := newOpaqueID()
	if err != nil {
		return nil, err
	}
	access, err := m.sign(identity, TokenAccess, sessionID, m.accessTTL)
	if err != nil {
		return nil, err
	}
	refresh, err := m.sign(identity, TokenRefresh, sessionID, m.refreshTTL)
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
	if !containsAudience(claims.Audience, audienceForPortal(portal)) || !validClaimsIdentity(claims) {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (m *Manager) sign(identity Identity, kind, sessionID string, ttl time.Duration) (string, error) {
	now := time.Now()
	jti, err := newOpaqueID()
	if err != nil {
		return "", err
	}
	principalID := identity.PrincipalID
	claims := Claims{
		Portal:           identity.Portal,
		Kind:             kind,
		Scope:            identity.Scope,
		PrincipalType:    identity.PrincipalType,
		PrincipalID:      principalID,
		Roles:            append([]string(nil), identity.Roles...),
		ClientPlatform:   identity.ClientPlatform,
		AuthContext:      identity.AuthContext,
		IdentityVersion:  identity.IdentityVersion,
		DataScopeVersion: identity.DataScopeVersion,
		MerchantID:       identity.MerchantID,
		StoreID:          identity.StoreID,
		MerchantAppID:    identity.MerchantAppID,
		IMSDKAppID:       identity.IMSDKAppID,
		SessionID:        sessionID,
		Channel:          identity.ClientPlatform,
		MerID:            identity.MerchantID,
		StoreAppID:       identity.MerchantAppID,
		Account:          identity.Account,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(principalID), 10),
			Audience:  jwt.ClaimStrings{audienceForPortal(identity.Portal)},
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Issuer:    "qixi-mergers",
			ID:        jti,
		},
	}
	switch identity.PrincipalType {
	case PrincipalCUser:
		claims.UID = principalID
	case PrincipalAdminUser, PrincipalStoreAccount, PrincipalMerchantUser:
		claims.AdminID = principalID
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(m.secret)
}

func audienceForPortal(portal string) string { return "qixi-mergers:" + portal }

func containsAudience(audiences jwt.ClaimStrings, expected string) bool {
	for _, audience := range audiences {
		if audience == expected {
			return true
		}
	}
	return false
}

func validIdentity(identity Identity) bool {
	if identity.Portal == "" || identity.Scope == "" || identity.PrincipalType == "" || identity.PrincipalID == 0 ||
		len(identity.Roles) == 0 || identity.ClientPlatform == "" || identity.AuthContext == "" || identity.IdentityVersion == 0 {
		return false
	}
	if identity.AuthContext == ContextStore {
		return identity.MerchantID != 0 && identity.StoreID != 0 && identity.MerchantAppID != ""
	}
	return true
}

func validClaimsIdentity(claims *Claims) bool {
	if claims == nil || claims.PrincipalID == 0 || claims.PrincipalType == "" || len(claims.Roles) == 0 ||
		claims.ClientPlatform == "" || claims.AuthContext == "" || claims.IdentityVersion == 0 || claims.SessionID == "" ||
		claims.ID == "" || claims.Subject != strconv.FormatUint(uint64(claims.PrincipalID), 10) {
		return false
	}
	if claims.AuthContext == ContextStore {
		return claims.MerchantID != 0 && claims.MerID == claims.MerchantID && claims.StoreID != 0 && claims.MerchantAppID != "" && claims.StoreAppID == claims.MerchantAppID
	}
	return true
}

func newOpaqueID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func firstNonZero(values ...uint) uint {
	for _, value := range values {
		if value != 0 {
			return value
		}
	}
	return 0
}
