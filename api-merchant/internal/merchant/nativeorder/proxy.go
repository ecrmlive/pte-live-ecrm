package nativeorder

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	errProxyBadRequest = errors.New("proxy bad request")
	errProxyConflict   = errors.New("proxy conflict")
	errProxyNotFound   = errors.New("proxy not found")
)

type proxyRequest struct {
	UserID         uint64 `json:"user_id"`
	ProductID      uint64 `json:"product_id"`
	Quantity       int    `json:"quantity"`
	Remark         string `json:"remark"`
	IdempotencyKey string `json:"idempotency_key"`
}

type proxyProduct struct {
	ProductID    uint64  `gorm:"column:product_id"`
	MerchantID   uint64  `gorm:"column:merchant_id"`
	StoreID      uint64  `gorm:"column:store_id"`
	MerchantName string  `gorm:"column:merchant_name"`
	StoreName    string  `gorm:"column:store_name"`
	Title        string  `gorm:"column:title"`
	CoverURL     string  `gorm:"column:cover_url"`
	SaleStatus   int     `gorm:"column:sale_status"`
	ProductType  int     `gorm:"column:product_type"`
}

type proxySKU struct {
	MerchantSKUID uint64          `gorm:"column:merchant_sku_id"`
	ProductID     uint64          `gorm:"column:product_id"`
	SKUKey        string          `gorm:"column:sku_key"`
	SpecSnapshot  json.RawMessage `gorm:"column:spec_snapshot"`
	Price         float64         `gorm:"column:price"`
	Stock         int             `gorm:"column:stock"`
	SaleStatus    int             `gorm:"column:sale_status"`
}

type proxyAddress struct {
	Recipient string `json:"recipient" gorm:"column:recipient"`
	Mobile    string `json:"mobile" gorm:"column:mobile"`
	Province  string `json:"province" gorm:"column:province"`
	City      string `json:"city" gorm:"column:city"`
	District  string `json:"district" gorm:"column:district"`
	Detail    string `json:"detail" gorm:"column:detail"`
	PostCode  int    `json:"post_code" gorm:"column:post_code"`
}

func (h *Handler) proxy(c *gin.Context) {
	var req proxyRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.UserID == 0 || req.ProductID == 0 || req.Quantity < 1 {
		response.Fail(c, http.StatusBadRequest, "代客下单参数不合法")
		return
	}
	req.Remark = strings.TrimSpace(req.Remark)
	if utf8.RuneCountInString(req.Remark) > 200 {
		response.Fail(c, http.StatusBadRequest, "订单备注不能超过 200 个字符")
		return
	}
	req.IdempotencyKey = strings.TrimSpace(req.IdempotencyKey)
	if req.IdempotencyKey == "" {
		req.IdempotencyKey = fmt.Sprintf("proxy:%d:%d:%s", middleware.StoreID(c), middleware.AdminID(c), randomToken(16))
	}
	if len(req.IdempotencyKey) < 12 || len(req.IdempotencyKey) > 128 {
		response.Fail(c, http.StatusBadRequest, "幂等键长度须为 12–128")
		return
	}

	storeID := uint64(middleware.StoreID(c))
	result, err := h.createProxyOrder(c, storeID, req)
	if errors.Is(err, errProxyBadRequest) {
		response.Fail(c, http.StatusBadRequest, result.Message)
		return
	}
	if errors.Is(err, errProxyNotFound) {
		response.Fail(c, http.StatusNotFound, result.Message)
		return
	}
	if errors.Is(err, errProxyConflict) {
		response.Fail(c, http.StatusConflict, result.Message)
		return
	}
	if err != nil {
		fail(c, "代客下单失败")
		return
	}
	response.OK(c, gin.H{
		"group_order_id": result.GroupOrderID,
		"order_id":       result.OrderID,
		"order_sn":       result.OrderSN,
		"pay_amount":     result.PayAmount,
		"paid":           0,
		"replayed":       result.Replayed,
	})
}

type proxyResult struct {
	GroupOrderID uint64
	OrderID      uint64
	OrderSN      string
	PayAmount    float64
	Replayed     bool
	Message      string
}

