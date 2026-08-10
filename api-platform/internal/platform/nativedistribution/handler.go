// Package nativedistribution exposes the platform/operations read model for
// business-owned distribution data, plus limited audited mutations for promoter
// list ops (clear parent / edit level). Commission balances are never adjusted here.
package nativedistribution

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	businessDB *gorm.DB
	adminDB    *gorm.DB
}

func NewHandler(businessDB, adminDB *gorm.DB) *Handler {
	return &Handler{businessDB: businessDB, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	read := middleware.RequireAdminRoles("platform", "operations")
	readSpread := middleware.RequireAdminMenu(h.adminDB, "marketing.spread.read")
	write := middleware.RequireAdminRoles("platform")
	r.GET("/distribution/promoters", read, readSpread, h.ListPromoters)
	r.GET("/distribution/promoters/:id/children", read, readSpread, h.ListChildren)
	r.GET("/distribution/promoters/:id/orders", read, readSpread, h.ListSpreadOrders)
	r.POST("/distribution/promoters/:id/clear-parent", write, readSpread, h.ClearParent)
	r.POST("/distribution/promoters/:id/level", write, readSpread, h.UpdatePromoterLevel)
	r.GET("/distribution/levels", read, readSpread, h.ListLevels)
	r.POST("/distribution/levels", write, readSpread, h.CreateLevel)
	r.GET("/distribution/levels/:id", read, readSpread, h.GetLevel)
	r.PUT("/distribution/levels/:id", write, readSpread, h.UpdateLevel)
	r.DELETE("/distribution/levels/:id", write, readSpread, h.DeleteLevel)
	r.GET("/distribution/withdraw-banks", read, readSpread, h.ListWithdrawBanks)
	r.POST("/distribution/withdraw-banks", write, readSpread, h.CreateWithdrawBank)
	r.PUT("/distribution/withdraw-banks/:id", write, readSpread, h.UpdateWithdrawBank)
	r.PUT("/distribution/withdraw-banks/:id/status", write, readSpread, h.SetWithdrawBankStatus)
	r.DELETE("/distribution/withdraw-banks/:id", write, readSpread, h.DeleteWithdrawBank)
	r.GET("/distribution/privileges", read, readSpread, h.ListPrivileges)
	r.POST("/distribution/privileges", write, readSpread, h.CreatePrivilege)
	r.PUT("/distribution/privileges/:id", write, readSpread, h.UpdatePrivilege)
	r.PUT("/distribution/privileges/:id/status", write, readSpread, h.SetPrivilegeStatus)
	r.DELETE("/distribution/privileges/:id", write, readSpread, h.DeletePrivilege)
	r.GET("/distribution/posters", read, readSpread, h.ListPosters)
	r.POST("/distribution/posters", write, readSpread, h.CreatePoster)
	r.PUT("/distribution/posters/:id", write, readSpread, h.UpdatePoster)
	r.PUT("/distribution/posters/:id/status", write, readSpread, h.SetPosterStatus)
	r.DELETE("/distribution/posters/:id", write, readSpread, h.DeletePoster)
	r.GET("/distribution/commissions", read, readSpread, h.ListCommissions)
	r.GET("/distribution/summary", read, readSpread, h.Summary)
}

type promoterRow struct {
	UserID              uint64    `gorm:"column:user_id" json:"user_id"`
	Nickname            string    `gorm:"column:nickname" json:"nickname"`
	Mobile              string    `gorm:"column:mobile" json:"mobile"`
	AvatarURL           string    `gorm:"column:avatar_url" json:"avatar_url"`
	Status              int8      `gorm:"column:status" json:"status"`
	LevelID             uint64    `gorm:"column:level_id" json:"level_id"`
	LevelName           string    `gorm:"column:level_name" json:"level_name"`
	UpdatedAt           time.Time `gorm:"column:updated_at" json:"updated_at"`
	DirectUserCount     int64     `gorm:"column:direct_user_count" json:"direct_user_count"`
	SpreadUserCount     int64     `gorm:"column:spread_user_count" json:"spread_user_count"`
	SpreadOrderCount    int64     `gorm:"column:spread_order_count" json:"spread_order_count"`
	SpreadOrderAmount   float64   `gorm:"column:spread_order_amount" json:"spread_order_amount"`
	PendingCommission   float64   `gorm:"column:pending_commission" json:"pending_commission"`
	AvailableCommission float64   `gorm:"column:available_commission" json:"available_commission"`
	SettledCommission   float64   `gorm:"column:settled_commission" json:"settled_commission"`
	CommissionAmount    float64   `gorm:"column:commission_amount" json:"commission_amount"`
	WithdrawnAmount     float64   `gorm:"column:withdrawn_amount" json:"withdrawn_amount"`
	WithdrawCount       int64     `gorm:"column:withdraw_count" json:"withdraw_count"`
	UnwithdrawnAmount   float64   `gorm:"column:unwithdrawn_amount" json:"unwithdrawn_amount"`
	ParentUserID        uint64    `gorm:"column:parent_user_id" json:"parent_user_id"`
	ParentNickname      string    `gorm:"column:parent_nickname" json:"parent_nickname"`
}

type commissionRow struct {
	ID          uint64     `gorm:"column:id" json:"commission_id"`
	UserID      uint64     `gorm:"column:user_id" json:"user_id"`
	OrderID     uint64     `gorm:"column:order_id" json:"order_id"`
	Amount      float64    `gorm:"column:amount" json:"amount"`
	Status      string     `gorm:"column:status" json:"status"`
	AvailableAt *time.Time `gorm:"column:available_at" json:"available_at"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
}

type childRow struct {
	UserID          uint64    `gorm:"column:user_id" json:"user_id"`
	Nickname        string    `gorm:"column:nickname" json:"nickname"`
	Mobile          string    `gorm:"column:mobile" json:"mobile"`
	AvatarURL       string    `gorm:"column:avatar_url" json:"avatar_url"`
	BoundAt         time.Time `gorm:"column:bound_at" json:"bound_at"`
	IsPromoter      int8      `gorm:"column:is_promoter" json:"is_promoter"`
	SpreadUserCount int64     `gorm:"column:spread_user_count" json:"spread_user_count"`
	PayCount        int64     `gorm:"column:pay_count" json:"pay_count"`
	PayAmount       float64   `gorm:"column:pay_amount" json:"pay_amount"`
	Level           int8      `gorm:"column:level" json:"level"`
}

type spreadOrderRow struct {
	OrderID     uint64     `gorm:"column:order_id" json:"order_id"`
	OrderNo     string     `gorm:"column:order_no" json:"order_no"`
	PayAmount   float64    `gorm:"column:pay_amount" json:"pay_amount"`
	Commission  float64    `gorm:"column:commission" json:"commission"`
	Status      string     `gorm:"column:status" json:"status"`
	BuyerID     uint64     `gorm:"column:buyer_id" json:"buyer_id"`
	BuyerName   string     `gorm:"column:buyer_name" json:"buyer_name"`
	PaidAt      *time.Time `gorm:"column:paid_at" json:"paid_at"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
}

type levelTaskItem struct {
	Name string  `json:"name"`
	Num  float64 `json:"num"`
	Info string  `json:"info"`
}

type levelTaskRule struct {
	SpreadUser   levelTaskItem `json:"spread_user"`
	PayMoney     levelTaskItem `json:"pay_money"`
	PayNum       levelTaskItem `json:"pay_num"`
	SpreadMoney  levelTaskItem `json:"spread_money"`
	SpreadPayNum levelTaskItem `json:"spread_pay_num"`
}

type levelRow struct {
	ID            uint64        `gorm:"column:id" json:"id"`
	Name          string        `gorm:"column:name" json:"name"`
	Rank          int           `gorm:"column:rank" json:"rank"`
	IconURL       string        `gorm:"column:icon_url" json:"icon_url"`
	TaskRuleRaw   string        `gorm:"column:task_rule" json:"-"`
	TaskRule      levelTaskRule `gorm:"-" json:"task_rule"`
	ExtensionOne  float64       `gorm:"column:extension_one" json:"extension_one"`
	ExtensionTwo  float64       `gorm:"column:extension_two" json:"extension_two"`
	Status        int8          `gorm:"column:status" json:"status"`
	PromoterCount int64         `gorm:"column:promoter_count" json:"promoter_count"`
}

type levelSaveInput struct {
	Name         string        `json:"name"`
	Rank         int           `json:"rank"`
	IconURL      string        `json:"icon_url"`
	TaskRule     levelTaskRule `json:"task_rule"`
	ExtensionOne float64       `json:"extension_one"`
	ExtensionTwo float64       `json:"extension_two"`
	Status       *int8         `json:"status"`
}

type distributionSummary struct {
	PromoterCount       int64   `gorm:"column:promoter_count" json:"promoter_count"`
	ActivePromoterCount int64   `gorm:"column:active_promoter_count" json:"active_promoter_count"`
	SpreadUserCount     int64   `gorm:"column:spread_user_count" json:"spread_user_count"`
	SpreadOrderCount    int64   `gorm:"column:spread_order_count" json:"spread_order_count"`
	SpreadOrderAmount   float64 `gorm:"column:spread_order_amount" json:"spread_order_amount"`
	WithdrawnAmount     float64 `gorm:"column:withdrawn_amount" json:"withdrawn_amount"`
	UnwithdrawnAmount   float64 `gorm:"column:unwithdrawn_amount" json:"unwithdrawn_amount"`
	PendingCommission   float64 `gorm:"column:pending_commission" json:"pending_commission"`
	AvailableCommission float64 `gorm:"column:available_commission" json:"available_commission"`
	SettledCommission   float64 `gorm:"column:settled_commission" json:"settled_commission"`
}

func (h *Handler) ListPromoters(c *gin.Context) {
	page, limit := paging(c)
	q, ok := h.applyPromoterFilters(c, h.promoterQuery(c))
	if !ok {
		return
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		failure(c)
		return
	}
	orderSQL := promoterOrderSQL(c.Query("sort_field"), c.Query("sort_order"))
	rows := make([]promoterRow, 0)
	selectSQL := `
p.user_id,p.status,COALESCE(p.level_id,0) AS level_id,p.updated_at,
COALESCE(u.nickname,'') AS nickname,COALESCE(u.mobile,'') AS mobile,COALESCE(prof.avatar_url,'') AS avatar_url,
COALESCE(lv.name,'') AS level_name,
COALESCE(rel.direct_user_count,0) AS direct_user_count,
COALESCE(rel.direct_user_count,0) AS spread_user_count,
COALESCE(ord.spread_order_count,0) AS spread_order_count,
COALESCE(ord.spread_order_amount,0) AS spread_order_amount,
COALESCE(ledger.pending_commission,0) AS pending_commission,
COALESCE(ledger.available_commission,0) AS available_commission,
COALESCE(ledger.settled_commission,0) AS settled_commission,
COALESCE(ledger.pending_commission,0)+COALESCE(ledger.available_commission,0)+COALESCE(ledger.settled_commission,0) AS commission_amount,
COALESCE(wd.withdrawn_amount,0) AS withdrawn_amount,
COALESCE(wd.withdraw_count,0) AS withdraw_count,
COALESCE(ledger.available_commission,0) AS unwithdrawn_amount,
COALESCE(self_rel.parent_user_id,0) AS parent_user_id,
COALESCE(parent.nickname,'') AS parent_nickname`
	if err := q.Select(selectSQL).Order(orderSQL).Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		failure(c)
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) ListCommissions(c *gin.Context) {
	page, limit := paging(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_commission_ledger")
	if status, ok := commissionStatus(c.Query("status")); !ok {
		response.Fail(c, http.StatusBadRequest, "佣金状态错误")
		return
	} else if status != "" {
		q = q.Where("status = ?", status)
	}
	if userID, provided, ok := queryID(c, "user_id"); !ok {
		response.Fail(c, http.StatusBadRequest, "用户 ID 错误")
		return
	} else if provided {
		q = q.Where("user_id = ?", userID)
	}
	if dateFrom := strings.TrimSpace(c.Query("date_from")); dateFrom != "" {
		q = q.Where("created_at >= ?", dateFrom+" 00:00:00")
	}
	if dateTo := strings.TrimSpace(c.Query("date_to")); dateTo != "" {
		q = q.Where("created_at <= ?", dateTo+" 23:59:59")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		failure(c)
		return
	}
	rows := make([]commissionRow, 0)
	if err := q.Select("id,user_id,order_id,amount,status,available_at,created_at").Order("id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		failure(c)
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Summary(c *gin.Context) {
	q, ok := h.applyPromoterFilters(c, h.promoterQuery(c))
	if !ok {
		return
	}
	var out distributionSummary
	// Aggregate over filtered promoters; commission/withdraw subqueries scoped to those user_ids.
	err := q.Select(`
COUNT(*) AS promoter_count,
COALESCE(SUM(CASE WHEN p.status = 1 THEN 1 ELSE 0 END),0) AS active_promoter_count,
COALESCE(SUM(COALESCE(rel.direct_user_count,0)),0) AS spread_user_count,
COALESCE(SUM(COALESCE(ord.spread_order_count,0)),0) AS spread_order_count,
COALESCE(SUM(COALESCE(ord.spread_order_amount,0)),0) AS spread_order_amount,
COALESCE(SUM(COALESCE(wd.withdrawn_amount,0)),0) AS withdrawn_amount,
COALESCE(SUM(COALESCE(ledger.available_commission,0)),0) AS unwithdrawn_amount,
COALESCE(SUM(COALESCE(ledger.pending_commission,0)),0) AS pending_commission,
COALESCE(SUM(COALESCE(ledger.available_commission,0)),0) AS available_commission,
COALESCE(SUM(COALESCE(ledger.settled_commission,0)),0) AS settled_commission`).Scan(&out).Error
	if err != nil {
		failure(c)
		return
	}
	response.OK(c, out)
}

func (h *Handler) ListLevels(c *gin.Context) {
	pageRaw := strings.TrimSpace(c.Query("page"))
	manageMode := pageRaw != "" || strings.TrimSpace(c.Query("name")) != "" || strings.TrimSpace(c.Query("keyword")) != ""
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_level AS lv").
		Joins(`LEFT JOIN (
			SELECT level_id, COUNT(*) AS promoter_count
			FROM qixi_crm_b_distribution_promoter WHERE level_id > 0 GROUP BY level_id
		) AS pc ON pc.level_id = lv.id`)
	if !manageMode {
		q = q.Where("lv.status = 1")
	}
	name := strings.TrimSpace(c.Query("name"))
	if name == "" {
		name = strings.TrimSpace(c.Query("keyword"))
	}
	if name != "" {
		q = q.Where("lv.name LIKE ?", "%"+name+"%")
	}
	// `rank` is a MySQL reserved word (window functions) — must quote.
	selectSQL := `lv.id,lv.name,lv.` + "`rank`" + `,COALESCE(lv.icon_url,'') AS icon_url,
COALESCE(CAST(lv.task_rule AS CHAR),'') AS task_rule,
COALESCE(lv.extension_one,0) AS extension_one,COALESCE(lv.extension_two,0) AS extension_two,
lv.status,COALESCE(pc.promoter_count,0) AS promoter_count`
	if !manageMode {
		rows := make([]levelRow, 0)
		err := q.Select(selectSQL).Order("lv.`rank` ASC,lv.id ASC").Scan(&rows).Error
		if err != nil {
			if isMissingTable(err) || isMissingColumn(err) {
				response.OK(c, gin.H{"list": []levelRow{}})
				return
			}
			failure(c)
			return
		}
		for i := range rows {
			rows[i].TaskRule = parseLevelTaskRule(rows[i].TaskRuleRaw)
		}
		response.OK(c, gin.H{"list": rows})
		return
	}
	page, limit := paging(c)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		if isMissingTable(err) || isMissingColumn(err) {
			response.OK(c, gin.H{"list": []levelRow{}, "total": 0, "page": page, "limit": limit})
			return
		}
		failure(c)
		return
	}
	rows := make([]levelRow, 0)
	if err := q.Select(selectSQL).Order("lv.`rank` ASC,lv.id ASC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		failure(c)
		return
	}
	for i := range rows {
		rows[i].TaskRule = parseLevelTaskRule(rows[i].TaskRuleRaw)
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) GetLevel(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "分销等级 ID 错误")
		return
	}
	row, err := h.loadLevel(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "分销等级不存在")
			return
		}
		failure(c)
		return
	}
	response.OK(c, row)
}

