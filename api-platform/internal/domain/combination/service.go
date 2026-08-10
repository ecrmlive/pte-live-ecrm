package combination

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	ListGroups(ctx context.Context, merID *uint, onlyOn bool, page, limit int) ([]ProductGroup, int64, error)
	GetGroup(ctx context.Context, id uint) (*ProductGroup, error)
	CreateGroup(ctx context.Context, g *ProductGroup) error
	UpdateGroup(ctx context.Context, g *ProductGroup) error
	SoftDeleteGroup(ctx context.Context, id uint) error
	LoadProductMeta(ctx context.Context, productID uint) (storeName, image, merName string, price, cost float64, merID uint, err error)

	CreateBuying(ctx context.Context, b *Buying) error
	GetBuying(ctx context.Context, id uint) (*Buying, error)
	UpdateBuying(ctx context.Context, b *Buying) error
	ListOpenBuyings(ctx context.Context, productGroupID uint, limit int) ([]Buying, error)
	ListBuyingsAdmin(ctx context.Context, q AdminBuyingQuery) ([]Buying, int64, error)

	CreateMember(ctx context.Context, m *Member) error
	GetMemberByOrder(ctx context.Context, orderID uint) (*Member, error)
	FindMember(ctx context.Context, buyingID, uid uint) (*Member, error)
	ListMembers(ctx context.Context, buyingID uint) ([]Member, error)
	UpdateMember(ctx context.Context, m *Member) error
	SoftDeleteMember(ctx context.Context, id uint) error
	SoftDeleteBuying(ctx context.Context, id uint) error
	ListOrderIDsByBuying(ctx context.Context, buyingID uint) ([]uint, error)
	BumpGroupSuccess(ctx context.Context, productGroupID uint) error
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListAdmin(ctx context.Context, merID *uint, page, limit int) (*PageResult[ProductGroup], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListGroups(ctx, merID, false, page, limit)
	if err != nil {
		return nil, err
	}
	_ = s.enrich(ctx, list)
	return &PageResult[ProductGroup]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListApp(ctx context.Context, page, limit int) (*PageResult[ProductGroup], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListGroups(ctx, nil, true, page, limit)
	if err != nil {
		return nil, err
	}
	_ = s.enrich(ctx, list)
	return &PageResult[ProductGroup]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) Get(ctx context.Context, id uint) (*ProductGroup, error) {
	g, err := s.store.GetGroup(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	tmp := []ProductGroup{*g}
	_ = s.enrich(ctx, tmp)
	*g = tmp[0]
	return g, nil
}

func (s *Service) Create(ctx context.Context, merID uint, in SaveInput) (*ProductGroup, error) {
	if merID == 0 || in.ProductID == 0 || in.Price <= 0 || in.BuyingCountNum < 2 {
		return nil, ErrBadParam
	}
	_, _, _, _, _, pMer, err := s.store.LoadProductMeta(ctx, in.ProductID)
	if err != nil {
		return nil, ErrBadParam
	}
	if pMer != merID {
		return nil, ErrBadParam
	}
	g := &ProductGroup{
		ProductID: in.ProductID, MerID: merID, Price: in.Price,
		BuyingCountNum: in.BuyingCountNum, BuyingNum: 1, OncePayCount: 1,
		Time: in.Time, Status: 1, IsShow: 1, ActionStatus: 1, ProductStatus: 1,
		CreateTime: time.Now(),
	}
	if g.Time <= 0 {
		g.Time = 24
	}
	g.StartTime, g.EndTime = parseRange(in.StartTime, in.EndTime)
	if !validActivityRange(in.StartTime, in.EndTime, g.StartTime, g.EndTime) {
		return nil, ErrBadParam
	}
	if in.IsShow != nil {
		g.IsShow = *in.IsShow
	}
	if in.Status != nil {
		g.Status = *in.Status
	}
	if err := s.store.CreateGroup(ctx, g); err != nil {
		return nil, err
	}
	return s.Get(ctx, g.ProductGroupID)
}

func (s *Service) Update(ctx context.Context, merID, id uint, in SaveInput) (*ProductGroup, error) {
	if in.Status != nil && *in.Status != 0 && *in.Status != 1 {
		return nil, ErrBadParam
	}
	g, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if merID > 0 && g.MerID != merID {
		return nil, ErrNotFound
	}
	if in.Price > 0 {
		g.Price = in.Price
	}
	if in.BuyingCountNum >= 2 {
		g.BuyingCountNum = in.BuyingCountNum
	}
	if in.Time > 0 {
		g.Time = in.Time
	}
	if in.StartTime != "" || in.EndTime != "" {
		if (in.StartTime != "" && !validActivityTime(in.StartTime)) || (in.EndTime != "" && !validActivityTime(in.EndTime)) {
			return nil, ErrBadParam
		}
		st, et := parseRange(in.StartTime, in.EndTime)
		if in.StartTime != "" {
			g.StartTime = st
		}
		if in.EndTime != "" {
			g.EndTime = et
		}
	}
	if g.EndTime.Before(g.StartTime) {
		return nil, ErrBadParam
	}
	if in.IsShow != nil {
		g.IsShow = *in.IsShow
	}
	if in.Status != nil {
		g.Status = *in.Status
	}
	if in.ProductStatus != nil {
		ps := *in.ProductStatus
		if ps != 0 && ps != 1 && ps != -1 && ps != -2 {
			return nil, ErrBadParam
		}
		g.ProductStatus = ps
	}
	if refusal := strings.TrimSpace(in.Refusal); refusal != "" {
		g.Refusal = refusal
	}
	if err := s.store.UpdateGroup(ctx, g); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// SetShow 拼团活动上下架（is_show 1上架 / 0下架）。
func (s *Service) SetShow(ctx context.Context, merID, id uint, show int) (*ProductGroup, error) {
	if show != 0 && show != 1 {
		return nil, ErrBadParam
	}
	v := show
	return s.Update(ctx, merID, id, SaveInput{IsShow: &v})
}

// Audit 平台审核拼团商品：1 通过 / -1 拒绝。
func (s *Service) Audit(ctx context.Context, id uint, status int, refusal string) (*ProductGroup, error) {
	if status != 1 && status != -1 {
		return nil, ErrBadParam
	}
	refusal = strings.TrimSpace(refusal)
	if status == -1 && refusal == "" {
		return nil, ErrBadParam
	}
	show := 0
	if status == 1 {
		show = 1
	}
	return s.Update(ctx, 0, id, SaveInput{ProductStatus: &status, IsShow: &show, Refusal: refusal})
}

// ForceOff 强制下架：product_status=-2，隐藏前台展示。
func (s *Service) ForceOff(ctx context.Context, ids []uint, reason string) error {
	reason = strings.TrimSpace(reason)
	if len(ids) == 0 || reason == "" {
		return ErrBadParam
	}
	off := -2
	show := 0
	st := 0
	for _, id := range ids {
		if _, err := s.Update(ctx, 0, id, SaveInput{
			ProductStatus: &off,
			IsShow:        &show,
			Status:        &st,
			Refusal:       reason,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, merID, id uint) error {
	g, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if merID > 0 && g.MerID != merID {
		return ErrNotFound
	}
	return s.store.SoftDeleteGroup(ctx, id)
}

// Quote 活动期内拼团价。
func (s *Service) Quote(ctx context.Context, productGroupID uint) (*ProductGroup, error) {
	g, err := s.Get(ctx, productGroupID)
	if err != nil {
		return nil, err
	}
	if !activeNow(g) {
		return nil, ErrInactive
	}
	return g, nil
}

func (s *Service) ListBuyings(ctx context.Context, productGroupID uint, limit int) ([]Buying, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	rows, err := s.store.ListOpenBuyings(ctx, productGroupID, limit)
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	out := make([]Buying, 0, len(rows))
	for i := range rows {
		if rows[i].EndTime > 0 && rows[i].EndTime < now {
			continue
		}
		rows[i].Remain = rows[i].BuyingCountNum - rows[i].YetBuyingNum
		if rows[i].Remain < 0 {
			rows[i].Remain = 0
		}
		ms, _ := s.store.ListMembers(ctx, rows[i].GroupBuyingID)
		rows[i].Members = ms
		out = append(out, rows[i])
	}
	return out, nil
}

func (s *Service) GetBuying(ctx context.Context, id uint) (*Buying, error) {
	b, err := s.store.GetBuying(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	b.Remain = b.BuyingCountNum - b.YetBuyingNum
	if b.Remain < 0 {
		b.Remain = 0
	}
	ms, _ := s.store.ListMembers(ctx, id)
	b.Members = ms
	tmp := []Buying{*b}
	_ = s.enrichBuyings(ctx, tmp)
	*b = tmp[0]
	b.Members = ms
	return b, nil
}

// ListAdminBuyings 平台监管：用户开团记录列表（拼团活动列表）。
func (s *Service) ListAdminBuyings(ctx context.Context, q AdminBuyingQuery) (*PageResult[Buying], error) {
	q.Page, q.Limit = normalize(q.Page, q.Limit)
	list, total, err := s.store.ListBuyingsAdmin(ctx, q)
	if err != nil {
		return nil, err
	}
	_ = s.enrichBuyings(ctx, list)
	return &PageResult[Buying]{List: list, Total: total, Page: q.Page, Limit: q.Limit}, nil
}

func (s *Service) enrichBuyings(ctx context.Context, list []Buying) error {
	for i := range list {
		list[i].Remain = list[i].BuyingCountNum - list[i].YetBuyingNum
		if list[i].Remain < 0 {
			list[i].Remain = 0
		}
		list[i].StatusText = buyingStatusText(list[i].Status)
		if list[i].EndTime > 0 {
			loc := time.Local
			if shanghai, err := time.LoadLocation("Asia/Shanghai"); err == nil {
				loc = shanghai
			}
			list[i].StopTime = time.Unix(list[i].EndTime, 0).In(loc).Format("2006-01-02 15:04:05")
		}
		g, err := s.store.GetGroup(ctx, list[i].ProductGroupID)
		if err == nil && g != nil {
			list[i].ProductID = g.ProductID
			list[i].Price = g.Price
			if list[i].MerID == 0 {
				list[i].MerID = g.MerID
			}
			name, img, merName, _, _, _, metaErr := s.store.LoadProductMeta(ctx, g.ProductID)
			if metaErr == nil {
				list[i].StoreName = name
				list[i].Image = img
				if list[i].MerName == "" {
					list[i].MerName = merName
				}
			}
		}
		if list[i].Nickname == "" {
			ms, _ := s.store.ListMembers(ctx, list[i].GroupBuyingID)
			for _, m := range ms {
				if m.IsInitiator == 1 || m.IsLeader == 1 {
					list[i].Nickname = m.Nickname
					list[i].UID = m.UID
					break
				}
			}
			if list[i].Nickname == "" && len(ms) > 0 {
				list[i].Nickname = ms[0].Nickname
				list[i].UID = ms[0].UID
			}
		}
	}
	return nil
}

func buyingStatusText(status int) string {
	switch status {
	case BuyingStatusDone:
		return "已完成"
	case BuyingStatusFailed:
		return "已失败"
	default:
		return "未完成"
	}
}

// BeginJoin 开团或参团占位（待支付）；返回 buyingID。
func (s *Service) BeginJoin(ctx context.Context, uid, productGroupID, joinBuyingID uint, nickname string) (buyingID uint, isLeader bool, err error) {
	g, err := s.Quote(ctx, productGroupID)
	if err != nil {
		return 0, false, err
	}
	if nickname == "" {
		nickname = fmt.Sprintf("用户%d", uid)
	}
	if joinBuyingID == 0 {
		end := time.Now().Add(time.Duration(g.Time) * time.Hour).Unix()
		b := &Buying{
			ProductGroupID: g.ProductGroupID, Status: 0,
			BuyingCountNum: g.BuyingCountNum, BuyingNum: g.BuyingNum,
			YetBuyingNum: 0, MerID: g.MerID, EndTime: end,
		}
		if err := s.store.CreateBuying(ctx, b); err != nil {
			return 0, false, err
		}
		return b.GroupBuyingID, true, nil
	}
	b, err := s.GetBuying(ctx, joinBuyingID)
	if err != nil {
		return 0, false, err
	}
	if b.ProductGroupID != productGroupID {
		return 0, false, ErrBadParam
	}
	if b.Status != 0 {
		return 0, false, ErrBuyingClosed
	}
	if b.EndTime > 0 && b.EndTime < time.Now().Unix() {
		return 0, false, ErrBuyingClosed
	}
	// 未支付成员也占席：以有效成员数与已支付人数取大，避免超员占位
	ms, err := s.store.ListMembers(ctx, joinBuyingID)
	if err != nil {
		return 0, false, err
	}
	seats := len(ms)
	if int(b.YetBuyingNum) > seats {
		seats = int(b.YetBuyingNum)
	}
	if seats >= int(b.BuyingCountNum) {
		return 0, false, ErrBuyingFull
	}
	if _, err := s.store.FindMember(ctx, joinBuyingID, uid); err == nil {
		return 0, false, ErrAlreadyJoined
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, false, err
	}
	return joinBuyingID, false, nil
}

func (s *Service) AttachMember(ctx context.Context, buyingID, productGroupID, uid, orderID uint, isLeader bool, nickname string) error {
	if nickname == "" {
		nickname = fmt.Sprintf("用户%d", uid)
	}
	b, err := s.store.GetBuying(ctx, buyingID)
	if err != nil {
		return err
	}
	if b.Status != 0 || b.IsDel == 1 {
		return ErrBuyingClosed
	}
	ms, err := s.store.ListMembers(ctx, buyingID)
	if err != nil {
		return err
	}
	if len(ms) >= int(b.BuyingCountNum) {
		return ErrBuyingFull
	}
	if _, err := s.store.FindMember(ctx, buyingID, uid); err == nil {
		return ErrAlreadyJoined
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	leader, init := 0, 0
	if isLeader {
		leader, init = 1, 1
	}
	return s.store.CreateMember(ctx, &Member{
		GroupBuyingID: buyingID, ProductGroupID: productGroupID,
		Status: 0, IsInitiator: init, IsLeader: int8(leader),
		OrderID: orderID, UID: uid, Nickname: nickname,
	})
}

// CancelUnpaid 未支付取消：释放团席位；若团无剩余成员则软删团次。
func (s *Service) CancelUnpaid(ctx context.Context, orderID uint) error {
	if orderID == 0 {
		return nil
	}
	m, err := s.store.GetMemberByOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	// 已计入支付人数的成员不回滚席位（不应走未支付取消路径）
	if m.Status == 1 {
		return nil
	}
	if err := s.store.SoftDeleteMember(ctx, m.ID); err != nil {
		return err
	}
	left, err := s.store.ListMembers(ctx, m.GroupBuyingID)
	if err != nil {
		return err
	}
	// 无剩余成员，或团长未支付退出：关闭团次，释放开放席位
	if len(left) > 0 && m.IsLeader != 1 {
		return nil
	}
	b, err := s.store.GetBuying(ctx, m.GroupBuyingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if b.Status != 0 || b.IsDel == 1 {
		return nil
	}
	return s.store.SoftDeleteBuying(ctx, b.GroupBuyingID)
}

// OnOrderPaid 支付成功后计人；满员则成团。
func (s *Service) OnOrderPaid(ctx context.Context, orderID uint) (success bool, orderIDs []uint, err error) {
	m, err := s.store.GetMemberByOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, nil
		}
		return false, nil, err
	}
	if m.Status == 1 {
		b, err := s.GetBuying(ctx, m.GroupBuyingID)
		if err != nil {
			return false, nil, err
		}
		ids, _ := s.store.ListOrderIDsByBuying(ctx, m.GroupBuyingID)
		return b.Status == 10, ids, nil
	}
	m.Status = 1
	if err := s.store.UpdateMember(ctx, m); err != nil {
		return false, nil, err
	}
	b, err := s.store.GetBuying(ctx, m.GroupBuyingID)
	if err != nil {
		return false, nil, err
	}
	b.YetBuyingNum++
	if b.YetBuyingNum >= b.BuyingCountNum {
		b.Status = 10
		success = true
	}
	if err := s.store.UpdateBuying(ctx, b); err != nil {
		return false, nil, err
	}
	if success {
		_ = s.store.BumpGroupSuccess(ctx, b.ProductGroupID)
	}
	ids, err := s.store.ListOrderIDsByBuying(ctx, m.GroupBuyingID)
	return success, ids, err
}

func (s *Service) enrich(ctx context.Context, list []ProductGroup) error {
	for i := range list {
		name, img, merName, price, _, _, err := s.store.LoadProductMeta(ctx, list[i].ProductID)
		if err != nil {
			continue
		}
		list[i].StoreName = name
		list[i].Image = img
		list[i].MerName = merName
		list[i].OtPrice = price
	}
	return nil
}

// ProductCost 下单用成本价。
func (s *Service) ProductCost(ctx context.Context, productID uint) (cost float64, err error) {
	_, _, _, _, cost, _, err = s.store.LoadProductMeta(ctx, productID)
	return cost, err
}

func activeNow(g *ProductGroup) bool {
	if g == nil || g.IsDel == 1 || g.IsShow != 1 || g.Status != 1 || g.ActionStatus != 1 {
		return false
	}
	now := time.Now()
	if !g.StartTime.IsZero() && now.Before(g.StartTime) {
		return false
	}
	if !g.EndTime.IsZero() && now.After(g.EndTime) {
		return false
	}
	return true
}

func parseRange(start, end string) (time.Time, time.Time) {
	now := time.Now()
	st, et := now, now.AddDate(0, 0, 30)
	if start != "" {
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local); err == nil {
			st = t
		} else if t, err := time.ParseInLocation("2006-01-02", start, time.Local); err == nil {
			st = t
		}
	}
	if end != "" {
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local); err == nil {
			et = t
		} else if t, err := time.ParseInLocation("2006-01-02", end, time.Local); err == nil {
			et = t.Add(23*time.Hour + 59*time.Minute)
		}
	}
	return st, et
}

func validActivityRange(start, end string, startTime, endTime time.Time) bool {
	return (start == "" || validActivityTime(start)) && (end == "" || validActivityTime(end)) && !endTime.Before(startTime)
}

func validActivityTime(raw string) bool {
	if _, err := time.ParseInLocation("2006-01-02 15:04:05", raw, time.Local); err == nil {
		return true
	}
	_, err := time.ParseInLocation("2006-01-02", raw, time.Local)
	return err == nil
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