func (h *Handler) createProxyOrder(c *gin.Context, storeID uint64, req proxyRequest) (proxyResult, error) {
	ctx := c.Request.Context()
	var existing struct {
		ID        uint64  `gorm:"column:id"`
		OrderNo   string  `gorm:"column:order_no"`
		PayAmount float64 `gorm:"column:pay_amount"`
	}
	err := h.db.WithContext(ctx).Table("qixi_crm_b_group_order").
		Select("id,order_no,pay_amount").
		Where("user_id = ? AND idempotency_key = ?", req.UserID, req.IdempotencyKey).
		Take(&existing).Error
	if err == nil {
		var child struct {
			ID      uint64 `gorm:"column:id"`
			OrderNo string `gorm:"column:order_no"`
		}
		_ = h.db.WithContext(ctx).Table("qixi_crm_b_order").
			Select("id,order_no").
			Where("group_order_id = ? AND store_id = ?", existing.ID, storeID).
			Order("id ASC").Limit(1).Take(&child).Error
		return proxyResult{
			GroupOrderID: existing.ID,
			OrderID:      child.ID,
			OrderSN:      child.OrderNo,
			PayAmount:    existing.PayAmount,
			Replayed:     true,
		}, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return proxyResult{}, err
	}

	var buyer struct {
		ID     uint64 `gorm:"column:id"`
		Status int    `gorm:"column:status"`
	}
	err = h.db.WithContext(ctx).Table("qixi_crm_b_user").Select("id,status").Where("id = ?", req.UserID).Take(&buyer).Error
	if errors.Is(err, gorm.ErrRecordNotFound) || buyer.ID == 0 {
		return proxyResult{Message: "客户不存在或已停用"}, errProxyNotFound
	}
	if err != nil {
		return proxyResult{}, err
	}
	if buyer.Status != 1 {
		return proxyResult{Message: "客户不存在或已停用"}, errProxyNotFound
	}
	var prior int64
	if err := h.db.WithContext(ctx).Table("qixi_crm_b_order").
		Where("store_id = ? AND user_id = ? AND status NOT IN ('pending_pay','cancelled')", storeID, req.UserID).
		Count(&prior).Error; err != nil {
		return proxyResult{}, err
	}
	if prior == 0 {
		return proxyResult{Message: "仅可为在本店有过有效订单的客户代客下单"}, errProxyBadRequest
	}

	var address proxyAddress
	err = h.db.WithContext(ctx).Table("qixi_crm_b_user_address").
		Select("recipient,mobile,province,city,district,detail,post_code").
		Where("user_id = ?", req.UserID).
		Order("is_default DESC, id DESC").
		Limit(1).Take(&address).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return proxyResult{Message: "客户暂无收货地址，请先引导客户维护地址"}, errProxyBadRequest
	}
	if err != nil {
		return proxyResult{}, err
	}
	addressJSON, err := json.Marshal(address)
	if err != nil {
		return proxyResult{}, err
	}

	var product proxyProduct
	err = h.db.WithContext(ctx).Table("qixi_crm_b_product_view").
		Where("product_id = ? AND store_id = ? AND sale_status = 1 AND product_type = 0", req.ProductID, storeID).
		Take(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return proxyResult{Message: "商品不存在、未上架或不属于本店"}, errProxyNotFound
	}
	if err != nil {
		return proxyResult{}, err
	}

	var skus []proxySKU
	if err := h.db.WithContext(ctx).Table("qixi_crm_b_product_sku_view").
		Where("product_id = ? AND sale_status = 1", req.ProductID).
		Order("merchant_sku_id ASC").Find(&skus).Error; err != nil {
		return proxyResult{}, err
	}
	if len(skus) == 0 {
		return proxyResult{Message: "商品无可售规格"}, errProxyBadRequest
	}
	if len(skus) > 1 {
		return proxyResult{Message: "多规格商品暂不支持代客下单，请选择单规格商品"}, errProxyBadRequest
	}
	sku := skus[0]
	if sku.MerchantSKUID == 0 {
		return proxyResult{Message: "商品规格映射缺失"}, errProxyBadRequest
	}
	if sku.Stock < req.Quantity {
		return proxyResult{Message: "库存不足"}, errProxyConflict
	}
	if len(sku.SpecSnapshot) == 0 {
		sku.SpecSnapshot = json.RawMessage(`{}`)
	}

	payAmount := roundMoney(sku.Price * float64(req.Quantity))
	expiresAt := time.Now().UTC().Add(30 * time.Minute)
	var out proxyResult
	err = h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		groupNo := orderNo("G")
		if err := tx.Exec(`INSERT INTO qixi_crm_b_group_order
			(order_no,user_id,total_amount,discount_amount,freight_amount,pay_amount,total_quantity,activity_type,points_amount,recipient_snapshot,pay_status,idempotency_key,remark)
			VALUES (?,?,?,0,0,?,?,0,0,?,'pending',?,?)`,
			groupNo, req.UserID, payAmount, payAmount, req.Quantity, string(addressJSON), req.IdempotencyKey, req.Remark,
		).Error; err != nil {
			if isDuplicateKey(err) {
				return errProxyConflict
			}
			return err
		}
		var groupID uint64
		if err := tx.Raw(`SELECT id FROM qixi_crm_b_group_order WHERE user_id = ? AND idempotency_key = ?`, req.UserID, req.IdempotencyKey).Scan(&groupID).Error; err != nil {
			return err
		}
		storeNo := orderNo("S")
		if err := tx.Exec(`INSERT INTO qixi_crm_b_order
			(group_order_id,order_no,merchant_id,merchant_name_snapshot,store_id,store_name_snapshot,user_id,total_amount,discount_amount,freight_amount,pay_amount,total_quantity,activity_type,points_amount,recipient_snapshot,remark,status)
			VALUES (?,?,?,?,?,?,?,?,0,0,?,?,0,0,?,?,'pending_pay')`,
			groupID, storeNo, product.MerchantID, product.MerchantName, storeID, product.StoreName, req.UserID,
			payAmount, payAmount, req.Quantity, string(addressJSON), req.Remark,
		).Error; err != nil {
			return err
		}
		var orderID uint64
		if err := tx.Raw(`SELECT id FROM qixi_crm_b_order WHERE order_no = ?`, storeNo).Scan(&orderID).Error; err != nil {
			return err
		}
		if err := tx.Exec(`INSERT INTO qixi_crm_b_order_item
			(order_id,product_id,merchant_sku_id,sku_key,title_snapshot,cover_url_snapshot,spec_snapshot,unit_price,quantity,refund_quantity)
			VALUES (?,?,?,?,?,?,?,?,?,0)`,
			orderID, product.ProductID, sku.MerchantSKUID, sku.SKUKey, product.Title, product.CoverURL, string(sku.SpecSnapshot), sku.Price, req.Quantity,
		).Error; err != nil {
			return err
		}
		stockKey := fmt.Sprintf("stock:reserve:%d:%d", orderID, sku.MerchantSKUID)
		if err := tx.Exec(`INSERT INTO qixi_crm_b_stock_command_outbox
			(action,order_id,store_id,merchant_sku_id,quantity,expires_at,idempotency_key,status)
			VALUES ('reserve',?,?,?,?,?,?,'pending')`,
			orderID, storeID, sku.MerchantSKUID, req.Quantity, expiresAt, stockKey,
		).Error; err != nil {
			return err
		}
		out = proxyResult{GroupOrderID: groupID, OrderID: orderID, OrderSN: storeNo, PayAmount: payAmount}
		return nil
	})
	if errors.Is(err, errProxyConflict) {
		var raced struct {
			ID        uint64  `gorm:"column:id"`
			OrderNo   string  `gorm:"column:order_no"`
			PayAmount float64 `gorm:"column:pay_amount"`
		}
		if loadErr := h.db.WithContext(ctx).Table("qixi_crm_b_group_order").
			Select("id,order_no,pay_amount").
			Where("user_id = ? AND idempotency_key = ?", req.UserID, req.IdempotencyKey).
			Take(&raced).Error; loadErr == nil {
			var child struct {
				ID      uint64 `gorm:"column:id"`
				OrderNo string `gorm:"column:order_no"`
			}
			_ = h.db.WithContext(ctx).Table("qixi_crm_b_order").
				Select("id,order_no").
				Where("group_order_id = ? AND store_id = ?", raced.ID, storeID).
				Order("id ASC").Limit(1).Take(&child).Error
			return proxyResult{
				GroupOrderID: raced.ID,
				OrderID:      child.ID,
				OrderSN:      child.OrderNo,
				PayAmount:    raced.PayAmount,
				Replayed:     true,
			}, nil
		}
		return proxyResult{Message: "代客下单幂等键冲突，请重试"}, errProxyConflict
	}
	return out, err
}

func orderNo(prefix string) string {
	return fmt.Sprintf("%s%s%s", prefix, time.Now().UTC().Format("20060102150405"), randomToken(4))
}

func randomToken(n int) string {
	buf := make([]byte, n)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}

func roundMoney(v float64) float64 {
	return float64(int64(v*100+0.5)) / 100
}

func isDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "uk_user_idempotency")
}
