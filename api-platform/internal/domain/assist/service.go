package assist

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	List(ctx context.Context, merID *uint, onlyOn bool, page, limit int) ([]ProductAssist, int64, error)
	Get(ctx context.Context, id uint) (*ProductAssist, error)
	Create(ctx context.Context, a *ProductAssist) error
	Update(ctx context.Context, a *ProductAssist) error
	SoftDelete(ctx context.Context, id uint) error
	DecStock(ctx context.Context, id uint, num int) error
	IncStock(ctx context.Context, id uint, num int) error
	LoadProductMeta(ctx context.Context, productID uint) (storeName, image, merName string, price, cost float64, merID uint, err error)
	LoadNickname(ctx context.Context, uid uint) (string, error)

	CreateSet(ctx context.Context, s *AssistSet) error
	GetSet(ctx context.Context, id uint) (*AssistSet, error)
	UpdateSet(ctx context.Context, s *AssistSet) error
	ListSetsByAssist(ctx context.Context, assistID uint, onlyOpen bool, limit int) ([]AssistSet, error)
	ListSetsAdmin(ctx context.Context, q AdminSetQuery) ([]AssistSet, int64, error)
	FindOpenSetByUID(ctx context.Context, assistID, uid uint) (*AssistSet, error)

	CreateHelper(ctx context.Context, u *AssistUser) error
	HasHelped(ctx context.Context, setID, uid uint) (bool, error)
	CountHelpsByUID(ctx context.Context, assistID, uid uint) (int64, error)
	ListHelpers(ctx context.Context, setID uint) ([]AssistUser, error)
	CountAssistStats(ctx context.Context, assistID uint) (success, pay, all int, err error)
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListAdmin(ctx context.Context, merID *uint, page, limit int) (*PageResult[ProductAssist], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.List(ctx, merID, false, page, limit)
	if err != nil {
		return nil, err
	}
	_ = s.enrich(ctx, list)
	return &PageResult[ProductAssist]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListApp(ctx context.Context, page, limit int) (*PageResult[ProductAssist], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.List(ctx, nil, true, page, limit)
	if err != nil {
		return nil, err
	}
	_ = s.enrich(ctx, list)
	return &PageResult[ProductAssist]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) Get(ctx context.Context, id uint) (*ProductAssist, error) {
	a, err := s.store.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	tmp := []ProductAssist{*a}
	_ = s.enrich(ctx, tmp)
	*a = tmp[0]
	return a, nil
}

func (s *Service) Create(ctx context.Context, merID uint, in SaveInput) (*ProductAssist, error) {
	if merID == 0 || in.ProductID == 0 || in.AssistPrice <= 0 {
		return nil, ErrBadParam
	}
	name, _, _, _, _, pMer, err := s.store.LoadProductMeta(ctx, in.ProductID)
	if err != nil || pMer != merID {
		return nil, ErrBadParam
	}
	title := strings.TrimSpace(in.StoreName)
	if title == "" {
		title = name + " · 助力"
	}
	start, end := parseRange(in.StartTime, in.EndTime)
	need := in.AssistCount
	if need < 1 {
		need = 1
	}
	per := in.AssistUserCount
	if per < 1 {
		per = 1
	}
	stock := in.Stock
	if stock <= 0 {
		stock = 100
	}
	a := &ProductAssist{
		StartTime: start, EndTime: end, Status: 1, AssistCount: need, AssistUserCount: per,
		ProductID: in.ProductID, AssistPrice: in.AssistPrice, Stock: stock, IsShow: 1,
		StoreName: title, MerID: merID, StoreInfo: strings.TrimSpace(in.StoreInfo),
		CreateTime: time.Now(), ProductStatus: 1, ActionStatus: 1,
	}
	if in.IsShow != nil {
		a.IsShow = int8(*in.IsShow)
	}
	if in.Status != nil {
		a.Status = *in.Status
	}
	if err := s.store.Create(ctx, a); err != nil {
		return nil, err
	}
	return s.Get(ctx, a.ProductAssistID)
}

func (s *Service) Update(ctx context.Context, merID, id uint, in SaveInput) (*ProductAssist, error) {
	a, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if merID > 0 && a.MerID != merID {
		return nil, ErrNotFound
	}
	if name := strings.TrimSpace(in.StoreName); name != "" {
		a.StoreName = name
	}
	if info := strings.TrimSpace(in.StoreInfo); info != "" {
		a.StoreInfo = info
	}
	if in.AssistPrice > 0 {
		a.AssistPrice = in.AssistPrice
	}
	if in.AssistCount > 0 {
		a.AssistCount = in.AssistCount
	}
	if in.AssistUserCount > 0 {
		a.AssistUserCount = in.AssistUserCount
	}
	if in.Stock > 0 {
		a.Stock = in.Stock
	}
	if in.StartTime != "" || in.EndTime != "" {
		start, end := parseRange(in.StartTime, in.EndTime)
		a.StartTime, a.EndTime = start, end
	}
	if in.IsShow != nil {
		a.IsShow = int8(*in.IsShow)
	}
	if in.Status != nil {
		a.Status = *in.Status
	}
	if err := s.store.Update(ctx, a); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// SetShow 助力活动上下架（is_show 1上架 / 0下架）。
func (s *Service) SetShow(ctx context.Context, merID, id uint, show int) (*ProductAssist, error) {
	if show != 0 && show != 1 {
		return nil, ErrBadParam
	}
	v := show
	return s.Update(ctx, merID, id, SaveInput{IsShow: &v})
}

func (s *Service) Delete(ctx context.Context, merID, id uint) error {
	a, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if merID > 0 && a.MerID != merID {
		return ErrNotFound
	}
	return s.store.SoftDelete(ctx, id)
}

// FindMine 返回当前用户进行中或已满员待下单的助力单。
func (s *Service) FindMine(ctx context.Context, uid, assistID uint) (*AssistSet, error) {
	if uid == 0 || assistID == 0 {
		return nil, ErrBadParam
	}
	existing, err := s.store.FindOpenSetByUID(ctx, assistID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s.GetSet(ctx, existing.ProductAssistSetID)
}

func (s *Service) StartSet(ctx context.Context, uid, assistID uint) (*AssistSet, error) {
	if uid == 0 || assistID == 0 {
		return nil, ErrBadParam
	}
	a, err := s.quote(ctx, assistID)
	if err != nil {
		return nil, err
	}
	if existing, err := s.store.FindOpenSetByUID(ctx, assistID, uid); err == nil && existing != nil {
		return s.GetSet(ctx, existing.ProductAssistSetID)
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	set := &AssistSet{
		ProductAssistID: assistID, ProductID: a.ProductID, UID: uid,
		Status: SetStatusRunning, AssistCount: a.AssistCount, AssistUserCount: a.AssistUserCount,
		YetAssistCount: 0, CreateTime: time.Now(), MerID: a.MerID,
	}
	if err := s.store.CreateSet(ctx, set); err != nil {
		return nil, err
	}
	return s.GetSet(ctx, set.ProductAssistSetID)
}

func (s *Service) Help(ctx context.Context, uid, setID uint) (*AssistSet, error) {
	if uid == 0 || setID == 0 {
		return nil, ErrBadParam
	}
	set, err := s.store.GetSet(ctx, setID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if set.Status != SetStatusRunning {
		return nil, ErrSetClosed
	}
	if set.UID == uid {
		return nil, ErrSelfHelp
	}
	a, err := s.quote(ctx, set.ProductAssistID)
	if err != nil {
		return nil, err
	}
	ok, err := s.store.HasHelped(ctx, setID, uid)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, ErrAlreadyHelped
	}
	n, err := s.store.CountHelpsByUID(ctx, set.ProductAssistID, uid)
	if err != nil {
		return nil, err
	}
	if int(n) >= a.AssistUserCount {
		return nil, ErrAlreadyHelped
	}
	nick, _ := s.store.LoadNickname(ctx, uid)
	if err := s.store.CreateHelper(ctx, &AssistUser{
		ProductAssistSetID: setID, ProductAssistID: set.ProductAssistID,
		UID: uid, Nickname: nick, CreateTime: time.Now(),
	}); err != nil {
		return nil, err
	}
	set.YetAssistCount++
	if set.YetAssistCount >= set.AssistCount {
		set.Status = SetStatusDone
	}
	if err := s.store.UpdateSet(ctx, set); err != nil {
		return nil, err
	}
	return s.GetSet(ctx, setID)
}

func (s *Service) GetSet(ctx context.Context, id uint) (*AssistSet, error) {
	set, err := s.store.GetSet(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	tmp := []AssistSet{*set}
	_ = s.enrichSets(ctx, tmp)
	*set = tmp[0]
	helpers, _ := s.store.ListHelpers(ctx, id)
	set.Helpers = helpers
	return set, nil
}

// ListAdminSets 平台监管：用户发起的助力实例列表。
func (s *Service) ListAdminSets(ctx context.Context, q AdminSetQuery) (*PageResult[AssistSet], error) {
	q.Page, q.Limit = normalize(q.Page, q.Limit)
	list, total, err := s.store.ListSetsAdmin(ctx, q)
	if err != nil {
		return nil, err
	}
	_ = s.enrichSets(ctx, list)
	return &PageResult[AssistSet]{List: list, Total: total, Page: q.Page, Limit: q.Limit}, nil
}

func (s *Service) enrichSets(ctx context.Context, list []AssistSet) error {
	for i := range list {
		if a, err := s.store.Get(ctx, list[i].ProductAssistID); err == nil && a != nil {
			if strings.TrimSpace(a.StoreName) != "" {
				list[i].StoreName = a.StoreName
			}
			list[i].AssistPrice = a.AssistPrice
			list[i].StartTime = a.StartTime
			list[i].EndTime = a.EndTime
			if list[i].MerID == 0 {
				list[i].MerID = a.MerID
			}
		}
		prodName, img, merName, _, _, _, err := s.store.LoadProductMeta(ctx, list[i].ProductID)
		if err == nil {
			list[i].Image = img
			if list[i].MerName == "" {
				list[i].MerName = merName
			}
			if list[i].StoreName == "" {
				list[i].StoreName = prodName
			}
		}
		if nick, err := s.store.LoadNickname(ctx, list[i].UID); err == nil {
			list[i].Nickname = nick
		}
	}
	return nil
}

func (s *Service) ListSets(ctx context.Context, assistID uint, limit int) ([]AssistSet, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	// 进行中 + 已满员待下单（不含已支付）
	list, err := s.store.ListSetsByAssist(ctx, assistID, false, limit)
	if err != nil {
		return nil, err
	}
	out := make([]AssistSet, 0, len(list))
	for i := range list {
		if list[i].Status != SetStatusRunning && list[i].Status != SetStatusDone {
			continue
		}
		if nick, err := s.store.LoadNickname(ctx, list[i].UID); err == nil {
			list[i].Nickname = nick
		}
		out = append(out, list[i])
	}
	return out, nil
}

// QuoteForOrder 仅 status=10 的助力单可下单。
func (s *Service) QuoteForOrder(ctx context.Context, setID, uid uint) (price float64, productID, merID, assistID uint, storeName, image, merName string, stock int, err error) {
	set, err := s.store.GetSet(ctx, setID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, 0, 0, "", "", "", 0, ErrNotFound
		}
		return 0, 0, 0, 0, "", "", "", 0, err
	}
	if set.UID != uid {
		return 0, 0, 0, 0, "", "", "", 0, ErrForbidden
	}
	if set.Status != SetStatusDone {
		return 0, 0, 0, 0, "", "", "", 0, ErrSetNotOpen
	}
	a, err := s.quote(ctx, set.ProductAssistID)
	if err != nil {
		return 0, 0, 0, 0, "", "", "", 0, err
	}
	return a.AssistPrice, a.ProductID, a.MerID, a.ProductAssistID, a.StoreName, a.Image, a.MerName, a.Stock, nil
}

func (s *Service) ProductCost(ctx context.Context, productID uint) (float64, error) {
	_, _, _, _, cost, _, err := s.store.LoadProductMeta(ctx, productID)
	return cost, err
}

func (s *Service) ReserveStock(ctx context.Context, assistID uint, num int) error {
	if num <= 0 {
		return ErrBadParam
	}
	a, err := s.quote(ctx, assistID)
	if err != nil {
		return err
	}
	if a.Stock < num {
		return ErrSoldOut
	}
	return s.store.DecStock(ctx, assistID, num)
}

// RestoreStock 取消/超时未支付时归还活动库存。
func (s *Service) RestoreStock(ctx context.Context, assistID uint, num int) error {
	if assistID == 0 || num <= 0 {
		return nil
	}
	return s.store.IncStock(ctx, assistID, num)
}

func (s *Service) MarkSetPaid(ctx context.Context, setID uint) error {
	set, err := s.store.GetSet(ctx, setID)
	if err != nil {
		return err
	}
	if set.Status == SetStatusPaid {
		return nil
	}
	set.Status = SetStatusPaid
	return s.store.UpdateSet(ctx, set)
}

func (s *Service) quote(ctx context.Context, id uint) (*ProductAssist, error) {
	a, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if a.Status != 1 || a.IsShow != 1 || a.ProductStatus != 1 || a.ActionStatus != 1 {
		return nil, ErrInactive
	}
	now := time.Now()
	if now.Before(a.StartTime) || now.After(a.EndTime) {
		return nil, ErrInactive
	}
	if a.AssistPrice <= 0 {
		return nil, ErrBadParam
	}
	if a.Stock <= 0 {
		return nil, ErrSoldOut
	}
	return a, nil
}

func (s *Service) enrich(ctx context.Context, list []ProductAssist) error {
	now := time.Now()
	for i := range list {
		_, img, merName, ot, _, _, err := s.store.LoadProductMeta(ctx, list[i].ProductID)
		if err == nil {
			list[i].Image = img
			if list[i].MerName == "" {
				list[i].MerName = merName
			}
			list[i].OtPrice = ot
		}
		success, pay, all, _ := s.store.CountAssistStats(ctx, list[i].ProductAssistID)
		list[i].Success = success
		list[i].Pay = pay
		list[i].All = all
		list[i].StockCount = list[i].Stock + pay
		list[i].AssistStatus = computeAssistStatus(&list[i], now)
		list[i].AssistStatusText = assistStatusText(list[i].AssistStatus)
		list[i].ProductStatusName = productStatusName(&list[i])
	}
	return nil
}

// computeAssistStatus 对齐 CRMEB ProductAssist.getAssistStatusAttr：0未开始 1进行中 2已结束。
func computeAssistStatus(a *ProductAssist, now time.Time) int8 {
	if a == nil {
		return 2
	}
	if a.ActionStatus == -1 || !now.Before(a.EndTime) {
		return 2
	}
	if now.Before(a.StartTime) {
		return 0
	}
	if a.ProductStatus != 1 || a.Status != 1 || a.IsShow != 1 {
		return 0
	}
	return 1
}

func assistStatusText(status int8) string {
	switch status {
	case 0:
		return "未开始"
	case 1:
		return "进行中"
	default:
		return "已结束"
	}
}

func productStatusName(a *ProductAssist) string {
	if a == nil {
		return "未知"
	}
	switch a.ProductStatus {
	case 0:
		return "待审核"
	case -1:
		return "审核未通过"
	case -2:
		return "强制下架"
	case 1:
		if a.IsShow == 1 {
			return "出售中"
		}
		return "仓库中"
	default:
		return "未知"
	}
}

func parseRange(startS, endS string) (time.Time, time.Time) {
	now := time.Now()
	start, end := now, now.AddDate(0, 0, 30)
	if t, ok := parseTime(startS); ok {
		start = t
	}
	if t, ok := parseTime(endS); ok {
		end = t
	}
	if !end.After(start) {
		end = start.AddDate(0, 0, 30)
	}
	return start, end
}

func parseTime(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05", "2006-01-02"} {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func normalize(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}
