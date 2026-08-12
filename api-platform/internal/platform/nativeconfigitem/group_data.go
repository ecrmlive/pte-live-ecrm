package nativeconfigitem

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// groupDataHandler keeps the CRMEB-compatible group/data hierarchy separate
// from ordinary platform configuration items.
type groupDataHandler struct{ adminDB *gorm.DB }

func newGroupDataHandler(adminDB *gorm.DB) *groupDataHandler {
	return &groupDataHandler{adminDB: adminDB}
}

type dataGroupRow struct {
	ID          uint64    `gorm:"column:id"`
	Name        string    `gorm:"column:name"`
	GroupKey    string    `gorm:"column:group_key"`
	Description string    `gorm:"column:description"`
	Fields      string    `gorm:"column:fields"`
	Sort        int       `gorm:"column:sort"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

type dataGroupItemRow struct {
	ID        uint64    `gorm:"column:id"`
	GroupID   uint64    `gorm:"column:group_id"`
	Data      string    `gorm:"column:data"`
	Sort      int       `gorm:"column:sort"`
	Status    int       `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type dataGroupInput struct {
	Name        string           `json:"name"`
	GroupKey    string           `json:"group_key"`
	Description string           `json:"description"`
	Fields      []map[string]any `json:"fields"`
	Sort        *int             `json:"sort"`
}

type dataGroupItemInput struct {
	Data   map[string]any `json:"data"`
	Sort   *int           `json:"sort"`
	Status *int           `json:"status"`
}

func (h *groupDataHandler) list(c *gin.Context) {
	page, limit := pageLimit(c)
	q := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_data_group").Where("is_del = 0")
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		q = q.Where("name LIKE ? OR group_key LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		groupFail(c)
		return
	}
	rows := make([]dataGroupRow, 0)
	if err := q.Select("id,name,group_key,description,fields,sort,created_at,updated_at").
		Order("sort DESC, id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		groupFail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, groupToJSON(row))
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *groupDataHandler) detail(c *gin.Context) {
	row, ok := h.loadGroup(c, c.Param("id"))
	if ok {
		response.OK(c, groupToJSON(*row))
	}
}

func (h *groupDataHandler) create(c *gin.Context) {
	var in dataGroupInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	name, key, description, fields, sort, err := normalizeGroupInput(in, 0)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	created := struct {
		ID          uint64 `gorm:"column:id;primaryKey"`
		Name        string `gorm:"column:name"`
		GroupKey    string `gorm:"column:group_key"`
		Description string `gorm:"column:description"`
		Fields      string `gorm:"column:fields"`
		Sort        int    `gorm:"column:sort"`
		IsDel       int    `gorm:"column:is_del"`
	}{Name: name, GroupKey: key, Description: description, Fields: fields, Sort: sort}
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_data_group").Create(&created).Error; err != nil {
		if isDuplicate(err) {
			response.Fail(c, http.StatusConflict, "数据组 Key 已存在")
			return
		}
		groupFail(c)
		return
	}
	row, ok := h.loadGroup(c, strconv.FormatUint(created.ID, 10))
	if ok {
		response.OK(c, groupToJSON(*row))
	}
}

func (h *groupDataHandler) update(c *gin.Context) {
	row, ok := h.loadGroup(c, c.Param("id"))
	if !ok {
		return
	}
	var in dataGroupInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	name, key, description, fields, sort, err := normalizeGroupInput(in, row.Sort)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	result := h.adminDB.WithContext(c.Request.Context()).Exec(`
		UPDATE qixi_crm_a_data_group
		SET name=?, group_key=?, description=?, fields=?, sort=?
		WHERE id=? AND is_del=0`, name, key, description, fields, sort, row.ID)
	if result.Error != nil {
		if isDuplicate(result.Error) {
			response.Fail(c, http.StatusConflict, "数据组 Key 已存在")
			return
		}
		groupFail(c)
		return
	}
	updated, ok := h.loadGroup(c, strconv.FormatUint(row.ID, 10))
	if ok {
		response.OK(c, groupToJSON(*updated))
	}
}

func (h *groupDataHandler) remove(c *gin.Context) {
	row, ok := h.loadGroup(c, c.Param("id"))
	if !ok {
		return
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE qixi_crm_a_data_group SET is_del=1 WHERE id=? AND is_del=0", row.ID).Error; err != nil {
			return err
		}
		return tx.Exec("UPDATE qixi_crm_a_data_group_item SET is_del=1, status=0 WHERE group_id=? AND is_del=0", row.ID).Error
	}); err != nil {
		groupFail(c)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *groupDataHandler) listItems(c *gin.Context) {
	group, ok := h.loadGroup(c, c.Param("id"))
	if !ok {
		return
	}
	page, limit := pageLimit(c)
	q := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_data_group_item").Where("group_id=? AND is_del=0", group.ID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		groupFail(c)
		return
	}
	rows := make([]dataGroupItemRow, 0)
	if err := q.Select("id,group_id,data,sort,status,created_at,updated_at").Order("sort DESC, id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		groupFail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, groupItemToJSON(row))
	}
	response.OK(c, gin.H{"group": groupToJSON(*group), "list": list, "total": total, "page": page, "limit": limit})
}

