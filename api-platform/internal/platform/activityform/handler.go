// Package activityform manages platform signup activities
// (CRMEB admin.store.marketing.StoreForm / ACTIVITY_TYPE_FORM).
package activityform

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	menuRead   = "marketing.application.read"
	menuManage = "marketing.application.manage"
	exportCap  = 5000
)

type Handler struct {
	adminDB *gorm.DB
}

func NewHandler(adminDB *gorm.DB) *Handler { return &Handler{adminDB: adminDB} }

func (h *Handler) Register(r gin.IRoutes) {
	access := middleware.RequireAdminRoles("platform", "operations")
	read := middleware.RequireAdminMenu(h.adminDB, menuRead)
	manage := middleware.RequireAdminMenu(h.adminDB, menuManage)

	r.GET("/marketing/applications/form-options", access, read, h.formOptions)
	r.GET("/marketing/applications", access, read, h.list)
	r.POST("/marketing/applications", access, manage, h.create)
	r.GET("/marketing/applications/:id", access, read, h.detail)
	r.PUT("/marketing/applications/:id", access, manage, h.update)
	r.PUT("/marketing/applications/:id/status", access, manage, h.setStatus)
	r.DELETE("/marketing/applications/:id", access, manage, h.remove)
	r.GET("/marketing/applications/:id/users", access, read, h.userList)
	r.GET("/marketing/applications/:id/users/export", access, read, h.userExport)
}

type activityRow struct {
	ID          uint64     `gorm:"column:id"`
	Name        string     `gorm:"column:name"`
	Info        string     `gorm:"column:info"`
	CoverURL    string     `gorm:"column:cover_url"`
	PosterURL   string     `gorm:"column:poster_url"`
	Color       string     `gorm:"column:color"`
	FormID      uint64     `gorm:"column:form_id"`
	Quota       uint32     `gorm:"column:quota"`
	Total       uint32     `gorm:"column:total"`
	Status      int        `gorm:"column:status"`
	Sort        int        `gorm:"column:sort"`
	StartsAt    *time.Time `gorm:"column:starts_at"`
	EndsAt      *time.Time `gorm:"column:ends_at"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	FormName    string     `gorm:"column:form_name"`
	FormPayload string     `gorm:"column:form_payload"`
}

type recordRow struct {
	ID         uint64    `gorm:"column:id"`
	ActivityID uint64    `gorm:"column:activity_id"`
	UserID     uint64    `gorm:"column:user_id"`
	Nickname   string    `gorm:"column:nickname"`
	Mobile     string    `gorm:"column:mobile"`
	Avatar     string    `gorm:"column:avatar"`
	FormValue  string    `gorm:"column:form_value"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

type formField struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

type upsertInput struct {
	Name      string `json:"name"`
	Info      string `json:"info"`
	CoverURL  string `json:"cover_url"`
	PosterURL string `json:"poster_url"`
	Color     string `json:"color"`
	FormID    uint64 `json:"form_id"`
	Quota     *int   `json:"quota"`
	Status    *int   `json:"status"`
	Sort      *int   `json:"sort"`
	StartsAt  string `json:"starts_at"`
	EndsAt    string `json:"ends_at"`
}

func (h *Handler) list(c *gin.Context) {
	page, limit := pageLimit(c)
	q := h.adminDB.WithContext(c.Request.Context()).
		Table("qixi_crm_a_signup_activity AS a").
		Joins("LEFT JOIN qixi_crm_a_config_item AS f ON f.id = a.form_id AND f.item_type = 'system_form' AND f.is_del = 0").
		Where("a.is_del = 0")

	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("a.name LIKE ? OR a.info LIKE ? OR CAST(a.id AS CHAR) LIKE ?", like, like, like)
	}
	if formID := strings.TrimSpace(c.Query("form_id")); formID != "" {
		id, err := strconv.ParseUint(formID, 10, 64)
		if err != nil || id == 0 {
			response.Fail(c, http.StatusBadRequest, "关联表单错误")
			return
		}
		q = q.Where("a.form_id = ?", id)
	}
	// 报名状态：0 未开始 / 1 进行中 / -1 已结束（对齐 CRMEB time_status）
	if raw := strings.TrimSpace(c.Query("activity_status")); raw != "" {
		now := time.Now()
		switch raw {
		case "0":
			q = q.Where("a.starts_at IS NOT NULL AND a.starts_at > ?", now)
		case "1":
			q = q.Where("(a.starts_at IS NULL OR a.starts_at <= ?) AND (a.ends_at IS NULL OR a.ends_at >= ?)", now, now)
		case "-1":
			q = q.Where("a.ends_at IS NOT NULL AND a.ends_at < ?", now)
		default:
			response.Fail(c, http.StatusBadRequest, "报名状态错误")
			return
		}
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c)
		return
	}
	rows := make([]activityRow, 0)
	if err := q.Select(`a.id,a.name,a.info,a.cover_url,a.poster_url,a.color,a.form_id,a.quota,a.total,
		a.status,a.sort,a.starts_at,a.ends_at,a.created_at,a.updated_at,
		IFNULL(f.name,'') AS form_name, IFNULL(f.payload,'{}') AS form_payload`).
		Order("a.sort DESC, a.id DESC").
		Offset((page - 1) * limit).Limit(limit).
		Scan(&rows).Error; err != nil {
		fail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, toActivityItem(row, false))
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) detail(c *gin.Context) {
	row, ok := h.load(c, c.Param("id"), true)
	if !ok {
		return
	}
	response.OK(c, toActivityItem(*row, true))
}