func (h *Handler) CreateLevel(c *gin.Context) {
	var in levelSaveInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "分销等级参数错误")
		return
	}
	payload, msg := normalizeLevelSave(in)
	if msg != "" {
		response.Fail(c, http.StatusBadRequest, msg)
		return
	}
	if exists, err := h.levelRankExists(c, payload.Rank, 0); err != nil {
		failure(c)
		return
	} else if exists {
		response.Fail(c, http.StatusBadRequest, "等级已存在")
		return
	}
	ruleJSON, err := json.Marshal(payload.TaskRule)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "升级任务格式错误")
		return
	}
	status := int8(1)
	if payload.Status != nil {
		status = *payload.Status
	}
	res := h.businessDB.WithContext(c.Request.Context()).Exec(
		"INSERT INTO qixi_crm_b_distribution_level (`name`,`rank`,icon_url,task_rule,extension_one,extension_two,status) VALUES (?,?,?,?,?,?,?)",
		payload.Name, payload.Rank, payload.IconURL, string(ruleJSON), payload.ExtensionOne, payload.ExtensionTwo, status,
	)
	if res.Error != nil {
		if isDuplicateKey(res.Error) {
			response.Fail(c, http.StatusBadRequest, "等级已存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "创建分销等级失败")
		return
	}
	var id uint64
	_ = h.businessDB.WithContext(c.Request.Context()).Raw("SELECT LAST_INSERT_ID()").Scan(&id).Error
	row, err := h.loadLevel(c, id)
	if err != nil {
		response.OK(c, gin.H{"id": id})
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateLevel(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "分销等级 ID 错误")
		return
	}
	var in levelSaveInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "分销等级参数错误")
		return
	}
	existing, err := h.loadLevel(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "分销等级不存在")
			return
		}
		failure(c)
		return
	}
	payload, msg := normalizeLevelSave(in)
	if msg != "" {
		response.Fail(c, http.StatusBadRequest, msg)
		return
	}
	if exists, err := h.levelRankExists(c, payload.Rank, id); err != nil {
		failure(c)
		return
	} else if exists {
		response.Fail(c, http.StatusBadRequest, "等级已存在")
		return
	}
	ruleJSON, err := json.Marshal(payload.TaskRule)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "升级任务格式错误")
		return
	}
	status := existing.Status
	if payload.Status != nil {
		status = *payload.Status
	}
	// `rank` is reserved — keep backticks in raw SQL.
	res := h.businessDB.WithContext(c.Request.Context()).Exec(
		"UPDATE qixi_crm_b_distribution_level SET `name`=?,`rank`=?,icon_url=?,task_rule=?,extension_one=?,extension_two=?,status=?,updated_at=NOW() WHERE id=?",
		payload.Name, payload.Rank, payload.IconURL, string(ruleJSON), payload.ExtensionOne, payload.ExtensionTwo, status, id,
	)
	if res.Error != nil {
		if isDuplicateKey(res.Error) {
			response.Fail(c, http.StatusBadRequest, "等级已存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "更新分销等级失败")
		return
	}
	row, err := h.loadLevel(c, id)
	if err != nil {
		response.OK(c, gin.H{"id": id})
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteLevel(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "分销等级 ID 错误")
		return
	}
	row, err := h.loadLevel(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "分销等级不存在")
			return
		}
		failure(c)
		return
	}
	if row.PromoterCount > 0 {
		response.Fail(c, http.StatusBadRequest, "该等级下有数据，不能进行删除操作！")
		return
	}
	res := h.businessDB.WithContext(c.Request.Context()).Exec("DELETE FROM qixi_crm_b_distribution_level WHERE id = ?", id)
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "删除分销等级失败")
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *Handler) ListChildren(c *gin.Context) {
	parentID, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "分销员 ID 错误")
		return
	}
	page, limit := paging(c)
	levelRaw := strings.TrimSpace(c.Query("level"))
	levelFilter := int8(0) // 0=全部, 1=一级, 2=二级
	switch levelRaw {
	case "", "0", "all":
		levelFilter = 0
	case "1":
		levelFilter = 1
	case "2":
		levelFilter = 2
	default:
		response.Fail(c, http.StatusBadRequest, "用户类型错误")
		return
	}

	// L1 = parent_user_id = promoter; L2 = parent is a direct child of promoter.
	var l1IDs []uint64
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_relation").
		Where("parent_user_id = ?", parentID).Pluck("user_id", &l1IDs).Error; err != nil {
		failure(c)
		return
	}

	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_relation AS r").
		Joins("LEFT JOIN qixi_crm_b_user AS u ON u.id = r.user_id").
		Joins("LEFT JOIN qixi_crm_b_user_profile AS prof ON prof.user_id = r.user_id").
		Joins("LEFT JOIN qixi_crm_b_distribution_promoter AS p ON p.user_id = r.user_id").
		Joins(`LEFT JOIN (
			SELECT parent_user_id, COUNT(*) AS spread_user_count
			FROM qixi_crm_b_distribution_relation WHERE parent_user_id IS NOT NULL GROUP BY parent_user_id
		) AS child_rel ON child_rel.parent_user_id = r.user_id`).
		Joins(`LEFT JOIN (
			SELECT user_id,
				COUNT(*) AS pay_count,
				COALESCE(SUM(pay_amount),0) AS pay_amount
			FROM qixi_crm_b_order
			WHERE status IN ('paid','awaiting_final','fulfilling','shipped','completed','aftersale')
			GROUP BY user_id
		) AS ord ON ord.user_id = r.user_id`)

	switch levelFilter {
	case 1:
		q = q.Where("r.parent_user_id = ?", parentID)
	case 2:
		if len(l1IDs) == 0 {
			response.OK(c, gin.H{"list": []childRow{}, "total": 0, "page": page, "limit": limit})
			return
		}
		q = q.Where("r.parent_user_id IN ?", l1IDs)
	default:
		parentSet := append([]uint64{parentID}, l1IDs...)
		q = q.Where("r.parent_user_id IN ?", parentSet)
	}

	if dateFrom := strings.TrimSpace(c.Query("date_from")); dateFrom != "" {
		q = q.Where("r.bound_at >= ?", dateFrom+" 00:00:00")
	}
	if dateTo := strings.TrimSpace(c.Query("date_to")); dateTo != "" {
		q = q.Where("r.bound_at <= ?", dateTo+" 23:59:59")
	}
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		if id, err := strconv.ParseUint(keyword, 10, 64); err == nil && id > 0 {
			q = q.Where("(r.user_id = ? OR u.nickname LIKE ? OR u.mobile LIKE ?)", id, "%"+keyword+"%", "%"+keyword+"%")
		} else {
			q = q.Where("(u.nickname LIKE ? OR u.mobile LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
		}
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		failure(c)
		return
	}

	orderSQL := "r.bound_at DESC,r.user_id DESC"
	switch strings.TrimSpace(c.Query("sort_field")) {
	case "spread_user_count":
		dir := "DESC"
		if strings.EqualFold(strings.TrimSpace(c.Query("sort_order")), "asc") {
			dir = "ASC"
		}
		orderSQL = "spread_user_count " + dir + ",r.user_id DESC"
	}

	rows := make([]childRow, 0)
	selectSQL := `
r.user_id,COALESCE(u.nickname,'') AS nickname,COALESCE(u.mobile,'') AS mobile,
COALESCE(prof.avatar_url,'') AS avatar_url,r.bound_at,
CASE WHEN p.status = 1 THEN 1 ELSE 0 END AS is_promoter,
COALESCE(child_rel.spread_user_count,0) AS spread_user_count,
COALESCE(ord.pay_count,0) AS pay_count,
COALESCE(ord.pay_amount,0) AS pay_amount,
CASE WHEN r.parent_user_id = ` + strconv.FormatUint(parentID, 10) + ` THEN 1 ELSE 2 END AS level`
	if err := q.Select(selectSQL).
		Order(orderSQL).Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		failure(c)
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) ListSpreadOrders(c *gin.Context) {
	promoterID, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "分销员 ID 错误")
		return
	}
	page, limit := paging(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_commission_ledger AS c").
		Joins("LEFT JOIN qixi_crm_b_order AS o ON o.id = c.order_id").
		Joins("LEFT JOIN qixi_crm_b_user AS u ON u.id = o.user_id").
		Where("c.user_id = ? AND c.order_id > 0 AND c.status <> ?", promoterID, "voided")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		failure(c)
		return
	}
	rows := make([]spreadOrderRow, 0)
	if err := q.Select(`c.order_id,COALESCE(o.order_no,'') AS order_no,COALESCE(o.pay_amount,0) AS pay_amount,
COALESCE(c.amount,0) AS commission,COALESCE(c.status,'') AS status,
COALESCE(o.user_id,0) AS buyer_id,COALESCE(u.nickname,'') AS buyer_name,o.paid_at,c.created_at`).
		Order("c.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		failure(c)
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

type clearParentInput struct {
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

type updateLevelInput struct {
	LevelID        uint64 `json:"level_id"`
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

var (
	errPromoterNotFound = errors.New("promoter not found")
	errPromoterConflict = errors.New("promoter mutation conflict")
)

func (h *Handler) ClearParent(c *gin.Context) {
	userID, ok := positiveID(c.Param("id"))
	var in clearParentInput
	if !ok || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "清除上级参数错误")
		return
	}
	in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 || len([]rune(in.IdempotencyKey)) < 8 || len([]rune(in.IdempotencyKey)) > 128 {
		response.Fail(c, http.StatusBadRequest, "原因或幂等键错误")
		return
	}
	out := gin.H{}
	err := h.businessDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var locked struct{ ID uint64 }
		if e := tx.Table("qixi_crm_b_user").Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").Where("id = ?", userID).Take(&locked).Error; e != nil {
			if errors.Is(e, gorm.ErrRecordNotFound) {
				return errPromoterNotFound
			}
			return e
		}
		var replay struct {
			ID         uint64 `gorm:"column:id"`
			Reason     string `gorm:"column:reason"`
			OperatorID uint64 `gorm:"column:operator_admin_id"`
		}
		if e := tx.Table("qixi_crm_b_distribution_relation_audit").Where("user_id=? AND idempotency_key=?", userID, in.IdempotencyKey).Take(&replay).Error; e == nil {
			if replay.Reason != in.Reason || replay.OperatorID != uint64(middleware.AdminID(c)) {
				return errPromoterConflict
			}
			out = gin.H{"user_id": userID, "parent_user_id": 0, "replayed": true}
			return nil
		} else if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}
		if e := tx.Exec("INSERT INTO qixi_crm_b_distribution_relation (user_id,parent_user_id,bound_at) VALUES (?,NULL,NOW()) ON DUPLICATE KEY UPDATE user_id=VALUES(user_id)", userID).Error; e != nil {
			return e
		}
		var relation struct {
			ParentUserID *uint64 `gorm:"column:parent_user_id"`
		}
		if e := tx.Table("qixi_crm_b_distribution_relation").Clauses(clause.Locking{Strength: "UPDATE"}).Select("parent_user_id").Where("user_id=?", userID).Take(&relation).Error; e != nil {
			return e
		}
		var previous any
		if relation.ParentUserID != nil {
			previous = *relation.ParentUserID
		}
		if e := tx.Table("qixi_crm_b_distribution_relation").Where("user_id=?", userID).Updates(map[string]any{"parent_user_id": nil, "bound_at": time.Now()}).Error; e != nil {
			return e
		}
		if e := tx.Table("qixi_crm_b_distribution_relation_audit").Create(map[string]any{
			"user_id": userID, "previous_parent_user_id": previous, "parent_user_id": nil,
			"reason": in.Reason, "operator_admin_id": middleware.AdminID(c), "idempotency_key": in.IdempotencyKey,
		}).Error; e != nil {
			return e
		}
		out = gin.H{"user_id": userID, "previous_parent_user_id": previous, "parent_user_id": 0, "replayed": false}
		return nil
	})
	switch {
	case err == nil:
		response.OK(c, out)
	case errors.Is(err, errPromoterNotFound):
		response.Fail(c, http.StatusNotFound, "用户不存在")
	case errors.Is(err, errPromoterConflict):
		response.Fail(c, http.StatusConflict, "清除上级幂等键冲突")
	default:
		response.Fail(c, http.StatusInternalServerError, "清除上级推广人失败")
	}
}