func (h *groupDataHandler) createItem(c *gin.Context) {
	group, ok := h.loadGroup(c, c.Param("id"))
	if !ok {
		return
	}
	var in dataGroupItemInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	data, sort, status, err := normalizeGroupItemInput(in, 0, 1)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	created := struct {
		ID      uint64 `gorm:"column:id;primaryKey"`
		GroupID uint64 `gorm:"column:group_id"`
		Data    string `gorm:"column:data"`
		Sort    int    `gorm:"column:sort"`
		Status  int    `gorm:"column:status"`
		IsDel   int    `gorm:"column:is_del"`
	}{GroupID: group.ID, Data: data, Sort: sort, Status: status}
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_data_group_item").Create(&created).Error; err != nil {
		groupFail(c)
		return
	}
	item, ok := h.loadItem(c, group.ID, strconv.FormatUint(created.ID, 10))
	if ok {
		response.OK(c, groupItemToJSON(*item))
	}
}

func (h *groupDataHandler) updateItem(c *gin.Context) {
	group, ok := h.loadGroup(c, c.Param("id"))
	if !ok {
		return
	}
	item, ok := h.loadItem(c, group.ID, c.Param("itemID"))
	if !ok {
		return
	}
	var in dataGroupItemInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	data, sort, status, err := normalizeGroupItemInput(in, item.Sort, item.Status)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Exec("UPDATE qixi_crm_a_data_group_item SET data=?,sort=?,status=? WHERE id=? AND group_id=? AND is_del=0", data, sort, status, item.ID, group.ID).Error; err != nil {
		groupFail(c)
		return
	}
	updated, ok := h.loadItem(c, group.ID, strconv.FormatUint(item.ID, 10))
	if ok {
		response.OK(c, groupItemToJSON(*updated))
	}
}

func (h *groupDataHandler) setItemStatus(c *gin.Context) {
	group, ok := h.loadGroup(c, c.Param("id"))
	if !ok {
		return
	}
	item, ok := h.loadItem(c, group.ID, c.Param("itemID"))
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
	if err := h.adminDB.WithContext(c.Request.Context()).Exec("UPDATE qixi_crm_a_data_group_item SET status=? WHERE id=? AND group_id=? AND is_del=0", *in.Status, item.ID, group.ID).Error; err != nil {
		groupFail(c)
		return
	}
	response.OK(c, gin.H{"id": item.ID, "status": *in.Status})
}

func (h *groupDataHandler) removeItem(c *gin.Context) {
	group, ok := h.loadGroup(c, c.Param("id"))
	if !ok {
		return
	}
	item, ok := h.loadItem(c, group.ID, c.Param("itemID"))
	if !ok {
		return
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Exec("UPDATE qixi_crm_a_data_group_item SET is_del=1,status=0 WHERE id=? AND group_id=? AND is_del=0", item.ID, group.ID).Error; err != nil {
		groupFail(c)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *groupDataHandler) loadGroup(c *gin.Context, rawID string) (*dataGroupRow, bool) {
	id, err := parsePositiveID(rawID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "数据组 ID 错误")
		return nil, false
	}
	var row dataGroupRow
	err = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_data_group").Select("id,name,group_key,description,fields,sort,created_at,updated_at").Where("id=? AND is_del=0", id).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "组合数据不存在")
		return nil, false
	}
	if err != nil {
		groupFail(c)
		return nil, false
	}
	return &row, true
}

