package nativeconfigitem

import (
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

// configClassificationHandler owns the platform's configurable categories and
// their ordinary (non-secret) configuration entries.
type configClassificationHandler struct{ adminDB *gorm.DB }

func newConfigClassificationHandler(adminDB *gorm.DB) *configClassificationHandler {
	return &configClassificationHandler{adminDB: adminDB}
}

type configClassificationRow struct {
	ID          uint64    `gorm:"column:id"`
	ParentID    uint64    `gorm:"column:parent_id"`
	Name        string    `gorm:"column:name"`
	ClassifyKey string    `gorm:"column:classify_key"`
	Description string    `gorm:"column:description"`
	Icon        string    `gorm:"column:icon"`
	Status      int       `gorm:"column:status"`
	Sort        int       `gorm:"column:sort"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

type configClassificationItemRow struct {
	ID               uint64    `gorm:"column:id"`
	ClassificationID uint64    `gorm:"column:classification_id"`
	Name             string    `gorm:"column:name"`
	ConfigKey        string    `gorm:"column:config_key"`
	FieldType        string    `gorm:"column:field_type"`
	BackendType      int       `gorm:"column:backend_type"`
	Content          string    `gorm:"column:content"`
	Description      string    `gorm:"column:description"`
	Status           int       `gorm:"column:status"`
	Sort             int       `gorm:"column:sort"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

type configClassificationInput struct {
	ClassifyKey string `json:"classify_key"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
	Sort        *int   `json:"sort"`
	Status      *int   `json:"status"`
}

type configClassificationItemInput struct {
	ClassificationID *uint64 `json:"classification_id"`
	ConfigKey        string  `json:"config_key"`
	FieldType        string  `json:"field_type"`
	BackendType      *int    `json:"backend_type"`
	Content          string  `json:"content"`
	Description      string  `json:"description"`
	Name             string  `json:"name"`
	Sort             *int    `json:"sort"`
	Status           *int    `json:"status"`
}

func (h *configClassificationHandler) list(c *gin.Context) {
	page, limit := pageLimit(c)
	q := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_config_classification").Where("is_del=0")
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		if status != "0" && status != "1" {
			response.Fail(c, http.StatusBadRequest, "显示状态错误")
			return
		}
		q = q.Where("status=?", status)
	}
	if name := strings.TrimSpace(c.Query("name")); name != "" {
		q = q.Where("name LIKE ? OR classify_key LIKE ?", "%"+name+"%", "%"+name+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		configClassificationFail(c)
		return
	}
	rows := make([]configClassificationRow, 0)
	if err := q.Select("id,parent_id,name,classify_key,description,icon,status,sort,created_at,updated_at").
		Order("sort DESC,id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		configClassificationFail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, configClassificationToJSON(row))
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *configClassificationHandler) detail(c *gin.Context) {
	row, ok := h.loadClassification(c, c.Param("id"))
	if ok {
		response.OK(c, configClassificationToJSON(*row))
	}
}

func (h *configClassificationHandler) create(c *gin.Context) {
	var input configClassificationInput
	if c.ShouldBindJSON(&input) != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	name, key, description, icon, status, sort, err := normalizeConfigClassificationInput(input, 1, 0)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	created := struct {
		ID          uint64 `gorm:"column:id;primaryKey"`
		Name        string `gorm:"column:name"`
		ClassifyKey string `gorm:"column:classify_key"`
		Description string `gorm:"column:description"`
		Icon        string `gorm:"column:icon"`
		Status      int    `gorm:"column:status"`
		Sort        int    `gorm:"column:sort"`
		IsDel       int    `gorm:"column:is_del"`
	}{Name: name, ClassifyKey: key, Description: description, Icon: icon, Status: status, Sort: sort}
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_config_classification").Create(&created).Error; err != nil {
		if isDuplicate(err) {
			response.Fail(c, http.StatusConflict, "配置分类 Key 已存在")
			return
		}
		configClassificationFail(c)
		return
	}
	row, ok := h.loadClassification(c, strconv.FormatUint(created.ID, 10))
	if ok {
		response.OK(c, configClassificationToJSON(*row))
	}
}

func (h *configClassificationHandler) update(c *gin.Context) {
	row, ok := h.loadClassification(c, c.Param("id"))
	if !ok {
		return
	}
	var input configClassificationInput
	if c.ShouldBindJSON(&input) != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	name, key, description, icon, status, sort, err := normalizeConfigClassificationInput(input, row.Status, row.Sort)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	result := h.adminDB.WithContext(c.Request.Context()).Exec(`UPDATE qixi_crm_a_config_classification
		SET name=?,classify_key=?,description=?,icon=?,status=?,sort=? WHERE id=? AND is_del=0`, name, key, description, icon, status, sort, row.ID)
	if result.Error != nil {
		if isDuplicate(result.Error) {
			response.Fail(c, http.StatusConflict, "配置分类 Key 已存在")
			return
		}
		configClassificationFail(c)
		return
	}
	updated, ok := h.loadClassification(c, strconv.FormatUint(row.ID, 10))
	if ok {
		response.OK(c, configClassificationToJSON(*updated))
	}
}

func (h *configClassificationHandler) setStatus(c *gin.Context) {
	row, ok := h.loadClassification(c, c.Param("id"))
	if !ok {
		return
	}
	var input struct {
		Status *int `json:"status"`
	}
	if c.ShouldBindJSON(&input) != nil || input.Status == nil || (*input.Status != 0 && *input.Status != 1) {
		response.Fail(c, http.StatusBadRequest, "显示状态错误")
		return
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Exec("UPDATE qixi_crm_a_config_classification SET status=? WHERE id=? AND is_del=0", *input.Status, row.ID).Error; err != nil {
		configClassificationFail(c)
		return
	}
	response.OK(c, gin.H{"id": row.ID, "status": *input.Status})
}

func (h *configClassificationHandler) remove(c *gin.Context) {
	row, ok := h.loadClassification(c, c.Param("id"))
	if !ok {
		return
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE qixi_crm_a_config_classification SET is_del=1,status=0 WHERE id=? AND is_del=0", row.ID).Error; err != nil {
			return err
		}
		return tx.Exec("UPDATE qixi_crm_a_config_classification_item SET is_del=1,status=0 WHERE classification_id=? AND is_del=0", row.ID).Error
	}); err != nil {
		configClassificationFail(c)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *configClassificationHandler) listItems(c *gin.Context) {
	classification, ok := h.loadClassification(c, c.Param("id"))
	if !ok {
		return
	}
	page, limit := pageLimit(c)
	q := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_config_classification_item").Where("classification_id=? AND is_del=0", classification.ID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		configClassificationFail(c)
		return
	}
	rows := make([]configClassificationItemRow, 0)
	if err := q.Select("id,classification_id,name,config_key,field_type,backend_type,content,description,status,sort,created_at,updated_at").Order("sort DESC,id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		configClassificationFail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, configClassificationItemToJSON(row))
	}
	response.OK(c, gin.H{"classification": configClassificationToJSON(*classification), "list": list, "total": total, "page": page, "limit": limit})
}

func (h *configClassificationHandler) createItem(c *gin.Context) {
	classification, ok := h.loadClassification(c, c.Param("id"))
	if !ok {
		return
	}
	var input configClassificationItemInput
	if c.ShouldBindJSON(&input) != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	name, key, fieldType, backendType, content, description, status, sort, err := normalizeConfigClassificationItemInput(input, "input", 0, 1, 0)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	created := struct {
		ID               uint64 `gorm:"column:id;primaryKey"`
		ClassificationID uint64 `gorm:"column:classification_id"`
		Name             string `gorm:"column:name"`
		ConfigKey        string `gorm:"column:config_key"`
		FieldType        string `gorm:"column:field_type"`
		BackendType      int    `gorm:"column:backend_type"`
		Content          string `gorm:"column:content"`
		Description      string `gorm:"column:description"`
		Status           int    `gorm:"column:status"`
		Sort             int    `gorm:"column:sort"`
		IsDel            int    `gorm:"column:is_del"`
	}{ClassificationID: classification.ID, Name: name, ConfigKey: key, FieldType: fieldType, BackendType: backendType, Content: content, Description: description, Status: status, Sort: sort}
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_config_classification_item").Create(&created).Error; err != nil {
		if isDuplicate(err) {
			response.Fail(c, http.StatusConflict, "该分类下配置 Key 已存在")
			return
		}
		configClassificationFail(c)
		return
	}
	item, ok := h.loadItem(c, classification.ID, strconv.FormatUint(created.ID, 10))
	if ok {
		response.OK(c, configClassificationItemToJSON(*item))
	}
}

func (h *configClassificationHandler) updateItem(c *gin.Context) {
	classification, ok := h.loadClassification(c, c.Param("id"))
	if !ok {
		return
	}
	item, ok := h.loadItem(c, classification.ID, c.Param("itemID"))
	if !ok {
		return
	}
	var input configClassificationItemInput
	if c.ShouldBindJSON(&input) != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	name, key, fieldType, backendType, content, description, status, sort, err := normalizeConfigClassificationItemInput(input, item.FieldType, item.BackendType, item.Status, item.Sort)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	result := h.adminDB.WithContext(c.Request.Context()).Exec(`UPDATE qixi_crm_a_config_classification_item
		SET name=?,config_key=?,field_type=?,backend_type=?,content=?,description=?,status=?,sort=? WHERE id=? AND classification_id=? AND is_del=0`, name, key, fieldType, backendType, content, description, status, sort, item.ID, classification.ID)
	if result.Error != nil {
		if isDuplicate(result.Error) {
			response.Fail(c, http.StatusConflict, "该分类下配置 Key 已存在")
			return
		}
		configClassificationFail(c)
		return
	}
	updated, ok := h.loadItem(c, classification.ID, strconv.FormatUint(item.ID, 10))
	if ok {
		response.OK(c, configClassificationItemToJSON(*updated))
	}
}

func (h *configClassificationHandler) setItemStatus(c *gin.Context) {
	classification, ok := h.loadClassification(c, c.Param("id"))
	if !ok {
		return
	}
	item, ok := h.loadItem(c, classification.ID, c.Param("itemID"))
	if !ok {
		return
	}
	var input struct {
		Status *int `json:"status"`
	}
	if c.ShouldBindJSON(&input) != nil || input.Status == nil || (*input.Status != 0 && *input.Status != 1) {
		response.Fail(c, http.StatusBadRequest, "显示状态错误")
		return
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Exec("UPDATE qixi_crm_a_config_classification_item SET status=? WHERE id=? AND classification_id=? AND is_del=0", *input.Status, item.ID, classification.ID).Error; err != nil {
		configClassificationFail(c)
		return
	}
	response.OK(c, gin.H{"id": item.ID, "status": *input.Status})
}

func (h *configClassificationHandler) removeItem(c *gin.Context) {
	classification, ok := h.loadClassification(c, c.Param("id"))
	if !ok {
		return
	}
	item, ok := h.loadItem(c, classification.ID, c.Param("itemID"))
	if !ok {
		return
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Exec("UPDATE qixi_crm_a_config_classification_item SET is_del=1,status=0 WHERE id=? AND classification_id=? AND is_del=0", item.ID, classification.ID).Error; err != nil {
		configClassificationFail(c)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *configClassificationHandler) listSettings(c *gin.Context) {
	page, limit := pageLimit(c)
	q := h.adminDB.WithContext(c.Request.Context()).
		Table("qixi_crm_a_config_classification_item AS item").
		Joins("JOIN qixi_crm_a_config_classification AS classification ON classification.id=item.classification_id AND classification.is_del=0").
		Where("item.is_del=0")
	if backendType := strings.TrimSpace(c.Query("backend_type")); backendType != "" {
		if backendType != "0" && backendType != "1" {
			response.Fail(c, http.StatusBadRequest, "后台类型错误")
			return
		}
		q = q.Where("item.backend_type=?", backendType)
	}
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		q = q.Where("item.name LIKE ? OR item.config_key LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		configClassificationFail(c)
		return
	}
	type settingListRow struct {
		configClassificationItemRow
		ClassificationName string `gorm:"column:classification_name"`
	}
	rows := make([]settingListRow, 0)
	if err := q.Select(`item.id,item.classification_id,item.name,item.config_key,item.field_type,item.backend_type,
		item.content,item.description,item.status,item.sort,item.created_at,item.updated_at,
		classification.name AS classification_name`).
		// 同一排序值保持配置初始化与创建的自然顺序，避免默认配置在列表中反向展示。
		Order("item.sort DESC,item.id ASC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		configClassificationFail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		data := configClassificationItemToJSON(row.configClassificationItemRow)
		data["classification_name"] = row.ClassificationName
		list = append(list, data)
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *configClassificationHandler) createSetting(c *gin.Context) {
	var input configClassificationItemInput
	if c.ShouldBindJSON(&input) != nil || input.ClassificationID == nil || *input.ClassificationID == 0 {
		response.Fail(c, http.StatusBadRequest, "请选择配置分类")
		return
	}
	classification, ok := h.loadClassification(c, strconv.FormatUint(*input.ClassificationID, 10))
	if !ok {
		return
	}
	name, key, fieldType, backendType, content, description, status, sort, err := normalizeConfigClassificationItemInput(input, "input", 0, 1, 0)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	created := struct {
		ID               uint64 `gorm:"column:id;primaryKey"`
		ClassificationID uint64 `gorm:"column:classification_id"`
		Name             string `gorm:"column:name"`
		ConfigKey        string `gorm:"column:config_key"`
		FieldType        string `gorm:"column:field_type"`
		BackendType      int    `gorm:"column:backend_type"`
		Content          string `gorm:"column:content"`
		Description      string `gorm:"column:description"`
		Status           int    `gorm:"column:status"`
		Sort             int    `gorm:"column:sort"`
		IsDel            int    `gorm:"column:is_del"`
	}{ClassificationID: classification.ID, Name: name, ConfigKey: key, FieldType: fieldType, BackendType: backendType, Content: content, Description: description, Status: status, Sort: sort}
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_config_classification_item").Create(&created).Error; err != nil {
		if isDuplicate(err) {
			response.Fail(c, http.StatusConflict, "该分类下配置 Key 已存在")
			return
		}
		configClassificationFail(c)
		return
	}
	createdRow, ok := h.loadItem(c, classification.ID, strconv.FormatUint(created.ID, 10))
	if ok {
		response.OK(c, configClassificationItemToJSON(*createdRow))
	}
}

func (h *configClassificationHandler) updateSetting(c *gin.Context) {
	id, err := parsePositiveID(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "配置 ID 错误")
		return
	}
	var existing configClassificationItemRow
	err = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_config_classification_item").
		Select("id,classification_id,name,config_key,field_type,backend_type,content,description,status,sort,created_at,updated_at").
		Where("id=? AND is_del=0", id).Take(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "配置不存在")
		return
	}
	if err != nil {
		configClassificationFail(c)
		return
	}
	var input configClassificationItemInput
	if c.ShouldBindJSON(&input) != nil || input.ClassificationID == nil || *input.ClassificationID == 0 {
		response.Fail(c, http.StatusBadRequest, "请选择配置分类")
		return
	}
	classification, ok := h.loadClassification(c, strconv.FormatUint(*input.ClassificationID, 10))
	if !ok {
		return
	}
	name, key, fieldType, backendType, content, description, status, sort, err := normalizeConfigClassificationItemInput(input, existing.FieldType, existing.BackendType, existing.Status, existing.Sort)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	result := h.adminDB.WithContext(c.Request.Context()).Exec(`UPDATE qixi_crm_a_config_classification_item
		SET classification_id=?,name=?,config_key=?,field_type=?,backend_type=?,content=?,description=?,status=?,sort=?
		WHERE id=? AND is_del=0`, classification.ID, name, key, fieldType, backendType, content, description, status, sort, existing.ID)
	if result.Error != nil {
		if isDuplicate(result.Error) {
			response.Fail(c, http.StatusConflict, "该分类下配置 Key 已存在")
			return
		}
		configClassificationFail(c)
		return
	}
	updated, err := h.loadSetting(c, existing.ID)
	if err != nil {
		return
	}
	response.OK(c, configClassificationItemToJSON(*updated))
}

func (h *configClassificationHandler) setSettingStatus(c *gin.Context) {
	id, err := parsePositiveID(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "配置 ID 错误")
		return
	}
	var input struct {
		Status *int `json:"status"`
	}
	if c.ShouldBindJSON(&input) != nil || input.Status == nil || (*input.Status != 0 && *input.Status != 1) {
		response.Fail(c, http.StatusBadRequest, "显示状态错误")
		return
	}
	result := h.adminDB.WithContext(c.Request.Context()).Exec("UPDATE qixi_crm_a_config_classification_item SET status=? WHERE id=? AND is_del=0", *input.Status, id)
	if result.Error != nil {
		configClassificationFail(c)
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "配置不存在")
		return
	}
	response.OK(c, gin.H{"id": id, "status": *input.Status})
}

func (h *configClassificationHandler) removeSetting(c *gin.Context) {
	id, err := parsePositiveID(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "配置 ID 错误")
		return
	}
	result := h.adminDB.WithContext(c.Request.Context()).Exec("UPDATE qixi_crm_a_config_classification_item SET is_del=1,status=0 WHERE id=? AND is_del=0", id)
	if result.Error != nil {
		configClassificationFail(c)
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "配置不存在")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *configClassificationHandler) loadSetting(c *gin.Context, id uint64) (*configClassificationItemRow, error) {
	var row configClassificationItemRow
	err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_config_classification_item").
		Select("id,classification_id,name,config_key,field_type,backend_type,content,description,status,sort,created_at,updated_at").
		Where("id=? AND is_del=0", id).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "配置不存在")
		return nil, err
	}
	if err != nil {
		configClassificationFail(c)
		return nil, err
	}
	return &row, nil
}

func (h *configClassificationHandler) loadClassification(c *gin.Context, rawID string) (*configClassificationRow, bool) {
	id, err := parsePositiveID(rawID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "配置分类 ID 错误")
		return nil, false
	}
	var row configClassificationRow
	err = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_config_classification").Select("id,parent_id,name,classify_key,description,icon,status,sort,created_at,updated_at").Where("id=? AND is_del=0", id).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "配置分类不存在")
		return nil, false
	}
	if err != nil {
		configClassificationFail(c)
		return nil, false
	}
	return &row, true
}

func (h *configClassificationHandler) loadItem(c *gin.Context, classificationID uint64, rawID string) (*configClassificationItemRow, bool) {
	id, err := parsePositiveID(rawID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "配置项 ID 错误")
		return nil, false
	}
	var row configClassificationItemRow
	err = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_config_classification_item").Select("id,classification_id,name,config_key,field_type,backend_type,content,description,status,sort,created_at,updated_at").Where("id=? AND classification_id=? AND is_del=0", id, classificationID).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "配置项不存在")
		return nil, false
	}
	if err != nil {
		configClassificationFail(c)
		return nil, false
	}
	return &row, true
}

func normalizeConfigClassificationInput(in configClassificationInput, defaultStatus, defaultSort int) (string, string, string, string, int, int, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || utf8.RuneCountInString(name) > 128 {
		return "", "", "", "", 0, 0, errors.New("配置分类名称必填且不超过 128 字")
	}
	key := strings.TrimSpace(in.ClassifyKey)
	if !validConfigKey(key, 64) {
		return "", "", "", "", 0, 0, errors.New("配置分类 Key 仅支持字母、数字、点、下划线和连字符")
	}
	description := strings.TrimSpace(in.Description)
	icon := strings.TrimSpace(in.Icon)
	if utf8.RuneCountInString(description) > 500 || utf8.RuneCountInString(icon) > 96 {
		return "", "", "", "", 0, 0, errors.New("配置分类说明或图标过长")
	}
	status, sort, err := normalizeConfigStatusAndSort(in.Status, in.Sort, defaultStatus, defaultSort)
	return name, key, description, icon, status, sort, err
}

func normalizeConfigClassificationItemInput(in configClassificationItemInput, defaultFieldType string, defaultBackendType, defaultStatus, defaultSort int) (string, string, string, int, string, string, int, int, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || utf8.RuneCountInString(name) > 128 {
		return "", "", "", 0, "", "", 0, 0, errors.New("配置名称必填且不超过 128 字")
	}
	key := strings.TrimSpace(in.ConfigKey)
	if !validConfigKey(key, 128) {
		return "", "", "", 0, "", "", 0, 0, errors.New("配置 Key 仅支持字母、数字、点、下划线和连字符")
	}
	fieldType := strings.TrimSpace(in.FieldType)
	if fieldType == "" {
		fieldType = defaultFieldType
	}
	if !validConfigFieldType(fieldType) {
		return "", "", "", 0, "", "", 0, 0, errors.New("配置类型错误")
	}
	backendType := defaultBackendType
	if in.BackendType != nil {
		backendType = *in.BackendType
	}
	if backendType != 0 && backendType != 1 {
		return "", "", "", 0, "", "", 0, 0, errors.New("后台类型错误")
	}
	content := strings.TrimSpace(in.Content)
	if content == "" || utf8.RuneCountInString(content) > 10000 {
		return "", "", "", 0, "", "", 0, 0, errors.New("配置内容必填且不超过 10000 字")
	}
	description := strings.TrimSpace(in.Description)
	if utf8.RuneCountInString(description) > 500 {
		return "", "", "", 0, "", "", 0, 0, errors.New("配置说明不超过 500 字")
	}
	status, sort, err := normalizeConfigStatusAndSort(in.Status, in.Sort, defaultStatus, defaultSort)
	return name, key, fieldType, backendType, content, description, status, sort, err
}

func validConfigFieldType(fieldType string) bool {
	switch fieldType {
	case "input", "textarea", "number", "radio", "switch", "image", "file":
		return true
	default:
		return false
	}
}

func normalizeConfigStatusAndSort(statusInput, sortInput *int, defaultStatus, defaultSort int) (int, int, error) {
	status, sort := defaultStatus, defaultSort
	if statusInput != nil {
		if *statusInput != 0 && *statusInput != 1 {
			return 0, 0, errors.New("显示状态错误")
		}
		status = *statusInput
	}
	if sortInput != nil {
		sort = *sortInput
	}
	return status, sort, nil
}

func validConfigKey(key string, max int) bool {
	if key == "" || utf8.RuneCountInString(key) > max {
		return false
	}
	for _, r := range key {
		if !(r == '_' || r == '-' || r == '.' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z') {
			return false
		}
	}
	return true
}

func configClassificationToJSON(row configClassificationRow) gin.H {
	return gin.H{"id": row.ID, "parent_id": row.ParentID, "name": row.Name, "classify_key": row.ClassifyKey, "description": row.Description, "icon": row.Icon, "status": row.Status, "sort": row.Sort, "created_at": row.CreatedAt.Format("2006-01-02 15:04:05"), "updated_at": row.UpdatedAt.Format("2006-01-02 15:04:05")}
}

func configClassificationItemToJSON(row configClassificationItemRow) gin.H {
	return gin.H{"id": row.ID, "classification_id": row.ClassificationID, "name": row.Name, "config_key": row.ConfigKey, "field_type": row.FieldType, "backend_type": row.BackendType, "content": row.Content, "description": row.Description, "status": row.Status, "sort": row.Sort, "created_at": row.CreatedAt.Format("2006-01-02 15:04:05"), "updated_at": row.UpdatedAt.Format("2006-01-02 15:04:05")}
}

func configClassificationFail(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "配置分类操作失败")
}