func (h *Handler) create(c *gin.Context) {
	var in upsertInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	name, info, cover, poster, color, formID, quota, status, sort, starts, ends, err := normalize(in, 1, 0)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if !h.formExists(c, formID) {
		response.Fail(c, http.StatusBadRequest, "关联系统表单不存在或已停用")
		return
	}
	created := struct {
		ID        uint64     `gorm:"column:id;primaryKey"`
		Name      string     `gorm:"column:name"`
		Info      string     `gorm:"column:info"`
		CoverURL  string     `gorm:"column:cover_url"`
		PosterURL string     `gorm:"column:poster_url"`
		Color     string     `gorm:"column:color"`
		FormID    uint64     `gorm:"column:form_id"`
		Quota     int        `gorm:"column:quota"`
		Total     int        `gorm:"column:total"`
		Status    int        `gorm:"column:status"`
		Sort      int        `gorm:"column:sort"`
		StartsAt  *time.Time `gorm:"column:starts_at"`
		EndsAt    *time.Time `gorm:"column:ends_at"`
		IsDel     int        `gorm:"column:is_del"`
	}{
		Name: name, Info: info, CoverURL: cover, PosterURL: poster, Color: color,
		FormID: formID, Quota: quota, Total: 0, Status: status, Sort: sort,
		StartsAt: starts, EndsAt: ends, IsDel: 0,
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_signup_activity").Create(&created).Error; err != nil {
		fail(c)
		return
	}
	row, ok := h.load(c, strconv.FormatUint(created.ID, 10), true)
	if !ok {
		return
	}
	response.OK(c, toActivityItem(*row, true))
}

func (h *Handler) update(c *gin.Context) {
	row, ok := h.load(c, c.Param("id"), false)
	if !ok {
		return
	}
	var in upsertInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	name, info, cover, poster, color, formID, quota, status, sort, starts, ends, err := normalize(in, row.Status, row.Sort)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if !h.formExists(c, formID) {
		response.Fail(c, http.StatusBadRequest, "关联系统表单不存在或已停用")
		return
	}
	res := h.adminDB.WithContext(c.Request.Context()).Exec(`
		UPDATE qixi_crm_a_signup_activity
		SET name=?, info=?, cover_url=?, poster_url=?, color=?, form_id=?, quota=?, status=?, sort=?, starts_at=?, ends_at=?
		WHERE id=? AND is_del=0`,
		name, info, cover, poster, color, formID, quota, status, sort, starts, ends, row.ID,
	)
	if res.Error != nil {
		fail(c)
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "记录不存在")
		return
	}
	updated, ok := h.load(c, strconv.FormatUint(row.ID, 10), true)
	if !ok {
		return
	}
	response.OK(c, toActivityItem(*updated, true))
}

func (h *Handler) setStatus(c *gin.Context) {
	row, ok := h.load(c, c.Param("id"), false)
	if !ok {
		return
	}
	var in struct {
		Status *int `json:"status"`
	}
	if c.ShouldBindJSON(&in) != nil || in.Status == nil || (*in.Status != 0 && *in.Status != 1) {
		response.Fail(c, http.StatusBadRequest, "状态错误")
		return
	}
	res := h.adminDB.WithContext(c.Request.Context()).Exec(
		`UPDATE qixi_crm_a_signup_activity SET status=? WHERE id=? AND is_del=0`,
		*in.Status, row.ID,
	)
	if res.Error != nil {
		fail(c)
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "记录不存在")
		return
	}
	response.OK(c, gin.H{"id": row.ID, "status": *in.Status})
}