func (h *Handler) UpdatePromoterLevel(c *gin.Context) {
	userID, ok := positiveID(c.Param("id"))
	var in updateLevelInput
	if !ok || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "等级参数错误")
		return
	}
	in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 || len([]rune(in.IdempotencyKey)) < 8 || len([]rune(in.IdempotencyKey)) > 128 {
		response.Fail(c, http.StatusBadRequest, "原因或幂等键错误")
		return
	}
	if in.LevelID > 0 {
		var level struct {
			ID uint64 `gorm:"column:id"`
		}
		if e := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_level").
			Select("id").Where("id = ? AND status = 1", in.LevelID).Take(&level).Error; e != nil {
			response.Fail(c, http.StatusBadRequest, "分销等级不存在")
			return
		}
	}
	res := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_promoter").
		Where("user_id = ?", userID).Updates(map[string]any{"level_id": in.LevelID, "updated_at": time.Now()})
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "更新分销员等级失败")
		return
	}
	if res.RowsAffected == 0 {
		// Upsert promoter row if missing (admin editing a user who is not yet in promoter table).
		if e := h.businessDB.WithContext(c.Request.Context()).Exec(
			"INSERT INTO qixi_crm_b_distribution_promoter (user_id,status,level_id,updated_at) VALUES (?,1,?,NOW()) ON DUPLICATE KEY UPDATE level_id=VALUES(level_id),updated_at=VALUES(updated_at)",
			userID, in.LevelID,
		).Error; e != nil {
			response.Fail(c, http.StatusInternalServerError, "更新分销员等级失败")
			return
		}
	}
	response.OK(c, gin.H{"user_id": userID, "level_id": in.LevelID, "reason": in.Reason, "idempotency_key": in.IdempotencyKey})
}

