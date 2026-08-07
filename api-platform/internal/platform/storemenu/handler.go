package storemenu

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// Handler 平台侧「店铺菜单」：读写商户库 qixi_crm_m_menu 目录（供店铺类型授权勾选的菜单树基线）。
type Handler struct{ merchantDB *gorm.DB }

func NewHandler(merchantDB *gorm.DB) *Handler { return &Handler{merchantDB: merchantDB} }

func (h *Handler) Register(r gin.IRoutes) {
	p := middleware.RequireAdminRoles("platform")
	r.GET("/merchant-menus", p, h.List)
	r.POST("/merchant-menus", p, h.Create)
	r.PUT("/merchant-menus/:id", p, h.Update)
	r.DELETE("/merchant-menus/:id", p, h.Delete)
}

type menuRow struct {
	MenuID    uint64 `json:"menu_id"`
	PID       uint64 `json:"pid"`
	Code      string `json:"code"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Icon      string `json:"icon"`
	MenuName  string `json:"menu_name"`
	Route     string `json:"route"`
	Sort      int    `json:"sort"`
	IsShow    int8   `json:"is_show"`
	IsMenu    int8   `json:"is_menu"`
	IsRoute   int8   `json:"is_route"`
	IsAgent   int8   `json:"is_agent"`
	CreatedAt string `json:"created_at"`
}

type menuInput struct {
	ParentID  *uint64 `json:"parent_id"`
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Path      string  `json:"path"`
	Component string  `json:"component"`
	Icon      string  `json:"icon"`
	IsMenu    *int8   `json:"is_menu"`
	IsRoute   *int8   `json:"is_route"`
	Sort      *int    `json:"sort"`
	Status    *int8   `json:"status"`
}

type dbMenuRow struct {
	ID        uint64    `gorm:"column:id"`
	ParentID  uint64    `gorm:"column:parent_id"`
	Code      string    `gorm:"column:code"`
	Name      string    `gorm:"column:name"`
	Path      string    `gorm:"column:path"`
	Component string    `gorm:"column:component"`
	Icon      string    `gorm:"column:icon"`
	IsMenu    int8      `gorm:"column:is_menu"`
	IsRoute   int8      `gorm:"column:is_route"`
	Sort      int       `gorm:"column:sort"`
	Status    int8      `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (h *Handler) List(c *gin.Context) {
	var rows []dbMenuRow
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_menu").
		Select("id,parent_id,code,name,path,component,icon,is_menu,is_route,sort,status,created_at").
		Order("sort DESC, id ASC").Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询店铺菜单失败")
		return
	}
	out := make([]menuRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, toMenuRow(row))
	}
	response.OK(c, gin.H{"list": out, "total": len(out)})
}

func (h *Handler) Create(c *gin.Context) {
	var req menuInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, errMsg := normalizeMenuInput(&req, true)
	if errMsg != "" {
		response.Fail(c, http.StatusBadRequest, errMsg)
		return
	}
	parentID := uint64(0)
	if req.ParentID != nil {
		parentID = *req.ParentID
	}
	if parentID > 0 {
		var parent dbMenuRow
		if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_menu").
			Select("id").Where("id = ?", parentID).Take(&parent).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.Fail(c, http.StatusBadRequest, "父级菜单不存在")
			} else {
				response.Fail(c, http.StatusInternalServerError, "校验父级菜单失败")
			}
			return
		}
	}
	insert := map[string]any{
		"parent_id": parentID,
		"code":      row.Code,
		"name":      row.Name,
		"path":      row.Path,
		"component": row.Component,
		"icon":      row.Icon,
		"is_menu":   row.IsMenu,
		"is_route":  row.IsRoute,
		"sort":      row.Sort,
		"status":    row.Status,
	}
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_menu").Create(insert).Error; err != nil {
		if isDuplicateKey(err) {
			response.Fail(c, http.StatusConflict, "菜单权限码已存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "创建店铺菜单失败")
		return
	}
	var created dbMenuRow
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_menu").
		Select("id,parent_id,code,name,path,component,icon,is_menu,is_route,sort,status,created_at").
		Where("code = ?", row.Code).Take(&created).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询店铺菜单失败")
		return
	}
	response.OK(c, toMenuRow(created))
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "菜单 ID 错误")
		return
	}
	var req menuInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, errMsg := normalizeMenuInput(&req, false)
	if errMsg != "" {
		response.Fail(c, http.StatusBadRequest, errMsg)
		return
	}
	updates := map[string]any{
		"code":      row.Code,
		"name":      row.Name,
		"path":      row.Path,
		"component": row.Component,
		"icon":      row.Icon,
		"is_menu":   row.IsMenu,
		"is_route":  row.IsRoute,
		"sort":      row.Sort,
		"status":    row.Status,
	}
	if req.ParentID != nil {
		parentID := *req.ParentID
		if parentID == id {
			response.Fail(c, http.StatusBadRequest, "父级菜单不能为自身")
			return
		}
		if parentID > 0 {
			var parent dbMenuRow
			if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_menu").
				Select("id").Where("id = ?", parentID).Take(&parent).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					response.Fail(c, http.StatusBadRequest, "父级菜单不存在")
				} else {
					response.Fail(c, http.StatusInternalServerError, "校验父级菜单失败")
				}
				return
			}
			if descendant, derr := h.isDescendant(c, id, parentID); derr != nil {
				response.Fail(c, http.StatusInternalServerError, "校验父级菜单失败")
				return
			} else if descendant {
				response.Fail(c, http.StatusBadRequest, "父级菜单不能为当前菜单的子级")
				return
			}
		}
		updates["parent_id"] = parentID
	}
	res := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_menu").Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		if isDuplicateKey(res.Error) {
			response.Fail(c, http.StatusConflict, "菜单权限码已存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "更新店铺菜单失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "菜单不存在")
		return
	}
	updated, err := h.loadMenuByID(c, id)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询店铺菜单失败")
		return
	}
	response.OK(c, updated)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "菜单 ID 错误")
		return
	}
	var childCount int64
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_menu").
		Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "校验子菜单失败")
		return
	}
	if childCount > 0 {
		response.Fail(c, http.StatusConflict, "请先删除子菜单")
		return
	}
	tx := h.merchantDB.WithContext(c.Request.Context()).Begin()
	if tx.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "删除店铺菜单失败")
		return
	}
	if err := tx.Table("qixi_crm_m_role_menu").Where("menu_id = ?", id).Delete(nil).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "清理角色菜单关联失败")
		return
	}
	res := tx.Table("qixi_crm_m_menu").Where("id = ?", id).Delete(nil)
	if res.Error != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "删除店铺菜单失败")
		return
	}
	if res.RowsAffected == 0 {
		tx.Rollback()
		response.Fail(c, http.StatusNotFound, "菜单不存在")
		return
	}
	if err := tx.Commit().Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除店铺菜单失败")
		return
	}
	response.OK(c, gin.H{"menu_id": id})
}