func (h *Handler) remove(c *gin.Context) {
	row, ok := h.load(c, c.Param("id"), false)
	if !ok {
		return
	}
	if timeStatus(row.StartsAt, row.EndsAt, row.Status) == 1 {
		response.Fail(c, http.StatusBadRequest, "活动进行中不能删除")
		return
	}
	err := h.adminDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`UPDATE qixi_crm_a_signup_record SET is_del=1 WHERE activity_id=? AND is_del=0`, row.ID).Error; err != nil {
			return err
		}
		res := tx.Exec(`UPDATE qixi_crm_a_signup_activity SET is_del=1, status=0 WHERE id=? AND is_del=0`, row.ID)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "记录不存在")
		return
	}
	if err != nil {
		fail(c)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) formOptions(c *gin.Context) {
	type opt struct {
		ID      uint64 `gorm:"column:id"`
		Name    string `gorm:"column:name"`
		Payload string `gorm:"column:payload"`
	}
	rows := make([]opt, 0)
	err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_config_item").
		Select("id,name,payload").
		Where("item_type = 'system_form' AND is_del = 0 AND status = 1").
		Order("sort ASC, id DESC").
		Scan(&rows).Error
	if err != nil {
		fail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, gin.H{
			"id":     row.ID,
			"name":   row.Name,
			"fields": parseFormFields(row.Payload),
		})
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) userList(c *gin.Context) {
	row, ok := h.load(c, c.Param("id"), true)
	if !ok {
		return
	}
	page, limit := pageLimit(c)
	keyword := strings.TrimSpace(c.Query("keyword"))
	q := h.recordQuery(c, row.ID, keyword)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c)
		return
	}
	rows := make([]recordRow, 0)
	if err := q.Select("id,activity_id,user_id,nickname,mobile,avatar,form_value,created_at").
		Order("id DESC").
		Offset((page - 1) * limit).Limit(limit).
		Scan(&rows).Error; err != nil {
		fail(c)
		return
	}
	fields := parseFormFields(row.FormPayload)
	list := make([]gin.H, 0, len(rows))
	for i, rec := range rows {
		list = append(list, toRecordItem(rec, (page-1)*limit+i+1, fields))
	}
	response.OK(c, gin.H{
		"list":   list,
		"total":  total,
		"page":   page,
		"limit":  limit,
		"fields": fields,
	})
}

func (h *Handler) userExport(c *gin.Context) {
	row, ok := h.load(c, c.Param("id"), true)
	if !ok {
		return
	}
	keyword := strings.TrimSpace(c.Query("keyword"))
	q := h.recordQuery(c, row.ID, keyword)
	rows := make([]recordRow, 0)
	if err := q.Select("id,activity_id,user_id,nickname,mobile,avatar,form_value,created_at").
		Order("id DESC").Limit(exportCap).Scan(&rows).Error; err != nil {
		fail(c)
		return
	}
	fields := parseFormFields(row.FormPayload)
	buf := &bytes.Buffer{}
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(buf)
	header := []string{"序号", "用户名称", "用户ID", "手机号"}
	for _, f := range fields {
		header = append(header, f.Label)
	}
	header = append(header, "提交时间")
	_ = w.Write(header)
	for i, rec := range rows {
		vals := parseFormValue(rec.FormValue)
		line := []string{
			strconv.Itoa(i + 1),
			rec.Nickname,
			strconv.FormatUint(rec.UserID, 10),
			rec.Mobile,
		}
		for _, f := range fields {
			line = append(line, stringifyFormCell(vals[f.Key]))
		}
		line = append(line, formatTime(rec.CreatedAt))
		_ = w.Write(line)
	}
	w.Flush()
	name := "报名用户_" + sanitizeFilePart(row.Name) + ".csv"
	response.OK(c, gin.H{
		"file_name": name,
		"content":   buf.String(),
	})
}

func (h *Handler) recordQuery(c *gin.Context, activityID uint64, keyword string) *gorm.DB {
	q := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_signup_record").
		Where("activity_id = ? AND is_del = 0", activityID)
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("nickname LIKE ? OR mobile LIKE ? OR CAST(user_id AS CHAR) LIKE ?", like, like, like)
	}
	return q
}

