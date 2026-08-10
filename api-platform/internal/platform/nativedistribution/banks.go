package nativedistribution

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type withdrawBankRow struct {
	ID        uint64    `gorm:"column:id" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Status    int8      `gorm:"column:status" json:"status"`
	Sort      int       `gorm:"column:sort" json:"sort"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type withdrawBankSaveInput struct {
	Name   string `json:"name"`
	Status *int8  `json:"status"`
	Sort   *int   `json:"sort"`
}

func (h *Handler) ListWithdrawBanks(c *gin.Context) {
	page, limit := paging(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_withdraw_bank")
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		q = q.Where("name LIKE ?", "%"+keyword+"%")
	}
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		switch status {
		case "0", "1":
			q = q.Where("status = ?", status)
		default:
			response.Fail(c, http.StatusBadRequest, "状态错误")
			return
		}
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		if isMissingTable(err) {
			response.OK(c, gin.H{"list": []withdrawBankRow{}, "total": 0, "page": page, "limit": limit})
			return
		}
		failure(c)
		return
	}
	rows := make([]withdrawBankRow, 0)
	if err := q.Select("id,name,status,sort,created_at,updated_at").
		Order("sort ASC, id DESC").
		Offset((page - 1) * limit).Limit(limit).
		Scan(&rows).Error; err != nil {
		failure(c)
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) CreateWithdrawBank(c *gin.Context) {
	var in withdrawBankSaveInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "提现银行参数错误")
		return
	}
	name, status, sort, msg := normalizeWithdrawBankSave(in, 1, 0)
	if msg != "" {
		response.Fail(c, http.StatusBadRequest, msg)
		return
	}
	res := h.businessDB.WithContext(c.Request.Context()).Exec(
		"INSERT INTO qixi_crm_b_withdraw_bank (`name`,`status`,`sort`) VALUES (?,?,?)",
		name, status, sort,
	)
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "创建提现银行失败")
		return
	}
	var id uint64
	_ = h.businessDB.WithContext(c.Request.Context()).Raw("SELECT LAST_INSERT_ID()").Scan(&id).Error
	row, err := h.loadWithdrawBank(c, id)
	if err != nil {
		response.OK(c, gin.H{"id": id, "name": name, "status": status, "sort": sort})
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateWithdrawBank(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "提现银行 ID 错误")
		return
	}
	existing, err := h.loadWithdrawBank(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "提现银行不存在")
			return
		}
		failure(c)
		return
	}
	var in withdrawBankSaveInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "提现银行参数错误")
		return
	}
	name, status, sort, msg := normalizeWithdrawBankSave(in, existing.Status, existing.Sort)
	if msg != "" {
		response.Fail(c, http.StatusBadRequest, msg)
		return
	}
	res := h.businessDB.WithContext(c.Request.Context()).Exec(
		"UPDATE qixi_crm_b_withdraw_bank SET name=?, status=?, sort=?, updated_at=NOW() WHERE id=?",
		name, status, sort, id,
	)
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "更新提现银行失败")
		return
	}
	row, err := h.loadWithdrawBank(c, id)
	if err != nil {
		response.OK(c, gin.H{"id": id, "name": name, "status": status, "sort": sort})
		return
	}
	response.OK(c, row)
}

func (h *Handler) SetWithdrawBankStatus(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "提现银行 ID 错误")
		return
	}
	if _, err := h.loadWithdrawBank(c, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "提现银行不存在")
			return
		}
		failure(c)
		return
	}
	var in struct {
		Status *int8 `json:"status"`
	}
	if c.ShouldBindJSON(&in) != nil || in.Status == nil || (*in.Status != 0 && *in.Status != 1) {
		response.Fail(c, http.StatusBadRequest, "状态错误")
		return
	}
	res := h.businessDB.WithContext(c.Request.Context()).Exec(
		"UPDATE qixi_crm_b_withdraw_bank SET status=?, updated_at=NOW() WHERE id=?",
		*in.Status, id,
	)
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "更新提现银行状态失败")
		return
	}
	row, err := h.loadWithdrawBank(c, id)
	if err != nil {
		response.OK(c, gin.H{"id": id, "status": *in.Status})
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteWithdrawBank(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "提现银行 ID 错误")
		return
	}
	if _, err := h.loadWithdrawBank(c, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "提现银行不存在")
			return
		}
		failure(c)
		return
	}
	res := h.businessDB.WithContext(c.Request.Context()).Exec("DELETE FROM qixi_crm_b_withdraw_bank WHERE id = ?", id)
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "删除提现银行失败")
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *Handler) loadWithdrawBank(c *gin.Context, id uint64) (withdrawBankRow, error) {
	var row withdrawBankRow
	err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_withdraw_bank").
		Select("id,name,status,sort,created_at,updated_at").
		Where("id = ?", id).Take(&row).Error
	return row, err
}

func normalizeWithdrawBankSave(in withdrawBankSaveInput, defaultStatus int8, defaultSort int) (string, int8, int, string) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return "", 0, 0, "请输入银行名称"
	}
	if utf8.RuneCountInString(name) > 64 {
		return "", 0, 0, "银行名称过长"
	}
	status := defaultStatus
	if in.Status != nil {
		if *in.Status != 0 && *in.Status != 1 {
			return "", 0, 0, "状态错误"
		}
		status = *in.Status
	}
	sort := defaultSort
	if in.Sort != nil {
		sort = *in.Sort
	}
	return name, status, sort, ""
}