func (h *Handler) loadLevel(c *gin.Context, id uint64) (levelRow, error) {
	var row levelRow
	err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_level AS lv").
		Joins(`LEFT JOIN (
			SELECT level_id, COUNT(*) AS promoter_count
			FROM qixi_crm_b_distribution_promoter WHERE level_id > 0 GROUP BY level_id
		) AS pc ON pc.level_id = lv.id`).
		Select(`lv.id,lv.name,lv.`+"`rank`"+`,COALESCE(lv.icon_url,'') AS icon_url,
COALESCE(CAST(lv.task_rule AS CHAR),'') AS task_rule,
COALESCE(lv.extension_one,0) AS extension_one,COALESCE(lv.extension_two,0) AS extension_two,
lv.status,COALESCE(pc.promoter_count,0) AS promoter_count`).
		Where("lv.id = ?", id).Take(&row).Error
	if err != nil {
		return row, err
	}
	row.TaskRule = parseLevelTaskRule(row.TaskRuleRaw)
	return row, nil
}

func (h *Handler) levelRankExists(c *gin.Context, rank int, excludeID uint64) (bool, error) {
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_level").Where("`rank` = ?", rank)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func normalizeLevelSave(in levelSaveInput) (levelSaveInput, string) {
	in.Name = strings.TrimSpace(in.Name)
	in.IconURL = strings.TrimSpace(in.IconURL)
	if in.Name == "" {
		return in, "请输入等级名称"
	}
	if len([]rune(in.Name)) > 64 {
		return in, "等级名称过长"
	}
	if in.Rank < 1 {
		return in, "请输入等级"
	}
	if in.ExtensionOne < 0 || in.ExtensionOne > 1000 || in.ExtensionTwo < 0 || in.ExtensionTwo > 1000 {
		return in, "返佣上浮比例须在 0-1000 之间"
	}
	in.TaskRule = normalizeTaskRule(in.TaskRule)
	if !hasLevelTask(in.TaskRule) {
		return in, "请至少设置一个升级任务"
	}
	if msg := validateTaskPair(in.TaskRule); msg != "" {
		return in, msg
	}
	if in.Status != nil && *in.Status != 0 && *in.Status != 1 {
		return in, "状态错误"
	}
	return in, ""
}

