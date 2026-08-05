package chat

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/infra/imclient"
	"gorm.io/gorm"
)

type Store interface {
	FindThreadByMerUID(ctx context.Context, merID, uid uint) (*Thread, error)
	GetThread(ctx context.Context, id uint) (*Thread, error)
	CreateThread(ctx context.Context, t *Thread) error
	UpdateThread(ctx context.Context, t *Thread) error
	ListThreadsByMer(ctx context.Context, merID uint, page, limit int) ([]Thread, int64, error)
	ListThreadsByUID(ctx context.Context, uid uint, page, limit int) ([]Thread, int64, error)
	CreateMessage(ctx context.Context, m *Message) error
	ListMessages(ctx context.Context, threadID uint, page, limit int) ([]Message, int64, error)
	EnsureIdentity(ctx context.Context, portal string, localID uint, imUserID string, imUserNum int64) (*Identity, error)
	FindCustomerServiceID(ctx context.Context, merID uint) (uint, error)
	LoadUserNickname(ctx context.Context, uid uint) (string, error)
	LoadMerName(ctx context.Context, merID uint) (string, error)
	LoadMerchantIMConfig(ctx context.Context, merID uint) (*MerchantIMConfig, error)
}

type MerchantIMConfig struct {
	SDKAppID     string
	APIPublicURL string
	WSPublicURL  string
}

type IMSettings struct {
	Mode         string // remote
	APIBase      string // server-to-server URL; may be a Docker alias
	APIPublicURL string // browser / mini-program reachable IM API URL
	WSPublicURL  string
	AppID        string
	Token        string
	Secret       string // retained for config compatibility; never used to issue a local UserSig
}

// ConfigResolver 由平台云服务配置中心实现。配置读取仅在服务端发生，绝不返回服务令牌。
type ConfigResolver interface {
	Values(ctx context.Context, group string) (map[string]string, error)
}

type Service struct {
	store    Store
	im       IMSettings
	resolver ConfigResolver
}

func NewService(store Store, im IMSettings, resolvers ...ConfigResolver) *Service {
	// Production path is remote-only; local fake UserSig is not used.
	if im.Mode == "" {
		im.Mode = "remote"
	}
	if im.AppID == "" {
		im.AppID = "30001"
	}
	var resolver ConfigResolver
	if len(resolvers) > 0 {
		resolver = resolvers[0]
	}
	return &Service{store: store, im: im, resolver: resolver}
}