func (h *groupDataHandler) loadItem(c *gin.Context, groupID uint64, rawID string) (*dataGroupItemRow, bool) {
	id, err := parsePositiveID(rawID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "数据 ID 错误")
		return nil, false
	}
	var row dataGroupItemRow
	err = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_data_group_item").Select("id,group_id,data,sort,status,created_at,updated_at").Where("id=? AND group_id=? AND is_del=0", id, groupID).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "组合数据项不存在")
		return nil, false
	}
	if err != nil {
		groupFail(c)
		return nil, false
	}
	return &row, true
}

func normalizeGroupInput(in dataGroupInput, defaultSort int) (string, string, string, string, int, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || utf8.RuneCountInString(name) > 128 {
		return "", "", "", "", 0, errors.New("数据组名称必填且不超过 128 字")
	}
	key := strings.TrimSpace(in.GroupKey)
	if key == "" || utf8.RuneCountInString(key) > 64 {
		return "", "", "", "", 0, errors.New("数据组 Key 必填且不超过 64 字")
	}
	for _, r := range key {
		if !(r == '_' || r == '-' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z') {
			return "", "", "", "", 0, errors.New("数据组 Key 仅支持字母、数字、下划线和连字符")
		}
	}
	description := strings.TrimSpace(in.Description)
	if utf8.RuneCountInString(description) > 500 {
		return "", "", "", "", 0, errors.New("数据组说明不超过 500 字")
	}
	fields := in.Fields
	if fields == nil {
		fields = []map[string]any{}
	}
	raw, err := json.Marshal(fields)
	if err != nil {
		return "", "", "", "", 0, errors.New("数据字段编码失败")
	}
	sort := defaultSort
	if in.Sort != nil {
		sort = *in.Sort
	}
	return name, key, description, string(raw), sort, nil
}

func normalizeGroupItemInput(in dataGroupItemInput, defaultSort, defaultStatus int) (string, int, int, error) {
	if in.Data == nil {
		return "", 0, 0, errors.New("请填写数据内容")
	}
	raw, err := json.Marshal(in.Data)
	if err != nil || len(raw) > 65535 {
		return "", 0, 0, errors.New("数据内容格式错误或过长")
	}
	sort := defaultSort
	if in.Sort != nil {
		sort = *in.Sort
	}
	status := defaultStatus
	if in.Status != nil {
		if *in.Status != 0 && *in.Status != 1 {
			return "", 0, 0, errors.New("状态错误")
		}
		status = *in.Status
	}
	return string(raw), sort, status, nil
}

func groupToJSON(row dataGroupRow) gin.H {
	fields := make([]map[string]any, 0)
	_ = json.Unmarshal([]byte(row.Fields), &fields)
	return gin.H{"id": row.ID, "name": row.Name, "group_key": row.GroupKey, "description": row.Description, "fields": fields, "sort": row.Sort, "created_at": row.CreatedAt.Format("2006-01-02 15:04:05"), "updated_at": row.UpdatedAt.Format("2006-01-02 15:04:05")}
}

func groupItemToJSON(row dataGroupItemRow) gin.H {
	data := map[string]any{}
	_ = json.Unmarshal([]byte(row.Data), &data)
	return gin.H{"id": row.ID, "group_id": row.GroupID, "data": data, "sort": row.Sort, "status": row.Status, "created_at": row.CreatedAt.Format("2006-01-02 15:04:05"), "updated_at": row.UpdatedAt.Format("2006-01-02 15:04:05")}
}

func parsePositiveID(raw string) (uint64, error) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func isDuplicate(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "duplicate")
}

func groupFail(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "组合数据操作失败")
}