func (h *Handler) load(c *gin.Context, rawID string, withForm bool) (*activityRow, bool) {
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "ID 错误")
		return nil, false
	}
	q := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_signup_activity AS a")
	selectCols := `a.id,a.name,a.info,a.cover_url,a.poster_url,a.color,a.form_id,a.quota,a.total,
		a.status,a.sort,a.starts_at,a.ends_at,a.created_at,a.updated_at,
		'' AS form_name, '{}' AS form_payload`
	if withForm {
		q = q.Joins("LEFT JOIN qixi_crm_a_config_item AS f ON f.id = a.form_id AND f.item_type = 'system_form' AND f.is_del = 0")
		selectCols = `a.id,a.name,a.info,a.cover_url,a.poster_url,a.color,a.form_id,a.quota,a.total,
			a.status,a.sort,a.starts_at,a.ends_at,a.created_at,a.updated_at,
			IFNULL(f.name,'') AS form_name, IFNULL(f.payload,'{}') AS form_payload`
	}
	var row activityRow
	err = q.Select(selectCols).Where("a.id = ? AND a.is_del = 0", id).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "记录不存在")
		return nil, false
	}
	if err != nil {
		fail(c)
		return nil, false
	}
	return &row, true
}

func (h *Handler) formExists(c *gin.Context, formID uint64) bool {
	var n int64
	_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_config_item").
		Where("id = ? AND item_type = 'system_form' AND is_del = 0 AND status = 1", formID).
		Count(&n).Error
	return n > 0
}

func normalize(in upsertInput, defaultStatus, defaultSort int) (string, string, string, string, string, uint64, int, int, int, *time.Time, *time.Time, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || utf8.RuneCountInString(name) > 128 {
		return "", "", "", "", "", 0, 0, 0, 0, nil, nil, errors.New("请输入活动名称")
	}
	info := strings.TrimSpace(in.Info)
	if utf8.RuneCountInString(info) > 500 {
		return "", "", "", "", "", 0, 0, 0, 0, nil, nil, errors.New("活动简介不超过 500 字")
	}
	cover := strings.TrimSpace(in.CoverURL)
	if cover == "" {
		return "", "", "", "", "", 0, 0, 0, 0, nil, nil, errors.New("请上传封面图")
	}
	if utf8.RuneCountInString(cover) > 1024 {
		return "", "", "", "", "", 0, 0, 0, 0, nil, nil, errors.New("封面图地址过长")
	}
	poster := strings.TrimSpace(in.PosterURL)
	if poster == "" {
		return "", "", "", "", "", 0, 0, 0, 0, nil, nil, errors.New("请上传活动分享海报")
	}
	if utf8.RuneCountInString(poster) > 1024 {
		return "", "", "", "", "", 0, 0, 0, 0, nil, nil, errors.New("分享海报地址过长")
	}
	color := strings.TrimSpace(in.Color)
	if utf8.RuneCountInString(color) > 32 {
		return "", "", "", "", "", 0, 0, 0, 0, nil, nil, errors.New("背景色格式错误")
	}
	if in.FormID == 0 {
		return "", "", "", "", "", 0, 0, 0, 0, nil, nil, errors.New("请关联系统表单")
	}
	starts, err := parseOptionalTime(in.StartsAt)
	if err != nil {
		return "", "", "", "", "", 0, 0, 0, 0, nil, nil, errors.New("开始时间格式错误")
	}
	ends, err := parseOptionalTime(in.EndsAt)
	if err != nil {
		return "", "", "", "", "", 0, 0, 0, 0, nil, nil, errors.New("结束时间格式错误")
	}
	if starts == nil || ends == nil {
		return "", "", "", "", "", 0, 0, 0, 0, nil, nil, errors.New("请选择活动起止日期")
	}
	if ends.Before(*starts) {
		return "", "", "", "", "", 0, 0, 0, 0, nil, nil, errors.New("结束时间不能早于开始时间")
	}
	quota := 0
	if in.Quota != nil {
		if *in.Quota < 0 {
			return "", "", "", "", "", 0, 0, 0, 0, nil, nil, errors.New("人数上限不能为负")
		}
		quota = *in.Quota
	}
	status := defaultStatus
	if in.Status != nil {
		if *in.Status != 0 && *in.Status != 1 {
			return "", "", "", "", "", 0, 0, 0, 0, nil, nil, errors.New("状态错误")
		}
		status = *in.Status
	}
	sort := defaultSort
	if in.Sort != nil {
		sort = *in.Sort
	}
	return name, info, cover, poster, color, in.FormID, quota, status, sort, starts, ends, nil
}