func (s *Service) OpenUserThread(ctx context.Context, uid uint, merID uint) (*Thread, error) {
	if uid == 0 || merID == 0 {
		return nil, ErrBadParam
	}
	row, err := s.store.FindThreadByMerUID(ctx, merID, uid)
	if err == nil {
		if row.Status == StatusClosed {
			now := time.Now()
			row.Status, row.LastMsg, row.LastTime = StatusOpen, "会话已重新开启", &now
			if updateErr := s.store.UpdateThread(ctx, row); updateErr != nil {
				return nil, updateErr
			}
		}
		_ = s.ensureIMConversation(ctx, row)
		return s.enrichThread(ctx, row), nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	now := time.Now()
	svcID, _ := s.store.FindCustomerServiceID(ctx, merID)
	t := &Thread{
		MerID: merID, UID: uid, ServiceID: svcID, Status: StatusOpen,
		ImConversationID: fmt.Sprintf("pending:%d:%d:%d", merID, uid, now.UnixNano()),
		LastMsg:          "会话已创建", LastTime: &now, CreateTime: now,
	}
	if svcID > 0 {
		t.AssignedAt = &now
	}
	if err := s.store.CreateThread(ctx, t); err != nil {
		// uk_user_store closes the concurrent-open race. Re-read the winner instead of creating duplicate business sessions.
		winner, findErr := s.store.FindThreadByMerUID(ctx, merID, uid)
		if findErr != nil {
			return nil, err
		}
		_ = s.ensureIMConversation(ctx, winner)
		return s.enrichThread(ctx, winner), nil
	}
	sys := &Message{
		ThreadID: t.ThreadID, MerID: merID, SenderRole: "system", SenderID: 0,
		MsgType: MsgSystem, Content: "您好，客服稍后为您服务", CreateTime: now,
	}
	_ = s.store.CreateMessage(ctx, sys)
	_ = s.ensureIMConversation(ctx, t)
	return s.enrichThread(ctx, t), nil
}

func (s *Service) ListUserThreads(ctx context.Context, uid uint, page, limit int) (*PageResult[Thread], error) {
	if uid == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListThreadsByUID(ctx, uid, page, limit)
	if err != nil {
		return nil, err
	}
	for i := range list {
		list[i] = *s.enrichThread(ctx, &list[i])
	}
	return &PageResult[Thread]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListServiceThreads(ctx context.Context, merID uint, page, limit int) (*PageResult[Thread], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListThreadsByMer(ctx, merID, page, limit)
	if err != nil {
		return nil, err
	}
	for i := range list {
		list[i] = *s.enrichThread(ctx, &list[i])
	}
	return &PageResult[Thread]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) Claim(ctx context.Context, merID, serviceID, threadID uint) (*Thread, error) {
	t, err := s.getOwned(ctx, merID, 0, threadID, false)
	if err != nil {
		return nil, err
	}
	if t.ServiceID == 0 || t.ServiceID == serviceID {
		t.ServiceID = serviceID
		now := time.Now()
		t.AssignedAt = &now
		if err := s.store.UpdateThread(ctx, t); err != nil {
			return nil, err
		}
		_ = s.ensureIMConversation(ctx, t)
	} else if t.ServiceID != serviceID {
		return nil, ErrForbidden
	}
	return s.enrichThread(ctx, t), nil
}

func (s *Service) ListMessages(ctx context.Context, merID, uid, serviceID, threadID uint, page, limit int) (*PageResult[Message], error) {
	t, err := s.getOwned(ctx, merID, uid, threadID, serviceID > 0)
	if err != nil {
		return nil, err
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListMessages(ctx, threadID, page, limit)
	if err != nil {
		return nil, err
	}
	if uid > 0 && t.UserUnread > 0 {
		t.UserUnread = 0
		_ = s.store.UpdateThread(ctx, t)
	}
	if serviceID > 0 && t.ServiceUnread > 0 {
		t.ServiceUnread = 0
		_ = s.store.UpdateThread(ctx, t)
	}
	return &PageResult[Message]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) SendUser(ctx context.Context, uid, threadID uint, in SendInput) (*Message, error) {
	t, err := s.getOwned(ctx, 0, uid, threadID, false)
	if err != nil {
		return nil, err
	}
	return s.send(ctx, t, RoleUser, uid, in)
}

func (s *Service) SendService(ctx context.Context, merID, serviceID, threadID uint, in SendInput) (*Message, error) {
	t, err := s.getOwned(ctx, merID, 0, threadID, true)
	if err != nil {
		return nil, err
	}
	if t.ServiceID == 0 {
		t.ServiceID = serviceID
		now := time.Now()
		t.AssignedAt = &now
	} else if t.ServiceID != serviceID {
		return nil, ErrForbidden
	}
	return s.send(ctx, t, RoleService, serviceID, in)
}

func (s *Service) send(ctx context.Context, t *Thread, role string, senderID uint, in SendInput) (*Message, error) {
	if t.Status == StatusClosed {
		return nil, ErrClosed
	}
	content := strings.TrimSpace(in.Content)
	if content == "" {
		return nil, ErrBadParam
	}
	msgType := strings.TrimSpace(in.MsgType)
	if msgType == "" {
		msgType = MsgText
	}
	// Text body must go through PTE E2EE SDK; /cs only keeps business metadata (order cards, etc.).
	if msgType == MsgText || msgType == "" {
		return nil, ErrTextViaCS
	}
	if msgType != MsgOrder && msgType != MsgSystem {
		return nil, ErrBadParam
	}
	now := time.Now()
	m := &Message{
		ThreadID: t.ThreadID, MerID: t.MerID, SenderRole: role, SenderID: senderID,
		MsgType: msgType, Content: content, CreateTime: now,
	}
	if err := s.store.CreateMessage(ctx, m); err != nil {
		return nil, err
	}
	preview := content
	if utf8.RuneCountInString(preview) > 80 {
		preview = string([]rune(preview)[:80])
	}
	t.LastMsg = preview
	t.LastTime = &now
	if role == RoleUser {
		t.ServiceUnread++
	} else {
		t.UserUnread++
	}
	_ = s.store.UpdateThread(ctx, t)
	_ = s.ensureIMConversation(ctx, t)
	return m, nil
}

func (s *Service) IssueCredential(ctx context.Context, portal string, localID uint) (*Credential, error) {
	return s.IssueCredentialForThread(ctx, portal, localID, 0)
}

func (s *Service) IssueCredentialForThread(ctx context.Context, portal string, localID uint, threadID uint) (*Credential, error) {
	if localID == 0 || (portal != "app" && portal != "service") {
		return nil, ErrBadParam
	}
	num := NumericIMUserID(portal, localID)
	ident := fmt.Sprintf("%d", num)
	if _, err := s.store.EnsureIdentity(ctx, portal, localID, ident, num); err != nil {
		return nil, err
	}

	var convID uint64
	var merchantID uint
	if threadID > 0 {
		var t *Thread
		var err error
		if portal == "app" {
			t, err = s.getOwned(ctx, 0, localID, threadID, false)
		} else {
			t, err = s.getOwned(ctx, 0, 0, threadID, true)
			if err == nil && t.ServiceID != localID {
				return nil, ErrForbidden
			}
		}
		if err != nil {
			return nil, err
		}
		_ = s.ensureIMConversation(ctx, t)
		convID, _ = strconv.ParseUint(t.ImConversationID, 10, 64)
		merchantID = t.MerID
	}

	im, err := s.effectiveIMSettingsForMerchant(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	remote := newRemoteClient(im)
	if remote == nil || !remote.Enabled() {
		return nil, ErrIMRemoteRequired
	}
	userType := "user"
	if portal == "service" {
		userType = "staff"
	}
	res, err := remote.IssueUserSig(ctx, num, ident, userType, "qixi-web", "web")
	if err != nil {
		return nil, fmt.Errorf("im usersig: %w", err)
	}
	ws := strings.TrimSpace(im.WSPublicURL)
	if ws == "" {
		ws = res.WsURL
	}
	apiURL := strings.TrimSpace(im.APIPublicURL)
	if apiURL == "" {
		return nil, fmt.Errorf("im api_public_url is required for client SDK")
	}
	if ws == "" {
		return nil, fmt.Errorf("im ws_public_url is required for client SDK")
	}
	return &Credential{
		Mode: "remote", AppID: firstNonEmpty(res.AppID, im.AppID), SDKAppID: res.SDKAppID,
		ImUserID: ident, Identifier: firstNonEmpty(res.Identifier, ident),
		UserSig: res.UserSig, ExpireAt: res.ExpireAt, APIURL: apiURL, WSURL: ws, WSHint: ws,
		ImConversationID: convID,
		Note:             "pte-live-im UserSig；文本经 Web SDK E2EE 发送，本仓 /cs 仅订单卡片等元数据",
	}, nil
}

func (s *Service) ensureIMConversation(ctx context.Context, t *Thread) error {
	if t == nil {
		return nil
	}
	im, err := s.effectiveIMSettingsForMerchant(ctx, t.MerID)
	if err != nil {
		return err
	}
	remote := newRemoteClient(im)
	if remote == nil || !remote.Enabled() {
		return nil
	}
	if convID, parseErr := strconv.ParseUint(t.ImConversationID, 10, 64); parseErr == nil && convID > 0 {
		return nil
	}
	if t.ServiceID == 0 {
		svcID, err := s.store.FindCustomerServiceID(ctx, t.MerID)
		if err != nil || svcID == 0 {
			return nil
		}
		t.ServiceID = svcID
	}
	userNum := NumericIMUserID("app", t.UID)
	svcNum := NumericIMUserID("service", t.ServiceID)
	_, _ = s.store.EnsureIdentity(ctx, "app", t.UID, fmt.Sprintf("%d", userNum), userNum)
	_, _ = s.store.EnsureIdentity(ctx, "service", t.ServiceID, fmt.Sprintf("%d", svcNum), svcNum)
	convID, err := remote.OpenSingle(ctx, userNum, svcNum)
	if err != nil {
		return err
	}
	t.ImConversationID = strconv.FormatUint(convID, 10)
	return s.store.UpdateThread(ctx, t)
}

func (s *Service) effectiveIMSettingsForMerchant(ctx context.Context, merID uint) (IMSettings, error) {
	im, err := s.effectiveIMSettings(ctx)
	if err != nil || merID == 0 {
		return im, err
	}
	merchant, err := s.store.LoadMerchantIMConfig(ctx, merID)
	if err != nil {
		return IMSettings{}, err
	}
	if merchant != nil {
		if value := strings.TrimSpace(merchant.SDKAppID); value != "" {
			im.AppID = value
		}
		if value := strings.TrimSpace(merchant.APIPublicURL); value != "" {
			im.APIPublicURL = value
		}
		if value := strings.TrimSpace(merchant.WSPublicURL); value != "" {
			im.WSPublicURL = value
		}
	}
	return im, nil
}

func (s *Service) effectiveIMSettings(ctx context.Context) (IMSettings, error) {
	im := s.im
	if s.resolver == nil {
		return im, nil
	}
	values, err := s.resolver.Values(ctx, "im")
	if err != nil {
		return IMSettings{}, err
	}
	if value := strings.TrimSpace(values["mode"]); value != "" {
		im.Mode = value
	}
	if value := strings.TrimSpace(values["api_base"]); value != "" {
		im.APIBase = value
	}
	if value := strings.TrimSpace(values["api_public_url"]); value != "" {
		im.APIPublicURL = value
	}
	if value := strings.TrimSpace(values["ws_public_url"]); value != "" {
		im.WSPublicURL = value
	}
	if value := strings.TrimSpace(values["app_id"]); value != "" {
		im.AppID = value
	}
	if value := strings.TrimSpace(values["integration_token"]); value != "" {
		im.Token = value
	}
	return im, nil
}

func newRemoteClient(im IMSettings) *imclient.Client {
	appID, _ := strconv.Atoi(im.AppID)
	if !strings.EqualFold(im.Mode, "remote") || strings.TrimSpace(im.APIBase) == "" || strings.TrimSpace(im.Token) == "" {
		return nil
	}
	return imclient.New(im.APIBase, im.Token, appID)
}

func (s *Service) getOwned(ctx context.Context, merID, uid, threadID uint, asService bool) (*Thread, error) {
	t, err := s.store.GetThread(ctx, threadID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if uid > 0 && t.UID != uid {
		return nil, ErrForbidden
	}
	if asService || merID > 0 {
		if merID > 0 && t.MerID != merID {
			return nil, ErrForbidden
		}
	}
	return t, nil
}

func (s *Service) enrichThread(ctx context.Context, t *Thread) *Thread {
	if t == nil {
		return nil
	}
	if name, err := s.store.LoadUserNickname(ctx, t.UID); err == nil {
		t.UserNickname = name
	}
	if name, err := s.store.LoadMerName(ctx, t.MerID); err == nil {
		t.MerName = name
	}
	return t
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func normalize(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return page, limit
}