func (h *Handler) loadMenuByID(c *gin.Context, id uint64) (menuRow, error) {
	var row dbMenuRow
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_menu").
		Select("id,parent_id,code,name,path,component,icon,is_menu,is_route,sort,status,created_at").
		Where("id = ?", id).Take(&row).Error; err != nil {
		return menuRow{}, err
	}
	return toMenuRow(row), nil
}

func (h *Handler) isDescendant(c *gin.Context, rootID, targetID uint64) (bool, error) {
	current := targetID
	for current > 0 {
		if current == rootID {
			return true, nil
		}
		var row struct {
			ParentID uint64 `gorm:"column:parent_id"`
		}
		if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_menu").
			Select("parent_id").Where("id = ?", current).Take(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}
		current = row.ParentID
	}
	return false, nil
}

func normalizeMenuInput(req *menuInput, create bool) (dbMenuRow, string) {
	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	path := strings.TrimSpace(req.Path)
	component := strings.TrimSpace(req.Component)
	icon := strings.TrimSpace(req.Icon)
	if create && code == "" {
		return dbMenuRow{}, "菜单权限码不能为空"
	}
	if create && name == "" {
		return dbMenuRow{}, "菜单名称不能为空"
	}
	if name != "" {
		if utf8.RuneCountInString(name) > 64 {
			return dbMenuRow{}, "菜单名称不能超过 64 字"
		}
	}
	if code != "" {
		if len(code) > 128 {
			return dbMenuRow{}, "菜单权限码不能超过 128 字符"
		}
	}
	isMenu := int8(1)
	if req.IsMenu != nil {
		if *req.IsMenu != 1 && *req.IsMenu != 2 {
			return dbMenuRow{}, "菜单类型仅支持 1（页面）或 2（按钮）"
		}
		isMenu = *req.IsMenu
	}
	isRoute := int8(1)
	if req.IsRoute != nil {
		if *req.IsRoute != 0 && *req.IsRoute != 1 {
			return dbMenuRow{}, "路由类型仅支持 0/1"
		}
		isRoute = *req.IsRoute
	}
	sort := 0
	if req.Sort != nil {
		sort = *req.Sort
	}
	status := int8(1)
	if req.Status != nil {
		if *req.Status != 0 && *req.Status != 1 {
			return dbMenuRow{}, "显示状态仅支持 0/1"
		}
		status = *req.Status
	}
	return dbMenuRow{
		Code:      code,
		Name:      name,
		Path:      path,
		Component: component,
		Icon:      icon,
		IsMenu:    isMenu,
		IsRoute:   isRoute,
		Sort:      sort,
		Status:    status,
	}, ""
}

func toMenuRow(row dbMenuRow) menuRow {
	route := strings.TrimSpace(row.Component)
	if route == "" {
		route = row.Code
	}
	createdAt := ""
	if !row.CreatedAt.IsZero() {
		createdAt = row.CreatedAt.Format("2006-01-02 15:04:05")
	}
	return menuRow{
		MenuID:    row.ID,
		PID:       row.ParentID,
		Code:      row.Code,
		Path:      row.Path,
		Component: row.Component,
		Icon:      row.Icon,
		MenuName:  row.Name,
		Route:     route,
		Sort:      row.Sort,
		IsShow:    row.Status,
		IsMenu:    row.IsMenu,
		IsRoute:   row.IsRoute,
		CreatedAt: createdAt,
	}
}

func isDuplicateKey(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