func normalizeTaskRule(rule levelTaskRule) levelTaskRule {
	rule.SpreadUser = normalizeTaskItem(rule.SpreadUser)
	rule.PayMoney = normalizeTaskItem(rule.PayMoney)
	rule.PayNum = normalizeTaskItem(rule.PayNum)
	rule.SpreadMoney = normalizeTaskItem(rule.SpreadMoney)
	rule.SpreadPayNum = normalizeTaskItem(rule.SpreadPayNum)
	return rule
}

func normalizeTaskItem(item levelTaskItem) levelTaskItem {
	item.Name = strings.TrimSpace(item.Name)
	item.Info = strings.TrimSpace(item.Info)
	if item.Num < 0 {
		item.Num = 0
	}
	return item
}

func hasLevelTask(rule levelTaskRule) bool {
	return rule.SpreadUser.Num > 0 || rule.PayMoney.Num > 0 || rule.PayNum.Num > 0 ||
		rule.SpreadMoney.Num > 0 || rule.SpreadPayNum.Num > 0
}

func validateTaskPair(rule levelTaskRule) string {
	check := func(item levelTaskItem) bool {
		hasNum := item.Num > 0
		hasName := item.Name != ""
		return hasNum == hasName
	}
	if !check(rule.SpreadUser) || !check(rule.PayMoney) || !check(rule.PayNum) ||
		!check(rule.SpreadMoney) || !check(rule.SpreadPayNum) {
		return "请输入相对应的任务或数量"
	}
	return ""
}

