package cart

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ProductView struct {
	ProductID     uint
	MerID         uint
	StoreName     string
	MerName       string
	Image         string
	Price         float64
	OtPrice       float64
	Cost          float64
	Stock         uint
	Status        int8
	IsShow        uint8
	MerStatus     int8
	Unique        string
	ProductType   uint8
	Integral      int // 积分商城兑换积分（product.integral）
	GoodsType     uint8
	SvipPriceType int8
	SvipPrice     float64
	MerSvipStatus int8
}

type Store interface {
	FindActive(ctx context.Context, uid uint, productID uint, unique string) (*Cart, error)
	Create(ctx context.Context, c *Cart) error
	UpdateNum(ctx context.Context, cartID uint64, uid uint, num uint) error
	SoftDelete(ctx context.Context, cartID uint64, uid uint) error
	ListByUID(ctx context.Context, uid uint) ([]Cart, error)
	GetByID(ctx context.Context, cartID uint64, uid uint) (*Cart, error)
	ListByIDs(ctx context.Context, uid uint, ids []uint64) ([]Cart, error)
	MarkPaid(ctx context.Context, ids []uint64) error
	ResolveProduct(ctx context.Context, productID uint, unique string) (*ProductView, error)

	ListAddresses(ctx context.Context, uid uint) ([]Address, error)
	GetAddress(ctx context.Context, id, uid uint) (*Address, error)
	CreateAddress(ctx context.Context, a *Address) error
	UpdateAddress(ctx context.Context, a *Address) error
	ClearDefaultAddress(ctx context.Context, uid uint) error
	SoftDeleteAddress(ctx context.Context, id, uid uint) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) Add(ctx context.Context, uid uint, in AddInput) (*Cart, error) {
	if uid == 0 || in.ProductID == 0 {
		return nil, ErrProductOff
	}
	if in.CartNum == 0 {
		in.CartNum = 1
	}
	unique := strings.TrimSpace(in.ProductAttrUnique)
	pv, err := s.store.ResolveProduct(ctx, in.ProductID, unique)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductOff
		}
		return nil, err
	}
	if pv.Status != 1 || pv.IsShow != 1 || pv.MerStatus != 1 {
		return nil, ErrProductOff
	}
	if unique == "" {
		unique = pv.Unique
	}
	if pv.Stock < in.CartNum {
		return nil, ErrStockNotEnough
	}
	if in.IsNew != 1 {
		existing, err := s.store.FindActive(ctx, uid, in.ProductID, unique)
		if err == nil && existing != nil {
			newNum := existing.CartNum + in.CartNum
			if pv.Stock < newNum {
				return nil, ErrStockNotEnough
			}
			if err := s.store.UpdateNum(ctx, existing.CartID, uid, newNum); err != nil {
				return nil, err
			}
			existing.CartNum = newNum
			return s.enrichOne(ctx, existing)
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	c := &Cart{
		UID: uid, MerID: pv.MerID, ProductID: in.ProductID,
		ProductAttrUnique: unique, CartNum: in.CartNum, IsNew: in.IsNew, CreateTime: time.Now(),
	}
	if err := s.store.Create(ctx, c); err != nil {
		return nil, err
	}
	return s.enrichOne(ctx, c)
}

func (s *Service) List(ctx context.Context, uid uint) ([]MerchantBucket, error) {
	rows, err := s.store.ListByUID(ctx, uid)
	if err != nil {
		return nil, err
	}
	buckets := map[uint]*MerchantBucket{}
	order := make([]uint, 0)
	for i := range rows {
		item, err := s.enrichOne(ctx, &rows[i])
		if err != nil {
			rows[i].IsFail = 1
			item = &rows[i]
		}
		b := buckets[item.MerID]
		if b == nil {
			b = &MerchantBucket{MerID: item.MerID, MerName: item.MerName, Items: []Cart{}}
			buckets[item.MerID] = b
			order = append(order, item.MerID)
		}
		if b.MerName == "" {
			b.MerName = item.MerName
		}
		b.Items = append(b.Items, *item)
		if item.IsFail == 0 {
			b.Subtotal += item.Price * float64(item.CartNum)
		}
	}
	out := make([]MerchantBucket, 0, len(order))
	for _, id := range order {
		out = append(out, *buckets[id])
	}
	return out, nil
}

func (s *Service) SetNum(ctx context.Context, uid uint, cartID uint64, num uint) error {
	if num == 0 {
		return ErrInvalidNum
	}
	c, err := s.store.GetByID(ctx, cartID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if c.IsPay == 1 {
		return ErrForbidden
	}
	pv, err := s.store.ResolveProduct(ctx, c.ProductID, c.ProductAttrUnique)
	if err != nil {
		return ErrProductOff
	}
	if pv.Stock < num {
		return ErrStockNotEnough
	}
	return s.store.UpdateNum(ctx, cartID, uid, num)
}

func (s *Service) Delete(ctx context.Context, uid uint, cartID uint64) error {
	if _, err := s.store.GetByID(ctx, cartID, uid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.SoftDelete(ctx, cartID, uid)
}

func (s *Service) MarkPaid(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	return s.store.MarkPaid(ctx, ids)
}

// LoadForCheckout 加载未支付购物车并校验可售（status/is_show/mer_status=1）
func (s *Service) LoadForCheckout(ctx context.Context, uid uint, ids []uint64) ([]Cart, error) {
	if len(ids) == 0 {
		return nil, ErrNotFound
	}
	rows, err := s.store.ListByIDs(ctx, uid, ids)
	if err != nil {
		return nil, err
	}
	if len(rows) != len(ids) {
		return nil, ErrNotFound
	}
	out := make([]Cart, 0, len(rows))
	for i := range rows {
		item, err := s.enrichOne(ctx, &rows[i])
		if err != nil || item.IsFail == 1 {
			return nil, ErrProductOff
		}
		if item.Stock < item.CartNum {
			return nil, ErrStockNotEnough
		}
		out = append(out, *item)
	}
	return out, nil
}

func (s *Service) enrichOne(ctx context.Context, c *Cart) (*Cart, error) {
	pv, err := s.store.ResolveProduct(ctx, c.ProductID, c.ProductAttrUnique)
	if err != nil {
		return nil, err
	}
	c.StoreName = pv.StoreName
	c.MerName = pv.MerName
	c.Image = pv.Image
	c.Price = pv.Price
	c.OtPrice = pv.OtPrice
	c.Cost = pv.Cost
	c.Stock = pv.Stock
	c.MerID = pv.MerID
	c.GoodsType = pv.GoodsType
	c.ProductType = int8(pv.ProductType)
	c.SvipPriceType = pv.SvipPriceType
	c.SvipPrice = pv.SvipPrice
	c.MerSvipStatus = pv.MerSvipStatus
	if c.ProductAttrUnique == "" {
		c.ProductAttrUnique = pv.Unique
	}
	if pv.Status != 1 || pv.IsShow != 1 || pv.MerStatus != 1 {
		c.IsFail = 1
	}
	return c, nil
}

func (s *Service) ListAddresses(ctx context.Context, uid uint) ([]Address, error) {
	return s.store.ListAddresses(ctx, uid)
}

func (s *Service) GetAddress(ctx context.Context, uid, id uint) (*Address, error) {
	a, err := s.store.GetAddress(ctx, id, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAddrNotFound
		}
		return nil, err
	}
	return a, nil
}

func (s *Service) CreateAddress(ctx context.Context, uid uint, in AddressInput) (*Address, error) {
	return s.SaveAddress(ctx, uid, 0, in)
}

func (s *Service) UpdateAddress(ctx context.Context, uid, id uint, in AddressInput) (*Address, error) {
	return s.SaveAddress(ctx, uid, id, in)
}

func (s *Service) SaveAddress(ctx context.Context, uid uint, id uint, in AddressInput) (*Address, error) {
	name := strings.TrimSpace(in.RealName)
	phone := strings.TrimSpace(in.Phone)
	detail := strings.TrimSpace(in.Detail)
	if name == "" || phone == "" || detail == "" {
		return nil, ErrAddrInvalid
	}
	isDefault := int8(0)
	if in.IsDefault != nil {
		isDefault = *in.IsDefault
	}
	if isDefault == 1 {
		_ = s.store.ClearDefaultAddress(ctx, uid)
	}
	if id == 0 {
		a := &Address{
			UID: uid, RealName: name, Phone: phone,
			Province: in.Province, City: in.City, District: in.District,
			Detail: detail, PostCode: in.PostCode, IsDefault: isDefault,
		}
		if err := s.store.CreateAddress(ctx, a); err != nil {
			return nil, err
		}
		return a, nil
	}
	a, err := s.GetAddress(ctx, uid, id)
	if err != nil {
		return nil, err
	}
	a.RealName, a.Phone, a.Detail = name, phone, detail
	a.Province, a.City, a.District = in.Province, in.City, in.District
	a.PostCode, a.IsDefault = in.PostCode, isDefault
	if err := s.store.UpdateAddress(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) DeleteAddress(ctx context.Context, uid, id uint) error {
	if _, err := s.GetAddress(ctx, uid, id); err != nil {
		return err
	}
	return s.store.SoftDeleteAddress(ctx, id, uid)
}
