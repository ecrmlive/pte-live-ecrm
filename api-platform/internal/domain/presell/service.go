package presell

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	List(ctx context.Context, merID *uint, onlyOn bool, page, limit int) ([]ProductPresell, int64, error)
	Get(ctx context.Context, id uint) (*ProductPresell, error)
	GetByProduct(ctx context.Context, productID uint) (*ProductPresell, error)
	Create(ctx context.Context, p *ProductPresell) error
	Update(ctx context.Context, p *ProductPresell) error
	SoftDelete(ctx context.Context, id uint) error
	DecStock(ctx context.Context, id uint, num int) error
	IncStock(ctx context.Context, id uint, num int) error
	IncSeles(ctx context.Context, id uint, num int) error
	LoadProductMeta(ctx context.Context, productID uint) (storeName, image, merName string, price, cost float64, merID uint, err error)

	CreatePresellOrder(ctx context.Context, o *PresellOrder) error
	GetPresellOrder(ctx context.Context, id uint) (*PresellOrder, error)
	GetPresellOrderByOrderID(ctx context.Context, orderID uint) (*PresellOrder, error)
	ListPresellOrdersByUID(ctx context.Context, uid uint, unpaidOnly bool, page, limit int) ([]PresellOrder, int64, error)
	MarkPresellOrderPaid(ctx context.Context, id uint, payType int8, at time.Time) error
	InvalidatePresellOrder(ctx context.Context, id uint) error
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListAdmin(ctx context.Context, merID *uint, page, limit int) (*PageResult[ProductPresell], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.List(ctx, merID, false, page, limit)
	if err != nil {
		return nil, err
	}
	_ = s.enrich(ctx, list)
	return &PageResult[ProductPresell]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListApp(ctx context.Context, page, limit int) (*PageResult[ProductPresell], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.List(ctx, nil, true, page, limit)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	for i := range list {
		list[i].InWindow = inWindow(&list[i], now)
	}
	_ = s.enrich(ctx, list)
	return &PageResult[ProductPresell]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) Get(ctx context.Context, id uint) (*ProductPresell, error) {
	p, err := s.store.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	tmp := []ProductPresell{*p}
	_ = s.enrich(ctx, tmp)
	*p = tmp[0]
	p.InWindow = inWindow(p, time.Now())
	fillFinal(p)
	return p, nil
}

func (s *Service) Create(ctx context.Context, merID uint, in SaveInput) (*ProductPresell, error) {
	if merID == 0 || in.ProductID == 0 || in.Price <= 0 {
		return nil, ErrBadParam
	}
	name, _, _, _, _, pMer, err := s.store.LoadProductMeta(ctx, in.ProductID)
	if err != nil || pMer != merID {
		return nil, ErrBadParam
	}
	title := strings.TrimSpace(in.StoreName)
	if title == "" {
		title = name + " · 预售"
	}
	start, end := parseRange(in.StartTime, in.EndTime)
	if !validActivityRange(in.StartTime, in.EndTime, start, end) {
		return nil, ErrBadParam
	}
	stock := in.Stock
	if stock <= 0 {
		stock = 100
	}
	ptype := in.PresellType
	if ptype != 2 {
		ptype = 1
	}
	finalPrice := in.FinalPrice
	if ptype == 2 {
		if finalPrice <= 0 {
			finalPrice = round2(in.Price - in.DownPrice)
		}
		if err := validateDeposit(in.Price, in.DownPrice, finalPrice); err != nil {
			return nil, err
		}
	}
	p := &ProductPresell{
		StartTime: start, EndTime: end, Status: 1, PresellType: ptype,
		PayCount: in.PayCount, DeliveryType: 1, DeliveryDay: in.DeliveryDay,
		ProductID: in.ProductID, Price: in.Price, DownPrice: in.DownPrice, FinalPrice: finalPrice,
		Stock: stock, IsShow: 1, StoreName: title,
		MerID: merID, StoreInfo: strings.TrimSpace(in.StoreInfo),
		CreateTime: time.Now(), ProductStatus: 1, ActionStatus: 1,
		FinalStartTime: strings.TrimSpace(in.FinalStartTime),
		FinalEndTime:   strings.TrimSpace(in.FinalEndTime),
	}
	if ptype == 2 && p.FinalStartTime == "" {
		p.FinalStartTime = time.Now().Format("2006-01-02 15:04:05")
	}
	if ptype == 2 && p.FinalEndTime == "" {
		p.FinalEndTime = time.Now().AddDate(0, 0, 30).Format("2006-01-02 15:04:05")
	}
	if ptype == 2 && !validFinalRange(p.FinalStartTime, p.FinalEndTime) {
		return nil, ErrBadParam
	}
	if in.DeliveryType > 0 {
		p.DeliveryType = in.DeliveryType
	}
	if in.IsShow != nil {
		p.IsShow = *in.IsShow
	}
	if in.Status != nil {
		p.Status = *in.Status
	}
	if err := s.store.Create(ctx, p); err != nil {
		return nil, err
	}
	return s.Get(ctx, p.ProductPresellID)
}

func (s *Service) Update(ctx context.Context, merID, id uint, in SaveInput) (*ProductPresell, error) {
	if in.Status != nil && *in.Status != 0 && *in.Status != 1 {
		return nil, ErrBadParam
	}
	p, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if merID > 0 && p.MerID != merID {
		return nil, ErrNotFound
	}
	if name := strings.TrimSpace(in.StoreName); name != "" {
		p.StoreName = name
	}
	if info := strings.TrimSpace(in.StoreInfo); info != "" {
		p.StoreInfo = info
	}
	if in.Price > 0 {
		p.Price = in.Price
	}
	if in.DownPrice > 0 {
		p.DownPrice = in.DownPrice
	}
	if in.FinalPrice > 0 {
		p.FinalPrice = in.FinalPrice
	}
	if in.PresellType == 1 || in.PresellType == 2 {
		p.PresellType = in.PresellType
	}
	if in.Stock > 0 {
		p.Stock = in.Stock
	}
	if in.StartTime != "" {
		t, ok := parseTime(in.StartTime)
		if !ok {
			return nil, ErrBadParam
		}
		p.StartTime = t
	}
	if in.EndTime != "" {
		t, ok := parseTime(in.EndTime)
		if !ok {
			return nil, ErrBadParam
		}
		p.EndTime = t
	}
	if strings.TrimSpace(in.FinalStartTime) != "" {
		p.FinalStartTime = strings.TrimSpace(in.FinalStartTime)
	}
	if strings.TrimSpace(in.FinalEndTime) != "" {
		p.FinalEndTime = strings.TrimSpace(in.FinalEndTime)
	}
	if in.PayCount >= 0 {
		p.PayCount = in.PayCount
	}
	if in.DeliveryDay >= 0 {
		p.DeliveryDay = in.DeliveryDay
	}
	if in.DeliveryType > 0 {
		p.DeliveryType = in.DeliveryType
	}
	if in.IsShow != nil {
		p.IsShow = *in.IsShow
	}
	if in.Status != nil {
		p.Status = *in.Status
	}
	if p.PresellType == 2 {
		if !validFinalRange(p.FinalStartTime, p.FinalEndTime) {
			return nil, ErrBadParam
		}
		if p.FinalPrice <= 0 {
			p.FinalPrice = round2(p.Price - p.DownPrice)
		}
		if err := validateDeposit(p.Price, p.DownPrice, p.FinalPrice); err != nil {
			return nil, err
		}
	}
	if p.EndTime.Before(p.StartTime) {
		return nil, ErrBadParam
	}
	if err := s.store.Update(ctx, p); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// SetShow 预售活动上下架（is_show 1上架 / 0下架）。
func (s *Service) SetShow(ctx context.Context, merID, id uint, show int) (*ProductPresell, error) {
	if show != 0 && show != 1 {
		return nil, ErrBadParam
	}
	v := show
	return s.Update(ctx, merID, id, SaveInput{IsShow: &v})
}

func (s *Service) Delete(ctx context.Context, merID, id uint) error {
	p, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if merID > 0 && p.MerID != merID {
		return ErrNotFound
	}
	return s.store.SoftDelete(ctx, id)
}

// QuoteByProduct 购物车路径：按商品查进行中的全款预售。
func (s *Service) QuoteByProduct(ctx context.Context, productID uint) (*ProductPresell, error) {
	p, err := s.store.GetByProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if p.PresellType != 1 {
		return nil, ErrNotFullPay
	}
	if !inWindow(p, time.Now()) {
		return nil, ErrInactive
	}
	tmp := []ProductPresell{*p}
	_ = s.enrich(ctx, tmp)
	*p = tmp[0]
	return p, nil
}

// Quote 活动窗口内报价（全款/定金）。
func (s *Service) Quote(ctx context.Context, productPresellID uint) (*ProductPresell, error) {
	p, err := s.Get(ctx, productPresellID)
	if err != nil {
		return nil, err
	}
	if p.PresellType != 1 && p.PresellType != 2 {
		return nil, ErrBadParam
	}
	if !inWindow(p, time.Now()) {
		return nil, ErrInactive
	}
	if p.Price <= 0 {
		return nil, ErrBadParam
	}
	fillFinal(p)
	if p.PresellType == 2 {
		if err := validateDeposit(p.Price, p.DownPrice, p.FinalPrice); err != nil {
			return nil, err
		}
	}
	if p.Stock <= 0 {
		return nil, ErrSoldOut
	}
	return p, nil
}

func (s *Service) ProductCost(ctx context.Context, productID uint) (float64, error) {
	_, _, _, _, cost, _, err := s.store.LoadProductMeta(ctx, productID)
	return cost, err
}

func (s *Service) ReserveStock(ctx context.Context, productPresellID uint, num int) error {
	if num <= 0 {
		return ErrBadParam
	}
	p, err := s.Quote(ctx, productPresellID)
	if err != nil {
		return err
	}
	if p.Stock < num {
		return ErrSoldOut
	}
	return s.store.DecStock(ctx, productPresellID, num)
}

func (s *Service) RestoreStock(ctx context.Context, productPresellID uint, num int) error {
	if num <= 0 {
		return ErrBadParam
	}
	return s.store.IncStock(ctx, productPresellID, num)
}

func (s *Service) OnOrderPaid(ctx context.Context, productPresellID uint, num int) error {
	if num <= 0 {
		return nil
	}
	return s.store.IncSeles(ctx, productPresellID, num)
}

// OnDepositPaid 定金支付成功：写尾款单（幂等）。
func (s *Service) OnDepositPaid(ctx context.Context, orderID, productPresellID, uid, merID uint, num int) (*PresellOrder, error) {
	p, err := s.Get(ctx, productPresellID)
	if err != nil {
		return nil, err
	}
	if p.PresellType != 2 {
		return nil, nil
	}
	if existing, err := s.store.GetPresellOrderByOrderID(ctx, orderID); err == nil && existing != nil {
		return existing, nil
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	fillFinal(p)
	finalPay := round2(p.FinalPrice * float64(num))
	if finalPay <= 0 {
		return nil, ErrBadParam
	}
	fs, fe := parseFinalWindow(p.FinalStartTime, p.FinalEndTime)
	o := &PresellOrder{
		PresellOrderSN:   genSN("PF"),
		UID:              uid,
		MerID:            merID,
		OrderID:          orderID,
		ProductPresellID: productPresellID,
		FinalStartTime:   fs,
		FinalEndTime:     fe,
		Paid:             0,
		Status:           1,
		PayPrice:         finalPay,
		CreateTime:       time.Now(),
	}
	if err := s.store.CreatePresellOrder(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Service) ListFinalOrders(ctx context.Context, uid uint, unpaidOnly bool, page, limit int) (*PageResult[PresellOrder], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListPresellOrdersByUID(ctx, uid, unpaidOnly, page, limit)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if p, err := s.Get(ctx, list[i].ProductPresellID); err == nil {
			list[i].StoreName = p.StoreName
		}
	}
	return &PageResult[PresellOrder]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetFinalOrder(ctx context.Context, id, uid uint) (*PresellOrder, error) {
	o, err := s.store.GetPresellOrder(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if uid > 0 && o.UID != uid {
		return nil, ErrNotFound
	}
	if p, err := s.Get(ctx, o.ProductPresellID); err == nil {
		o.StoreName = p.StoreName
	}
	return o, nil
}

func (s *Service) MarkFinalTimeout(ctx context.Context, o *PresellOrder) error {
	if o == nil || o.Paid == 1 {
		return nil
	}
	return s.store.InvalidatePresellOrder(ctx, o.PresellOrderID)
}

func (s *Service) MarkFinalPaid(ctx context.Context, id uint, payType int8) error {
	return s.store.MarkPresellOrderPaid(ctx, id, payType, time.Now())
}

func (s *Service) enrich(ctx context.Context, list []ProductPresell) error {
	for i := range list {
		_, img, merName, ot, _, _, err := s.store.LoadProductMeta(ctx, list[i].ProductID)
		if err != nil {
			continue
		}
		list[i].Image = img
		list[i].MerName = merName
		list[i].OtPrice = ot
		fillFinal(&list[i])
	}
	return nil
}

func fillFinal(p *ProductPresell) {
	if p == nil || p.PresellType != 2 {
		return
	}
	if p.FinalPrice <= 0 && p.Price > p.DownPrice {
		p.FinalPrice = round2(p.Price - p.DownPrice)
	}
}

func validateDeposit(price, down, final float64) error {
	if down <= 0 || final <= 0 || down >= price {
		return ErrBadParam
	}
	diff := down + final - price
	if diff < 0 {
		diff = -diff
	}
	if diff > 0.02 {
		return ErrBadParam
	}
	return nil
}

func parseFinalWindow(startS, endS string) (time.Time, time.Time) {
	now := time.Now()
	fs, fe := now, now.AddDate(0, 0, 30)
	if t, ok := parseTime(startS); ok {
		fs = t
	}
	if t, ok := parseTime(endS); ok {
		fe = t
	}
	if !fe.After(fs) {
		fe = fs.AddDate(0, 0, 30)
	}
	return fs, fe
}

func round2(v float64) float64 {
	return float64(int64(v*100+0.5)) / 100
}

func genSN(prefix string) string {
	return prefix + time.Now().Format("20060102150405") + strings.ReplaceAll(time.Now().Format(".000"), ".", "")
}

func inWindow(p *ProductPresell, now time.Time) bool {
	if p == nil || p.IsDel != 0 || p.Status != 1 || p.IsShow != 1 || p.ProductStatus != 1 || p.ActionStatus != 1 {
		return false
	}
	if now.Before(p.StartTime) || now.After(p.EndTime) {
		return false
	}
	return true
}

func parseRange(startS, endS string) (time.Time, time.Time) {
	now := time.Now()
	start := now
	end := now.AddDate(0, 0, 30)
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

// validActivityRange 保留空时间的创建默认值，但拒绝非法输入和倒置时间窗。
func validActivityRange(startS, endS string, start, end time.Time) bool {
	startRaw, endRaw := strings.TrimSpace(startS), strings.TrimSpace(endS)
	if startRaw != "" && !validTime(startRaw) || endRaw != "" && !validTime(endRaw) {
		return false
	}
	if startRaw != "" && endRaw != "" {
		rawStart, _ := parseTime(startRaw)
		rawEnd, _ := parseTime(endRaw)
		if rawEnd.Before(rawStart) {
			return false
		}
	}
	return !end.Before(start)
}

func validFinalRange(startS, endS string) bool {
	start, ok := parseTime(startS)
	if !ok {
		return false
	}
	end, ok := parseTime(endS)
	return ok && !end.Before(start)
}

func validTime(s string) bool {
	_, ok := parseTime(s)
	return ok
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