func toActivityItem(row activityRow, withFields bool) gin.H {
	ts := timeStatus(row.StartsAt, row.EndsAt, row.Status)
	item := gin.H{
		"id":                   row.ID,
		"name":                 row.Name,
		"info":                 row.Info,
		"cover_url":            row.CoverURL,
		"poster_url":           row.PosterURL,
		"color":                row.Color,
		"form_id":              row.FormID,
		"form_name":            row.FormName,
		"quota":                row.Quota,
		"total":                row.Total,
		"status":               row.Status,
		"sort":                 row.Sort,
		"starts_at":            formatOptTime(row.StartsAt),
		"ends_at":              formatOptTime(row.EndsAt),
		"created_at":           formatTime(row.CreatedAt),
		"updated_at":           formatTime(row.UpdatedAt),
		"activity_status":      ts,
		"activity_status_text": activityStatusText(ts),
		"signup_count_text":    signupCountText(row.Total, row.Quota),
	}
	if withFields {
		item["form_fields"] = parseFormFields(row.FormPayload)
	}
	return item
}

func toRecordItem(rec recordRow, index int, fields []formField) gin.H {
	vals := parseFormValue(rec.FormValue)
	formCols := gin.H{}
	for _, f := range fields {
		formCols[f.Key] = stringifyFormCell(vals[f.Key])
	}
	return gin.H{
		"id":          rec.ID,
		"index":       index,
		"activity_id": rec.ActivityID,
		"user_id":     rec.UserID,
		"nickname":    rec.Nickname,
		"mobile":      rec.Mobile,
		"avatar":      rec.Avatar,
		"form_value":  vals,
		"form_cols":   formCols,
		"created_at":  formatTime(rec.CreatedAt),
	}
}

func timeStatus(starts, ends *time.Time, _ int) int {
	// 报名状态仅按起止时间：0 未开始 / 1 进行中 / -1 已结束；是否显示用 status 字段。
	now := time.Now()
	if starts != nil && starts.After(now) {
		return 0
	}
	if ends != nil && ends.Before(now) {
		return -1
	}
	return 1
}

func activityStatusText(ts int) string {
	switch ts {
	case 0:
		return "未开始"
	case 1:
		return "进行中"
	case -1:
		return "已结束"
	default:
		return "已关闭"
	}
}

func signupCountText(total, quota uint32) string {
	if quota == 0 {
		return strconv.FormatUint(uint64(total), 10) + "/不限制"
	}
	return strconv.FormatUint(uint64(total), 10) + "/" + strconv.FormatUint(uint64(quota), 10)
}

func parseFormFields(payload string) []formField {
	out := make([]formField, 0)
	if strings.TrimSpace(payload) == "" {
		return out
	}
	var raw map[string]any
	if err := json.Unmarshal([]byte(payload), &raw); err != nil {
		return out
	}
	arr, _ := raw["fields"].([]any)
	for _, item := range arr {
		switch v := item.(type) {
		case string:
			key := strings.TrimSpace(v)
			if key == "" {
				continue
			}
			out = append(out, formField{Key: key, Label: key, Type: "text"})
		case map[string]any:
			key := strings.TrimSpace(asString(v["key"]))
			if key == "" {
				continue
			}
			label := strings.TrimSpace(asString(v["label"]))
			if label == "" {
				label = key
			}
			typ := strings.TrimSpace(asString(v["type"]))
			if typ == "" {
				typ = "text"
			}
			out = append(out, formField{Key: key, Label: label, Type: typ})
		}
	}
	return out
}

func parseFormValue(raw string) map[string]any {
	out := map[string]any{}
	if strings.TrimSpace(raw) == "" {
		return out
	}
	_ = json.Unmarshal([]byte(raw), &out)
	return out
}

func stringifyFormCell(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		if t {
			return "是"
		}
		return "否"
	default:
		b, err := json.Marshal(t)
		if err != nil {
			return ""
		}
		return string(b)
	}
}

func asString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	default:
		return ""
	}
}

func parseOptionalTime(raw string) (*time.Time, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return nil, nil
	}
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		time.RFC3339,
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			if layout == "2006-01-02" {
				t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
			}
			return &t, nil
		}
	}
	return nil, errors.New("time format")
}

func formatOptTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return formatTime(*t)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(time.Local).Format("2006-01-02 15:04:05")
}

func pageLimit(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

func sanitizeFilePart(name string) string {
	s := strings.TrimSpace(name)
	if s == "" {
		return "activity"
	}
	replacer := strings.NewReplacer("/", "_", "\\", "_", ":", "_", "*", "_", "?", "_", "\"", "_", "<", "_", ">", "_", "|", "_")
	return replacer.Replace(s)
}

func fail(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "操作失败")
}
