package circle

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrNotFound          = errors.New("记录不存在")
	ErrBadParam          = errors.New("参数不合法")
	ErrHasChildren       = errors.New("存在下级区域，不能删除")
	ErrAlreadyAudited    = errors.New("代理申请已处理")
	ErrAgentBound        = errors.New("代理已关联区域，不能撤销")
	ErrAgentBalance      = errors.New("代理仍有佣金余额，不能撤销")
	ErrAgentAdminBound   = errors.New("代理仍关联统一后台账号，不能撤销")
	ErrAgentRevoked      = errors.New("代理资格已撤销")
	ErrAgentNotApproved  = errors.New("仅已审核通过的代理可以重置后台密码")
	ErrAgentAdminUnbound = errors.New("代理未关联启用的统一后台账号")
	ErrCommandConflict   = errors.New("幂等键与既有代理命令不一致")
)

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

type Store interface {
	ListCircles(context.Context, CircleListFilter, int, int) ([]Circle, int64, error)
	GetCircle(context.Context, uint) (*Circle, error)
	CreateCircle(context.Context, *Circle) error
	UpdateCircle(context.Context, *Circle) error
	DeleteCircle(context.Context, uint) error
	CountCircleChildren(context.Context, uint) (int64, error)

	ListAgents(context.Context, string, *int8, *int8, int, int) ([]Agent, int64, error)
	GetAgent(context.Context, uint) (*Agent, error)
	CreateAgent(context.Context, *Agent) error
	UpdateAgent(context.Context, *Agent) error
	AuditAgent(context.Context, uint, int8, string, uint, time.Time) error
	RevokeAgent(context.Context, uint, string, string, uint, time.Time) (bool, error)
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func normalizePage(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}

func (s *Service) ListCircles(ctx context.Context, filter CircleListFilter, page, limit int) (*PageResult[Circle], error) {
	page, limit = normalizePage(page, limit)
	filter.Keyword = strings.TrimSpace(filter.Keyword)
	rows, total, err := s.store.ListCircles(ctx, filter, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[Circle]{List: rows, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetCircle(ctx context.Context, id uint) (*Circle, error) {
	if id == 0 {
		return nil, ErrBadParam
	}
	row, err := s.store.GetCircle(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return row, nil
}

func (s *Service) CreateCircle(ctx context.Context, in CircleInput) (*Circle, error) {
	row, err := s.circleFromInput(ctx, 0, in)
	if err != nil {
		return nil, err
	}
	if err := s.store.CreateCircle(ctx, row); err != nil {
		return nil, err
	}
	row.Path = "/" + itoa(row.CircleID) + "/"
	if row.PID != 0 {
		parent, _ := s.GetCircle(ctx, row.PID)
		if parent != nil {
			row.Path = parent.Path + itoa(row.CircleID) + "/"
			row.Level = parent.Level + 1
		}
	}
	if err := s.store.UpdateCircle(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) UpdateCircle(ctx context.Context, id uint, in CircleInput) (*Circle, error) {
	if id == 0 || in.PID == id {
		return nil, ErrBadParam
	}
	old, err := s.GetCircle(ctx, id)
	if err != nil {
		return nil, err
	}
	row, err := s.circleFromInput(ctx, id, in)
	if err != nil {
		return nil, err
	}
	row.Path, row.Level = old.Path, old.Level
	if in.PID != old.PID {
		return nil, errors.New("区域创建后不支持直接调整父级，请先删除下级后重建")
	}
	if err := s.store.UpdateCircle(ctx, row); err != nil {
		return nil, err
	}
	return s.GetCircle(ctx, id)
}

func (s *Service) circleFromInput(ctx context.Context, id uint, in CircleInput) (*Circle, error) {
	name := strings.TrimSpace(in.Name)
	nameLimit := 64
	if in.Type == 1 {
		nameLimit = 20
	}
	if name == "" || len([]rune(name)) > nameLimit || (in.Status != 0 && in.Status != 1) || (in.Type != 0 && in.Type != 1) || (in.CommissionType != 0 && in.CommissionType != 1) || in.CommissionRate < 0 || in.CommissionRate > 100 {
		return nil, ErrBadParam
	}
	if in.Type == 1 && (in.CircleAgentID == 0 || in.RoleID == 0 || in.BusinessStoreCategory == 0 || in.BusinessStoreType == 0) {
		return nil, ErrBadParam
	}
	level := uint8(0)
	if in.PID != 0 {
		parent, err := s.GetCircle(ctx, in.PID)
		if err != nil {
			return nil, err
		}
		if parent.Type != in.Type || parent.Level >= 2 {
			return nil, ErrBadParam
		}
		level = parent.Level + 1
	}
	if in.CircleAgentID != 0 {
		agent, err := s.GetAgent(ctx, in.CircleAgentID)
		if err != nil {
			return nil, err
		}
		if agent.Status != AgentApproved || agent.Type != in.Type {
			return nil, ErrBadParam
		}
	}
	return &Circle{
		CircleID: id, PID: in.PID, Name: name, CircleAgentID: in.CircleAgentID,
		CommissionType: in.CommissionType, CommissionRate: in.CommissionRate,
		Remark: strings.TrimSpace(in.Remark), Sort: in.Sort, Status: in.Status, Type: in.Type,
		RoleID: in.RoleID, BusinessStoreCategory: in.BusinessStoreCategory, BusinessStoreType: in.BusinessStoreType,
		GoodsType: joinInts(in.GoodsTypes), PlatformCategoryIDs: joinUints(in.PlatformCategoryIDs), Level: level,
	}, nil
}

func (s *Service) DeleteCircle(ctx context.Context, id uint) error {
	if _, err := s.GetCircle(ctx, id); err != nil {
		return err
	}
	children, err := s.store.CountCircleChildren(ctx, id)
	if err != nil {
		return err
	}
	if children > 0 {
		return ErrHasChildren
	}
	return s.store.DeleteCircle(ctx, id)
}

func (s *Service) ListAgents(ctx context.Context, keyword string, status, agentType *int8, page, limit int) (*PageResult[Agent], error) {
	page, limit = normalizePage(page, limit)
	rows, total, err := s.store.ListAgents(ctx, strings.TrimSpace(keyword), status, agentType, page, limit)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		scrubAgentPayment(&rows[i])
	}
	return &PageResult[Agent]{List: rows, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetAgent(ctx context.Context, id uint) (*Agent, error) {
	row, err := s.getAgent(ctx, id)
	if err != nil {
		return nil, err
	}
	scrubAgentPayment(row)
	return row, nil
}

func (s *Service) getAgent(ctx context.Context, id uint) (*Agent, error) {
	if id == 0 {
		return nil, ErrBadParam
	}
	row, err := s.store.GetAgent(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return row, nil
}

func (s *Service) CreateAgent(ctx context.Context, in AgentInput) (*Agent, error) {
	row, err := agentFromInput(in)
	if err != nil {
		return nil, err
	}
	if err := s.store.CreateAgent(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) UpdateAgent(ctx context.Context, id uint, in AgentInput) (*Agent, error) {
	old, err := s.getAgent(ctx, id)
	if err != nil {
		return nil, err
	}
	if old.Status != AgentPending {
		return nil, ErrAlreadyAudited
	}
	row, err := agentFromInput(in)
	if err != nil {
		return nil, err
	}
	row.CircleAgentID = id
	// 支付账号、开户行、二维码只支持非空覆盖，避免脱敏读模型在编辑时把既有资料清空。
	if row.PaymentAccount == "" {
		row.PaymentAccount = old.PaymentAccount
	}
	if row.PaymentBank == "" {
		row.PaymentBank = old.PaymentBank
	}
	if row.PaymentQRImg == "" {
		row.PaymentQRImg = old.PaymentQRImg
	}
	if err := s.store.UpdateAgent(ctx, row); err != nil {
		return nil, err
	}
	return s.GetAgent(ctx, id)
}

func scrubAgentPayment(row *Agent) {
	if row == nil {
		return
	}
	row.PaymentConfigured = strings.TrimSpace(row.PaymentAccount) != "" || strings.TrimSpace(row.PaymentBank) != "" || strings.TrimSpace(row.PaymentQRImg) != ""
	row.PaymentAccount = ""
	row.PaymentBank = ""
	row.PaymentQRImg = ""
}

func (s *Service) AuditAgent(ctx context.Context, id uint, in AuditInput, adminID uint) error {
	if in.Status != AgentApproved && in.Status != AgentRejected {
		return ErrBadParam
	}
	if in.Status == AgentRejected && strings.TrimSpace(in.AuditReason) == "" {
		return ErrBadParam
	}
	row, err := s.GetAgent(ctx, id)
	if err != nil {
		return err
	}
	if row.Status != AgentPending {
		return ErrAlreadyAudited
	}
	return s.store.AuditAgent(ctx, id, in.Status, strings.TrimSpace(in.AuditReason), adminID, time.Now())
}

func (s *Service) RevokeAgent(ctx context.Context, id uint, reason, idempotencyKey string, adminID uint) (bool, error) {
	reason, idempotencyKey = strings.TrimSpace(reason), strings.TrimSpace(idempotencyKey)
	if id == 0 || len([]rune(reason)) < 2 || len([]rune(reason)) > 500 || len([]rune(idempotencyKey)) < 8 || len([]rune(idempotencyKey)) > 128 {
		return false, ErrBadParam
	}
	return s.store.RevokeAgent(ctx, id, reason, idempotencyKey, adminID, time.Now())
}

func agentFromInput(in AgentInput) (*Agent, error) {
	name, phone := strings.TrimSpace(in.Name), strings.TrimSpace(in.Phone)
	if name == "" || phone == "" || len([]rune(name)) > 64 || len(phone) > 16 || in.Type < 0 || in.Type > 1 || in.PaymentMethod > 2 {
		return nil, ErrBadParam
	}
	status := AgentPending
	// 平台「添加商户管理员」会带登录账号：直接按已通过资格创建，避免额外审核。
	if strings.TrimSpace(in.Account) != "" && strings.TrimSpace(in.Password) != "" {
		if len([]rune(strings.TrimSpace(in.Password))) < 6 {
			return nil, ErrBadParam
		}
		status = AgentApproved
	}
	return &Agent{UID: in.UID, Name: name, Phone: phone, Qualification: strings.TrimSpace(in.Qualification), Remark: strings.TrimSpace(in.Remark), PaymentMethod: in.PaymentMethod, PaymentName: strings.TrimSpace(in.PaymentName), PaymentAccount: strings.TrimSpace(in.PaymentAccount), PaymentBank: strings.TrimSpace(in.PaymentBank), PaymentQRImg: strings.TrimSpace(in.PaymentQRImg), Type: in.Type, BusinessName: strings.TrimSpace(in.BusinessName), BusinessStoreCategory: in.BusinessStoreCategory, BusinessStoreType: in.BusinessStoreType, Status: status}, nil
}

func itoa(v uint) string { return strconv.FormatUint(uint64(v), 10) }

func joinInts(values []int) string {
	if len(values) == 0 {
		return ""
	}
	parts := make([]string, 0, len(values))
	for _, v := range values {
		parts = append(parts, strconv.Itoa(v))
	}
	return strings.Join(parts, ",")
}

func joinUints(values []uint) string {
	if len(values) == 0 {
		return ""
	}
	parts := make([]string, 0, len(values))
	for _, v := range values {
		parts = append(parts, itoa(v))
	}
	return strings.Join(parts, ",")
}

func ParseIntCSV(raw string) []int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err == nil {
			out = append(out, n)
		}
	}
	return out
}

func ParseUintCSV(raw string) []uint {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]uint, 0, len(parts))
	for _, part := range parts {
		n, err := strconv.ParseUint(strings.TrimSpace(part), 10, 64)
		if err == nil && n > 0 {
			out = append(out, uint(n))
		}
	}
	return out
}