func parseLevelTaskRule(raw string) levelTaskRule {
	raw = strings.TrimSpace(raw)
	out := levelTaskRule{}
	if raw == "" || raw == "null" {
		return out
	}
	_ = json.Unmarshal([]byte(raw), &out)
	return normalizeTaskRule(out)
}

func isMissingColumn(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unknown column") || strings.Contains(msg, "1054")
}

func isDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "1062")
}

func (h *Handler) promoterQuery(c *gin.Context) *gorm.DB {
	return h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_promoter AS p").
		Joins("LEFT JOIN qixi_crm_b_user AS u ON u.id = p.user_id").
		Joins("LEFT JOIN qixi_crm_b_user_profile AS prof ON prof.user_id = p.user_id").
		Joins("LEFT JOIN qixi_crm_b_distribution_level AS lv ON lv.id = p.level_id").
		Joins("LEFT JOIN qixi_crm_b_distribution_relation AS self_rel ON self_rel.user_id = p.user_id").
		Joins("LEFT JOIN qixi_crm_b_user AS parent ON parent.id = self_rel.parent_user_id").
		Joins(`LEFT JOIN (
			SELECT user_id,
				COALESCE(SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END),0) AS pending_commission,
				COALESCE(SUM(CASE WHEN status = 'available' THEN amount ELSE 0 END),0) AS available_commission,
				COALESCE(SUM(CASE WHEN status = 'settled' THEN amount ELSE 0 END),0) AS settled_commission
			FROM qixi_crm_b_commission_ledger GROUP BY user_id
		) AS ledger ON ledger.user_id = p.user_id`).
		Joins(`LEFT JOIN (
			SELECT parent_user_id, COUNT(*) AS direct_user_count
			FROM qixi_crm_b_distribution_relation WHERE parent_user_id IS NOT NULL GROUP BY parent_user_id
		) AS rel ON rel.parent_user_id = p.user_id`).
		Joins(`LEFT JOIN (
			SELECT c.user_id,
				COUNT(DISTINCT c.order_id) AS spread_order_count,
				COALESCE(SUM(o.pay_amount),0) AS spread_order_amount
			FROM qixi_crm_b_commission_ledger AS c
			LEFT JOIN qixi_crm_b_order AS o ON o.id = c.order_id
			WHERE c.order_id > 0 AND c.status <> 'voided'
			GROUP BY c.user_id
		) AS ord ON ord.user_id = p.user_id`).
		Joins(`LEFT JOIN (
			SELECT user_id,
				COALESCE(SUM(CASE WHEN status IN ('paid','approved') THEN amount ELSE 0 END),0) AS withdrawn_amount,
				COALESCE(SUM(CASE WHEN status IN ('paid','approved') THEN 1 ELSE 0 END),0) AS withdraw_count
			FROM qixi_crm_b_withdrawal_application GROUP BY user_id
		) AS wd ON wd.user_id = p.user_id`)
}

