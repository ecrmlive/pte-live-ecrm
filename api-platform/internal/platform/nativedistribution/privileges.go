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

type privilegeRow struct {
	ID        uint64    `gorm:"column:id" json:"id"`
	Title     string    `gorm:"column:title" json:"title"`
	ImgURL    string    `gorm:"column:img_url" json:"img_url"`
	Status    int8      `gorm:"column:status" json:"status"`
	Sort      int       `gorm:"column:sort" json:"sort"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type privilegeSaveInput struct {
	Title  string `json:"title"`
	ImgURL string `json:"img_url"`
	Status *int8  `json:"status"`
	Sort   *int   `json:"sort"`
}

func (h *Handler) ListPrivileges(c *gin.Context) {
	page, limit := paging(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_privilege")
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		q = q.Where("title LIKE ?", "%"+keyword+"%")
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
			response.OK(c, gin.H{"list": []privilegeRow{}, "total": 0, "page": page, "limit": limit})
			return
		}
		failure(c)
		return
	}
	rows := make([]privilegeRow, 0)
	if err := q.Select("id,title,img_url,status,sort,created_at,updated_at").
		Order("sort ASC, id DESC").
		Offset((page - 1) * limit).Limit(limit).
		Scan(&rows).Error; err != nil {
		failure(c)
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) CreatePrivilege(c *gin.Context) {
	var in privilegeSaveInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "分销特权参数错误")
		return
	}
	title, imgURL, status, sort, msg := normalizePrivilegeSave(in, 1, 0)
	if msg != "" {
		response.Fail(c, http.StatusBadRequest, msg)
		return
	}
	res := h.businessDB.WithContext(c.Request.Context()).Exec(
		"INSERT INTO qixi_crm_b_distribution_privilege (`title`,`img_url`,`status`,`sort`) VALUES (?,?,?,?)",
		title, imgURL, status, sort,
	)
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "创建分销特权失败")
		return
	}
	var id uint64
	_ = h.businessDB.WithContext(c.Request.Context()).Raw("SELECT LAST_INSERT_ID()").Scan(&id).Error
	row, err := h.loadPrivilege(c, id)
	if err != nil {
		response.OK(c, gin.H{"id": id, "title": title, "img_url": imgURL, "status": status, "sort": sort})
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdatePrivilege(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "分销特权 ID 错误")
		return
	}
	existing, err := h.loadPrivilege(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "分销特权不存在")
			return
		}
		failure(c)
		return
	}
	var in privilegeSaveInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "分销特权参数错误")
		return
	}
	title, imgURL, status, sort, msg := normalizePrivilegeSave(in, existing.Status, existing.Sort)
	if msg != "" {
		response.Fail(c, http.StatusBadRequest, msg)
		return
	}
	res := h.businessDB.WithContext(c.Request.Context()).Exec(
		"UPDATE qixi_crm_b_distribution_privilege SET title=?, img_url=?, status=?, sort=?, updated_at=NOW() WHERE id=?",
		title, imgURL, status, sort, id,
	)
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "更新分销特权失败")
		return
	}
	row, err := h.loadPrivilege(c, id)
	if err != nil {
		response.OK(c, gin.H{"id": id, "title": title, "img_url": imgURL, "status": status, "sort": sort})
		return
	}
	response.OK(c, row)
}

func (h *Handler) SetPrivilegeStatus(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "分销特权 ID 错误")
		return
	}
	if _, err := h.loadPrivilege(c, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "分销特权不存在")
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
		"UPDATE qixi_crm_b_distribution_privilege SET status=?, updated_at=NOW() WHERE id=?",
		*in.Status, id,
	)
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "更新分销特权状态失败")
		return
	}
	row, err := h.loadPrivilege(c, id)
	if err != nil {
		response.OK(c, gin.H{"id": id, "status": *in.Status})
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeletePrivilege(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "分销特权 ID 错误")
		return
	}
	if _, err := h.loadPrivilege(c, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "分销特权不存在")
			return
		}
		failure(c)
		return
	}
	res := h.businessDB.WithContext(c.Request.Context()).Exec(
		"DELETE FROM qixi_crm_b_distribution_privilege WHERE id = ?", id,
	)
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "删除分销特权失败")
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *Handler) loadPrivilege(c *gin.Context, id uint64) (privilegeRow, error) {
	var row privilegeRow
	err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_privilege").
		Select("id,title,img_url,status,sort,created_at,updated_at").
		Where("id = ?", id).Take(&row).Error
	return row, err
}

func normalizePrivilegeSave(in privilegeSaveInput, defaultStatus int8, defaultSort int) (string, string, int8, int, string) {
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return "", "", 0, 0, "请输入标题"
	}
	if utf8.RuneCountInString(title) > 20 {
		return "", "", 0, 0, "标题过长"
	}
	imgURL := strings.TrimSpace(in.ImgURL)
	if imgURL == "" {
		return "", "", 0, 0, "请选择图片"
	}
	if utf8.RuneCountInString(imgURL) > 1024 {
		return "", "", 0, 0, "图片地址过长"
	}
	status := defaultStatus
	if in.Status != nil {
		if *in.Status != 0 && *in.Status != 1 {
			return "", "", 0, 0, "状态错误"
		}
		status = *in.Status
	}
	sort := defaultSort
	if in.Sort != nil {
		sort = *in.Sort
	}
	return title, imgURL, status, sort, ""
}
