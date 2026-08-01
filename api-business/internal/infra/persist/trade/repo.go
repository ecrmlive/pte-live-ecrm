package tradepersist

import (
	"context"
	"errors"
	"time"

	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/cart"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/catalog"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/trade"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) WithTx(fn func(tx trade.Store) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(&Repo{db: tx})
	})
}

func (r *Repo) CreateGroupOrder(ctx context.Context, g *trade.GroupOrder) error {
	return r.db.WithContext(ctx).Create(g).Error
}

func (r *Repo) CreateStoreOrder(ctx context.Context, o *trade.StoreOrder) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *Repo) CreateOrderProduct(ctx context.Context, p *trade.OrderProduct) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *Repo) GetGroupOrder(ctx context.Context, id uint) (*trade.GroupOrder, error) {
	var g trade.GroupOrder
	err := r.db.WithContext(ctx).Where("group_order_id = ? AND is_del = 0", id).First(&g).Error
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *Repo) ListGroupOrders(ctx context.Context, uid uint, page, limit int) ([]trade.GroupOrder, int64, error) {
	page, limit = normalize(page, limit)
	q := r.db.WithContext(ctx).Model(&trade.GroupOrder{}).Where("uid = ? AND is_del = 0", uid)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []trade.GroupOrder
	err := q.Order("group_order_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) ListExpiredUnpaidGroups(ctx context.Context, before time.Time, limit int) ([]trade.GroupOrder, error) {
	if limit <= 0 {
		limit = 50
	}
	var rows []trade.GroupOrder
	err := r.db.WithContext(ctx).
		Where("paid = 0 AND is_del = 0 AND create_time < ?", before).
		Order("group_order_id ASC").Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *Repo) ListStoreOrdersByGroup(ctx context.Context, groupID uint) ([]trade.StoreOrder, error) {
	var rows []trade.StoreOrder
	err := r.db.WithContext(ctx).Where("group_order_id = ? AND is_del = 0", groupID).
		Order("order_id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) ListOrderProductsByOrder(ctx context.Context, orderID uint) ([]trade.OrderProduct, error) {
	var rows []trade.OrderProduct
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&rows).Error
	return rows, err
}

func (r *Repo) ListOrderProductsByOrders(ctx context.Context, orderIDs []uint) ([]trade.OrderProduct, error) {
	if len(orderIDs) == 0 {
		return nil, nil
	}
	var rows []trade.OrderProduct
	err := r.db.WithContext(ctx).Where("order_id IN ?", orderIDs).Find(&rows).Error
	return rows, err
}

func (r *Repo) GetStoreOrder(ctx context.Context, id uint) (*trade.StoreOrder, error) {
	var o trade.StoreOrder
	err := r.db.WithContext(ctx).Where("order_id = ?", id).First(&o).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *Repo) ListStoreOrders(ctx context.Context, merID *uint, page, limit int) ([]trade.StoreOrder, int64, error) {
	return r.ListStoreOrdersFiltered(ctx, merID, nil, nil, page, limit)
}

func (r *Repo) ListStoreOrdersFiltered(ctx context.Context, merID *uint, paid, status *int8, page, limit int) ([]trade.StoreOrder, int64, error) {
	page, limit = normalize(page, limit)
	q := r.db.WithContext(ctx).Model(&trade.StoreOrder{}).Where("is_del = 0")
	if merID != nil {
		q = q.Where("mer_id = ?", *merID)
	}
	if paid != nil {
		q = q.Where("paid = ?", *paid)
	}
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []trade.StoreOrder
	err := q.Order("order_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) ListPlatformOrdersByRegions(ctx context.Context, regionIDs []uint, paid *int8, page, limit int) ([]trade.StoreOrder, int64, error) {
	page, limit = normalize(page, limit)
	if len(regionIDs) == 0 {
		return []trade.StoreOrder{}, 0, nil
	}
	q := r.db.WithContext(ctx).Table("qixi_m_app_store_order AS o").
		Joins("JOIN qixi_m_admin_merchant AS m ON m.mer_id = o.mer_id").
		Where("o.is_del = 0 AND m.is_del = 0 AND m.region_id IN ?", regionIDs)
	if paid != nil {
		q = q.Where("o.paid = ?", *paid)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []trade.StoreOrder
	err := q.Select("o.*").Order("o.order_id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error
	return rows, total, err
}

func (r *Repo) GetPlatformOrderByRegions(ctx context.Context, orderID uint, regionIDs []uint) (*trade.StoreOrder, error) {
	if len(regionIDs) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var row trade.StoreOrder
	err := r.db.WithContext(ctx).Table("qixi_m_app_store_order AS o").
		Select("o.*").Joins("JOIN qixi_m_admin_merchant AS m ON m.mer_id = o.mer_id").
		Where("o.order_id = ? AND o.is_del = 0 AND m.is_del = 0 AND m.region_id IN ?", orderID, regionIDs).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) MerchantName(ctx context.Context, merID uint) (string, error) {
	var name string
	err := r.db.WithContext(ctx).Table("qixi_m_admin_merchant").Select("mer_name").Where("mer_id = ?", merID).Scan(&name).Error
	return name, err
}

func (r *Repo) MarkGroupPaid(ctx context.Context, id uint, payType int8, payTime time.Time) error {
	return r.db.WithContext(ctx).Model(&trade.GroupOrder{}).
		Where("group_order_id = ? AND paid = 0", id).
		Updates(map[string]interface{}{"paid": 1, "pay_type": payType, "pay_time": payTime}).Error
}

func (r *Repo) MarkChildrenPaid(ctx context.Context, groupID uint, payType int8, payTime time.Time) error {
	return r.db.WithContext(ctx).Model(&trade.StoreOrder{}).
		Where("group_order_id = ? AND paid = 0", groupID).
		Updates(map[string]interface{}{
			"paid": 1, "pay_type": payType, "pay_time": payTime, "status": trade.OrderStatusAwaitShip,
		}).Error
}

func (r *Repo) MarkCartsPaid(ctx context.Context, uid uint, cartIDs []uint64) error {
	if len(cartIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&cart.Cart{}).
		Where("uid = ? AND cart_id IN ?", uid, cartIDs).
		Update("is_pay", 1).Error
}

func (r *Repo) DeductProductStock(ctx context.Context, productID uint, num uint) error {
	res := r.db.WithContext(ctx).Model(&catalog.Product{}).
		Where("product_id = ? AND stock >= ?", productID, num).
		Updates(map[string]interface{}{
			"stock": gorm.Expr("stock - ?", num),
			"sales": gorm.Expr("sales + ?", num),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return trade.ErrStockNotEnough
	}
	return nil
}

func (r *Repo) DeductSKUStock(ctx context.Context, productID uint, unique string, num uint) error {
	q := r.db.WithContext(ctx).Model(&catalog.AttrValue{}).
		Where("product_id = ? AND stock >= ?", productID, num)
	if unique != "" {
		q = q.Where("`unique` = ?", unique)
	}
	return q.Updates(map[string]interface{}{
		"stock": gorm.Expr("stock - ?", num),
		"sales": gorm.Expr("sales + ?", num),
	}).Error
}

func (r *Repo) GetUserBalance(ctx context.Context, uid uint) (float64, error) {
	var bal float64
	err := r.db.WithContext(ctx).Table("qixi_m_app_user").Select("now_money").
		Where("uid = ?", uid).Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&bal).Error
	return bal, err
}

func (r *Repo) DeductUserBalance(ctx context.Context, uid uint, amount float64) error {
	res := r.db.WithContext(ctx).Table("qixi_m_app_user").
		Where("uid = ? AND now_money >= ?", uid, amount).
		Update("now_money", gorm.Expr("now_money - ?", amount))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return trade.ErrBalanceNotEnough
	}
	return nil
}

func (r *Repo) GetUserIntegral(ctx context.Context, uid uint) (int, error) {
	var n int
	err := r.db.WithContext(ctx).Table("qixi_m_app_user").Select("integral").
		Where("uid = ?", uid).Scan(&n).Error
	return n, err
}

func (r *Repo) DeductUserIntegral(ctx context.Context, uid uint, amount int) error {
	if amount <= 0 {
		return nil
	}
	res := r.db.WithContext(ctx).Table("qixi_m_app_user").
		Where("uid = ? AND integral >= ?", uid, amount).
		Update("integral", gorm.Expr("integral - ?", amount))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return trade.ErrIntegralNotEnough
	}
	return nil
}

func (r *Repo) AddUserIntegral(ctx context.Context, uid uint, amount int) (int, error) {
	if amount <= 0 {
		return r.GetUserIntegral(ctx, uid)
	}
	if err := r.db.WithContext(ctx).Table("qixi_m_app_user").
		Where("uid = ?", uid).
		Update("integral", gorm.Expr("integral + ?", amount)).Error; err != nil {
		return 0, err
	}
	return r.GetUserIntegral(ctx, uid)
}

func (r *Repo) MerchantsIntegralEnabled(ctx context.Context, merIDs []uint) (bool, error) {
	if len(merIDs) == 0 {
		return true, nil
	}
	var n int64
	err := r.db.WithContext(ctx).Table("qixi_m_admin_merchant").
		Where("mer_id IN ? AND mer_integral_status = 0", merIDs).
		Count(&n).Error
	if err != nil {
		return false, err
	}
	return n == 0, nil
}

func (r *Repo) GetUserSVIP(ctx context.Context, uid uint) (int8, *time.Time, error) {
	var row struct {
		IsSvip      int8       `gorm:"column:is_svip"`
		SvipEndtime *time.Time `gorm:"column:svip_endtime"`
	}
	err := r.db.WithContext(ctx).Table("qixi_m_app_user").
		Select("is_svip, svip_endtime").Where("uid = ?", uid).Take(&row).Error
	if err != nil {
		return -1, nil, err
	}
	return row.IsSvip, row.SvipEndtime, nil
}

func (r *Repo) MerchantsSVIPCouponMerge(ctx context.Context, merIDs []uint) (map[uint]int8, error) {
	out := map[uint]int8{}
	if len(merIDs) == 0 {
		return out, nil
	}
	type row struct {
		MerID           uint `gorm:"column:mer_id"`
		SvipCouponMerge int8 `gorm:"column:svip_coupon_merge"`
	}
	var rows []row
	err := r.db.WithContext(ctx).Table("qixi_m_admin_merchant").
		Select("mer_id, svip_coupon_merge").Where("mer_id IN ?", merIDs).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, x := range rows {
		out[x.MerID] = x.SvipCouponMerge
	}
	return out, nil
}

func (r *Repo) HasBill(ctx context.Context, uid uint, category, typ, linkID string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("qixi_m_app_user_bill").
		Where("uid = ? AND category = ? AND type = ? AND link_id = ?", uid, category, typ, linkID).
		Count(&n).Error
	return n > 0, err
}

func (r *Repo) CreateBill(ctx context.Context, b *trade.UserBill) error {
	if b.CreateTime.IsZero() {
		b.CreateTime = time.Now()
	}
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *Repo) LoadPointsProduct(ctx context.Context, productID uint, unique string) (*trade.PointsProductView, error) {
	var p catalog.Product
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND is_del = 0 AND product_type = ?", productID, catalog.ProductTypePoints).
		First(&p).Error
	if err != nil {
		return nil, err
	}
	if p.Status != 1 || p.IsShow != 1 || p.MerStatus != 1 {
		return nil, trade.ErrNotPointsProduct
	}
	need := p.Integral
	if need <= 0 {
		need = int(p.OtPrice) // 兼容仅写 ot_price 的种子
	}
	pv := &trade.PointsProductView{
		ProductID:   p.ProductID,
		MerID:       p.MerID,
		StoreName:   p.StoreName,
		Image:       p.Image,
		Price:       p.Price,
		Stock:       p.Stock,
		Integral:    need,
		ProductType: p.ProductType,
	}
	_ = r.db.WithContext(ctx).Table("qixi_m_admin_merchant").Select("mer_name").
		Where("mer_id = ?", p.MerID).Scan(&pv.MerName).Error
	var sku catalog.AttrValue
	q := r.db.WithContext(ctx).Where("product_id = ?", productID)
	if unique != "" {
		q = q.Where("`unique` = ?", unique)
	}
	if err := q.Order("value_id ASC").First(&sku).Error; err == nil {
		pv.Unique = sku.Unique
		pv.Price = sku.Price
		pv.Stock = sku.Stock
		pv.Cost = sku.Cost
		if sku.Image != "" {
			pv.Image = sku.Image
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return pv, nil
}

func (r *Repo) MarkCouponUsersUsed(ctx context.Context, uid uint, ids []uint, at time.Time) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	res := r.db.WithContext(ctx).Table("qixi_m_app_store_coupon_user").
		Where("uid = ? AND coupon_user_id IN ? AND status = 0", uid, ids).
		Updates(map[string]interface{}{
			"status":   1,
			"use_time": at,
		})
	return res.RowsAffected, res.Error
}

func (r *Repo) UpdateOrderStatus(ctx context.Context, orderID uint, status int8) error {
	return r.db.WithContext(ctx).Model(&trade.StoreOrder{}).
		Where("order_id = ?", orderID).Update("status", status).Error
}

func (r *Repo) MarkOrderVerified(ctx context.Context, orderID, merID, serviceID uint, at time.Time) error {
	updates := map[string]interface{}{
		"status":      trade.OrderStatusDone,
		"verify_time": at,
	}
	if serviceID > 0 {
		updates["verify_service_id"] = serviceID
	}
	res := r.db.WithContext(ctx).Model(&trade.StoreOrder{}).
		Where("order_id = ? AND mer_id = ? AND paid = 1 AND is_del = 0", orderID, merID).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return trade.ErrBadStatus
	}
	return nil
}

func (r *Repo) GetStoreOrderByVerifyCode(ctx context.Context, merID uint, code string) (*trade.StoreOrder, error) {
	var o trade.StoreOrder
	err := r.db.WithContext(ctx).
		Where("mer_id = ? AND verify_code = ? AND is_del = 0", merID, code).
		First(&o).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *Repo) ListStoreOrdersInStatuses(ctx context.Context, merID uint, paid *int8, statuses []int8, page, limit int) ([]trade.StoreOrder, int64, error) {
	page, limit = normalize(page, limit)
	q := r.db.WithContext(ctx).Model(&trade.StoreOrder{}).Where("is_del = 0 AND mer_id = ?", merID)
	if paid != nil {
		q = q.Where("paid = ?", *paid)
	}
	if len(statuses) > 0 {
		q = q.Where("status IN ?", statuses)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []trade.StoreOrder
	err := q.Order("order_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) SoftDeleteGroup(ctx context.Context, id, uid uint) error {
	return r.db.WithContext(ctx).Model(&trade.GroupOrder{}).
		Where("group_order_id = ? AND uid = ? AND paid = 0", id, uid).
		Update("is_del", 1).Error
}

func (r *Repo) SoftDeleteOrdersByGroup(ctx context.Context, groupID uint) error {
	return r.db.WithContext(ctx).Model(&trade.StoreOrder{}).
		Where("group_order_id = ?", groupID).Update("is_del", 1).Error
}

func (r *Repo) DeliverOrder(ctx context.Context, orderID, merID uint, name, deliveryID, deliveryType string) error {
	res := r.db.WithContext(ctx).Model(&trade.StoreOrder{}).
		Where("order_id = ? AND mer_id = ? AND paid = 1 AND status = ?", orderID, merID, trade.OrderStatusAwaitShip).
		Updates(map[string]interface{}{
			"delivery_name": name, "delivery_id": deliveryID, "delivery_type": deliveryType,
			"status": trade.OrderStatusShipped,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return trade.ErrBadStatus
	}
	return nil
}

func (r *Repo) CreatePresellOrder(ctx context.Context, o *trade.PresellOrder) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *Repo) GetPresellOrder(ctx context.Context, id uint) (*trade.PresellOrder, error) {
	var row trade.PresellOrder
	err := r.db.WithContext(ctx).Where("presell_order_id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) GetPresellOrderByOrderID(ctx context.Context, orderID uint) (*trade.PresellOrder, error) {
	var row trade.PresellOrder
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) ListPresellOrdersByUID(ctx context.Context, uid uint, unpaidOnly bool, page, limit int) ([]trade.PresellOrder, int64, error) {
	page, limit = normalize(page, limit)
	q := r.db.WithContext(ctx).Model(&trade.PresellOrder{}).Where("uid = ?", uid)
	if unpaidOnly {
		q = q.Where("paid = 0 AND status = 1")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []trade.PresellOrder
	err := q.Order("presell_order_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) MarkPresellOrderPaid(ctx context.Context, id uint, payType int8, at time.Time) error {
	return r.db.WithContext(ctx).Model(&trade.PresellOrder{}).
		Where("presell_order_id = ? AND paid = 0", id).
		Updates(map[string]interface{}{"paid": 1, "pay_type": payType, "pay_time": at}).Error
}

func (r *Repo) InvalidatePresellOrder(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&trade.PresellOrder{}).
		Where("presell_order_id = ? AND paid = 0", id).
		Update("status", 0).Error
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

var _ trade.Store = (*Repo)(nil)