func (h *Handler) applyPromoterFilters(c *gin.Context, q *gorm.DB) (*gorm.DB, bool) {
	if status, ok := promoterStatus(c.Query("status")); !ok {
		response.Fail(c, http.StatusBadRequest, "推广员状态错误")
		return q, false
	} else if status != nil {
		q = q.Where("p.status = ?", *status)
	}
	if userID, provided, ok := queryID(c, "user_id"); !ok {
		response.Fail(c, http.StatusBadRequest, "用户 ID 错误")
		return q, false
	} else if provided {
		q = q.Where("p.user_id = ?", userID)
	}
	if levelID, provided, ok := queryID(c, "level_id"); !ok {
		response.Fail(c, http.StatusBadRequest, "等级 ID 错误")
		return q, false
	} else if provided {
		q = q.Where("p.level_id = ?", levelID)
	}
	if dateFrom := strings.TrimSpace(c.Query("date_from")); dateFrom != "" {
		q = q.Where("p.updated_at >= ?", dateFrom+" 00:00:00")
	}
	if dateTo := strings.TrimSpace(c.Query("date_to")); dateTo != "" {
		q = q.Where("p.updated_at <= ?", dateTo+" 23:59:59")
	}
	keyword := strings.TrimSpace(c.Query("keyword"))
	if keyword == "" {
		return q, true
	}
	switch strings.TrimSpace(c.Query("keyword_type")) {
	case "", "nickname":
		q = q.Where("u.nickname LIKE ?", "%"+keyword+"%")
	case "uid", "user_id":
		id, err := strconv.ParseUint(keyword, 10, 64)
		if err != nil || id == 0 {
			response.Fail(c, http.StatusBadRequest, "用户 ID 错误")
			return q, false
		}
		q = q.Where("p.user_id = ?", id)
	case "phone", "mobile":
		q = q.Where("u.mobile LIKE ?", "%"+keyword+"%")
	default:
		response.Fail(c, http.StatusBadRequest, "搜索类型错误")
		return q, false
	}
	return q, true
}

func promoterOrderSQL(field, order string) string {
	dir := "DESC"
	if strings.EqualFold(strings.TrimSpace(order), "asc") {
		dir = "ASC"
	}
	switch strings.TrimSpace(field) {
	case "commission_amount":
		return "commission_amount " + dir + ",p.user_id DESC"
	case "withdrawn_amount":
		return "withdrawn_amount " + dir + ",p.user_id DESC"
	case "unwithdrawn_amount":
		return "unwithdrawn_amount " + dir + ",p.user_id DESC"
	default:
		return "p.updated_at DESC,p.user_id DESC"
	}
}

func promoterStatus(raw string) (*int8, bool) {
	switch strings.TrimSpace(raw) {
	case "":
		return nil, true
	case "0":
		value := int8(0)
		return &value, true
	case "1":
		value := int8(1)
		return &value, true
	default:
		return nil, false
	}
}

func commissionStatus(raw string) (string, bool) {
	switch strings.TrimSpace(raw) {
	case "", "pending", "available", "settled", "voided":
		return strings.TrimSpace(raw), true
	default:
		return "", false
	}
}

func queryID(c *gin.Context, field string) (uint64, bool, bool) {
	raw := strings.TrimSpace(c.Query(field))
	if raw == "" {
		return 0, false, true
	}
	value, err := strconv.ParseUint(raw, 10, 64)
	return value, true, err == nil && value > 0
}

func positiveID(raw string) (uint64, bool) {
	value, err := strconv.ParseUint(strings.TrimSpace(raw), 10, 64)
	return value, err == nil && value > 0
}

func paging(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}

func failure(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "分销监管数据查询失败")
}

func isMissingTable(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "doesn't exist") || strings.Contains(msg, "1146")
}
