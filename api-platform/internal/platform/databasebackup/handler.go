package databasebackup

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const defaultBackupDir = "/app/backups"

var tableNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]{2,127}$`)

type Handler struct {
	adminDB    *gorm.DB
	backupDir  string
	businessDB *gorm.DB
}

type databaseScope string

const (
	scopeAdmin    databaseScope = "admin"
	scopeBusiness databaseScope = "business"
)

type databaseTable struct {
	DatabaseScope string     `json:"database_scope"`
	Engine        string     `json:"engine"`
	RowCount      int64      `json:"row_count"`
	SizeBytes     int64      `json:"size_bytes"`
	TableComment  string     `json:"table_comment"`
	TableName     string     `json:"table_name"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type tableDetailColumn struct {
	ColumnComment string  `json:"column_comment"`
	ColumnDefault *string `json:"column_default"`
	ColumnName    string  `json:"column_name"`
	ColumnType    string  `json:"column_type"`
	Extra         string  `json:"extra"`
	IsNullable    string  `json:"is_nullable"`
}

type backupRecord struct {
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
	CreatedBy     uint       `gorm:"column:created_by" json:"created_by"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	FileName      string     `gorm:"column:file_name" json:"file_name"`
	ID            uint64     `gorm:"column:id" json:"id"`
	Scope         string     `gorm:"column:database_scope" json:"database_scope"`
	SizeBytes     int64      `gorm:"column:size_bytes" json:"size_bytes"`
	Status        string     `gorm:"column:status" json:"status"`
	TableCount    int        `gorm:"column:table_count" json:"table_count"`
	TableNamesRaw string     `gorm:"column:table_names_json" json:"-"`
	TableNames    []string   `gorm:"-" json:"table_names"`
}

type tableActionInput struct {
	Scope  databaseScope `json:"scope"`
	Tables []string      `json:"tables"`
}

func NewHandler(adminDB, businessDB *gorm.DB) *Handler {
	return newHandler(adminDB, businessDB, defaultBackupDir)
}

func newHandler(adminDB, businessDB *gorm.DB, backupDir string) *Handler {
	return &Handler{adminDB: adminDB, businessDB: businessDB, backupDir: backupDir}
}

func (h *Handler) Register(r gin.IRoutes) {
	platformOnly := middleware.RequireAdminRoles("platform")
	read := middleware.RequireAdminMenu(h.adminDB, "maintain.backup")
	manage := middleware.RequireAdminMenu(h.adminDB, "maintain.backup.manage")

	r.GET("/maintain/database/tables", platformOnly, read, h.ListTables)
	r.GET("/maintain/database/tables/:scope/:name", platformOnly, read, h.GetTableDetail)
	r.POST("/maintain/database/backups", platformOnly, manage, h.CreateBackup)
	r.GET("/maintain/database/backups", platformOnly, read, h.ListBackups)
	r.GET("/maintain/database/backups/:id/download", platformOnly, read, h.DownloadBackup)
	r.DELETE("/maintain/database/backups/:id", platformOnly, manage, h.DeleteBackup)
	r.POST("/maintain/database/tables/optimize", platformOnly, manage, h.OptimizeTables)
	r.POST("/maintain/database/tables/repair", platformOnly, manage, h.RepairTables)
}

func (h *Handler) ListTables(c *gin.Context) {
	page, limit := parsePage(c)
	rows, err := h.allTables(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "读取数据库表失败")
		return
	}

	total := len(rows)
	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}
	pageRows := rows[start:end]
	for index := range pageRows {
		count, countErr := h.countRows(c.Request.Context(), databaseScope(pageRows[index].DatabaseScope), pageRows[index].TableName)
		if countErr != nil {
			response.Fail(c, http.StatusInternalServerError, "读取数据行数失败")
			return
		}
		pageRows[index].RowCount = count
	}
	response.OK(c, gin.H{"list": pageRows, "total": total})
}

func (h *Handler) GetTableDetail(c *gin.Context) {
	scope, db, err := h.databaseFor(databaseScope(c.Param("scope")))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "数据库范围错误")
		return
	}
	name := c.Param("name")
	if !validTableName(scope, name) || !h.tableExists(c.Request.Context(), db, name) {
		response.Fail(c, http.StatusNotFound, "数据表不存在")
		return
	}
	var columns []tableDetailColumn
	err = db.WithContext(c.Request.Context()).Raw(`
SELECT column_name, column_type, column_default, is_nullable, extra, column_comment
FROM information_schema.columns
WHERE table_schema = DATABASE() AND table_name = ?
ORDER BY ordinal_position`, name).Scan(&columns).Error
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "读取表结构失败")
		return
	}
	response.OK(c, gin.H{"columns": columns, "database_scope": scope, "table_name": name})
}

func (h *Handler) CreateBackup(c *gin.Context) {
	var in tableActionInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	scope, db, tables, err := h.validateTables(c.Request.Context(), in)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := os.MkdirAll(h.backupDir, 0o750); err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建备份目录失败")
		return
	}
	fileName := backupFileName(scope, tables)
	finalPath := filepath.Join(h.backupDir, fileName)
	if !isSafeBackupPath(h.backupDir, finalPath) {
		response.Fail(c, http.StatusBadRequest, "备份文件名错误")
		return
	}
	temporary, err := os.CreateTemp(h.backupDir, ".backup-*.sql")
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建备份文件失败")
		return
	}
	temporaryPath := temporary.Name()
	defer func() { _ = os.Remove(temporaryPath) }()

	tx := db.WithContext(c.Request.Context()).Begin(&sql.TxOptions{ReadOnly: true})
	if tx.Error != nil {
		_ = temporary.Close()
		response.Fail(c, http.StatusInternalServerError, "创建备份快照失败")
		return
	}
	defer func() { _ = tx.Rollback().Error }()
	if _, err = io.WriteString(temporary, "-- qixi live ecrm local database backup\nSET FOREIGN_KEY_CHECKS=0;\n"); err == nil {
		for _, table := range tables {
			if err = dumpTable(c.Request.Context(), temporary, tx, table); err != nil {
				break
			}
		}
	}
	if err == nil {
		_, err = io.WriteString(temporary, "SET FOREIGN_KEY_CHECKS=1;\n")
	}
	if closeErr := temporary.Close(); err == nil {
		err = closeErr
	}
	if err == nil {
		err = tx.Commit().Error
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "生成数据库备份失败")
		return
	}
	if err = os.Rename(temporaryPath, finalPath); err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存备份文件失败")
		return
	}
	info, err := os.Stat(finalPath)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "读取备份文件失败")
		return
	}
	namesJSON, _ := json.Marshal(tables)
	record := map[string]any{
		"file_name":        fileName,
		"database_scope":   string(scope),
		"table_names_json": string(namesJSON),
		"table_count":      len(tables),
		"size_bytes":       info.Size(),
		"status":           "ready",
		"created_by":       middleware.AdminID(c),
	}
	if err = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_database_backup").Create(record).Error; err != nil {
		_ = os.Remove(finalPath)
		response.Fail(c, http.StatusInternalServerError, "登记备份记录失败")
		return
	}
	h.writeMaintenanceLog(c.Request.Context(), "backup", scope, tables, middleware.AdminID(c), "本地逻辑备份已生成")
	response.OK(c, gin.H{"file_name": fileName, "size_bytes": info.Size(), "table_count": len(tables)})
}

func (h *Handler) ListBackups(c *gin.Context) {
	page, limit := parsePage(c)
	var total int64
	query := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_database_backup").Where("status <> ?", "deleted")
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "读取备份记录失败")
		return
	}
	var rows []backupRecord
	if err := query.Order("id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "读取备份记录失败")
		return
	}
	for index := range rows {
		_ = json.Unmarshal([]byte(rows[index].TableNamesRaw), &rows[index].TableNames)
	}
	response.OK(c, gin.H{"list": rows, "total": total})
}

func (h *Handler) DownloadBackup(c *gin.Context) {
	record, ok := h.loadReadyBackup(c)
	if !ok {
		return
	}
	path := filepath.Join(h.backupDir, record.FileName)
	if !isSafeBackupPath(h.backupDir, path) {
		response.Fail(c, http.StatusBadRequest, "备份文件错误")
		return
	}
	if _, err := os.Stat(path); err != nil {
		response.Fail(c, http.StatusNotFound, "备份文件不存在")
		return
	}
	c.FileAttachment(path, record.FileName)
}

func (h *Handler) DeleteBackup(c *gin.Context) {
	record, ok := h.loadReadyBackup(c)
	if !ok {
		return
	}
	path := filepath.Join(h.backupDir, record.FileName)
	if !isSafeBackupPath(h.backupDir, path) {
		response.Fail(c, http.StatusBadRequest, "备份文件错误")
		return
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		response.Fail(c, http.StatusInternalServerError, "删除备份文件失败")
		return
	}
	now := time.Now()
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_database_backup").Where("id = ?", record.ID).
		Updates(map[string]any{"status": "deleted", "deleted_at": now}).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除备份记录失败")
		return
	}
	h.writeMaintenanceLog(c.Request.Context(), "delete", databaseScope(record.Scope), record.TableNames, middleware.AdminID(c), "已删除本地备份文件")
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) OptimizeTables(c *gin.Context) {
	h.maintainTables(c, "optimize")
}

func (h *Handler) RepairTables(c *gin.Context) {
	h.maintainTables(c, "repair")
}

func (h *Handler) maintainTables(c *gin.Context, action string) {
	var in tableActionInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	scope, db, tables, err := h.validateTables(c.Request.Context(), in)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	results := make([]string, 0, len(tables))
	for _, table := range tables {
		engine, engineErr := h.tableEngine(c.Request.Context(), db, table)
		if engineErr != nil {
			response.Fail(c, http.StatusInternalServerError, "读取数据表引擎失败")
			return
		}
		statement := "OPTIMIZE TABLE " + quoteIdentifier(table)
		if action == "repair" {
			if strings.EqualFold(engine, "InnoDB") {
				statement = "ANALYZE TABLE " + quoteIdentifier(table)
			} else {
				statement = "REPAIR TABLE " + quoteIdentifier(table)
			}
		}
		if err := db.WithContext(c.Request.Context()).Exec(statement).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "执行表维护失败")
			return
		}
		if action == "repair" && strings.EqualFold(engine, "InnoDB") {
			results = append(results, table+"（已执行统计分析）")
		} else {
			results = append(results, table)
		}
	}
	detail := "已执行 " + action + "：" + strings.Join(results, "、")
	h.writeMaintenanceLog(c.Request.Context(), action, scope, tables, middleware.AdminID(c), detail)
	response.OK(c, gin.H{"ok": true, "tables": results})
}

func (h *Handler) allTables(ctx context.Context) ([]databaseTable, error) {
	adminRows, err := listTableMetadata(ctx, h.adminDB, scopeAdmin)
	if err != nil {
		return nil, err
	}
	businessRows, err := listTableMetadata(ctx, h.businessDB, scopeBusiness)
	if err != nil {
		return nil, err
	}
	rows := append(adminRows, businessRows...)
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].TableName < rows[j].TableName
	})
	return rows, nil
}

func listTableMetadata(ctx context.Context, db *gorm.DB, scope databaseScope) ([]databaseTable, error) {
	if db == nil {
		return nil, errors.New("database unavailable")
	}
	prefix := "qixi_crm_a_%"
	if scope == scopeBusiness {
		prefix = "qixi_crm_b_%"
	}
	var rows []databaseTable
	err := db.WithContext(ctx).Raw(`
SELECT table_name, table_comment, engine, COALESCE(data_length, 0) + COALESCE(index_length, 0) AS size_bytes, update_time AS updated_at
FROM information_schema.tables
WHERE table_schema = DATABASE() AND table_type = 'BASE TABLE' AND table_name LIKE ?
ORDER BY table_name`, prefix).Scan(&rows).Error
	for index := range rows {
		rows[index].DatabaseScope = string(scope)
		if strings.TrimSpace(rows[index].TableComment) == "" {
			rows[index].TableComment = "—"
		}
	}
	return rows, err
}

func (h *Handler) countRows(ctx context.Context, scope databaseScope, table string) (int64, error) {
	_, db, err := h.databaseFor(scope)
	if err != nil {
		return 0, err
	}
	var count int64
	err = db.WithContext(ctx).Raw("SELECT COUNT(*) FROM " + quoteIdentifier(table)).Scan(&count).Error
	return count, err
}

func (h *Handler) tableExists(ctx context.Context, db *gorm.DB, name string) bool {
	var total int64
	err := db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ? AND table_type = 'BASE TABLE'`, name).Scan(&total).Error
	return err == nil && total == 1
}

func (h *Handler) tableEngine(ctx context.Context, db *gorm.DB, name string) (string, error) {
	var engine string
	err := db.WithContext(ctx).Raw(`SELECT engine FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?`, name).Scan(&engine).Error
	return engine, err
}

func (h *Handler) validateTables(ctx context.Context, input tableActionInput) (databaseScope, *gorm.DB, []string, error) {
	scope, db, err := h.databaseFor(input.Scope)
	if err != nil {
		return "", nil, nil, errors.New("数据库范围错误")
	}
	if len(input.Tables) == 0 || len(input.Tables) > 30 {
		return "", nil, nil, errors.New("请选择 1 至 30 张数据表")
	}
	unique := make(map[string]struct{}, len(input.Tables))
	tables := make([]string, 0, len(input.Tables))
	for _, name := range input.Tables {
		name = strings.TrimSpace(name)
		if !validTableName(scope, name) || !h.tableExists(ctx, db, name) {
			return "", nil, nil, errors.New("存在无效数据表")
		}
		if _, exists := unique[name]; exists {
			continue
		}
		unique[name] = struct{}{}
		tables = append(tables, name)
	}
	if len(tables) == 0 {
		return "", nil, nil, errors.New("请选择数据表")
	}
	sort.Strings(tables)
	return scope, db, tables, nil
}

func (h *Handler) databaseFor(scope databaseScope) (databaseScope, *gorm.DB, error) {
	switch scope {
	case scopeAdmin:
		return scope, h.adminDB, nil
	case scopeBusiness:
		return scope, h.businessDB, nil
	default:
		return "", nil, errors.New("unknown scope")
	}
}

func (h *Handler) loadReadyBackup(c *gin.Context) (backupRecord, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "备份记录错误")
		return backupRecord{}, false
	}
	var record backupRecord
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_database_backup").Where("id = ? AND status = ?", id, "ready").Take(&record).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "备份记录不存在")
		return backupRecord{}, false
	}
	_ = json.Unmarshal([]byte(record.TableNamesRaw), &record.TableNames)
	return record, true
}

func (h *Handler) writeMaintenanceLog(ctx context.Context, action string, scope databaseScope, tables []string, operatorID uint, detail string) {
	tablesJSON, _ := json.Marshal(tables)
	_ = h.adminDB.WithContext(ctx).Table("qixi_crm_a_database_maintenance_log").Create(map[string]any{
		"action": action, "database_scope": string(scope), "table_names_json": string(tablesJSON),
		"operator_admin_id": operatorID, "detail": detail,
	}).Error
}

func dumpTable(ctx context.Context, writer io.Writer, db *gorm.DB, table string) error {
	row := db.WithContext(ctx).Raw("SHOW CREATE TABLE " + quoteIdentifier(table)).Row()
	var actualName, createStatement string
	if err := row.Scan(&actualName, &createStatement); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "\nDROP TABLE IF EXISTS %s;\n%s;\n", quoteIdentifier(actualName), createStatement); err != nil {
		return err
	}
	rows, err := db.WithContext(ctx).Raw("SELECT * FROM " + quoteIdentifier(table)).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	columnSQL := make([]string, 0, len(columns))
	for _, column := range columns {
		columnSQL = append(columnSQL, quoteIdentifier(column))
	}
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for index := range values {
			pointers[index] = &values[index]
		}
		if err := rows.Scan(pointers...); err != nil {
			return err
		}
		valueSQL := make([]string, 0, len(values))
		for _, value := range values {
			valueSQL = append(valueSQL, sqlValue(value))
		}
		if _, err := fmt.Fprintf(writer, "INSERT INTO %s (%s) VALUES (%s);\n", quoteIdentifier(table), strings.Join(columnSQL, ","), strings.Join(valueSQL, ",")); err != nil {
			return err
		}
	}
	return rows.Err()
}

func sqlValue(value any) string {
	switch typed := value.(type) {
	case nil:
		return "NULL"
	case []byte:
		return "X'" + hex.EncodeToString(typed) + "'"
	case time.Time:
		return "'" + typed.Format("2006-01-02 15:04:05.999999") + "'"
	case bool:
		if typed {
			return "1"
		}
		return "0"
	case string:
		return "'" + escapeSQLString(typed) + "'"
	default:
		return fmt.Sprint(typed)
	}
}

func escapeSQLString(value string) string {
	replacer := strings.NewReplacer("\\", "\\\\", "'", "''", "\x00", "\\0", "\n", "\\n", "\r", "\\r", "\x1a", "\\Z")
	return replacer.Replace(value)
}

func quoteIdentifier(value string) string {
	return "`" + strings.ReplaceAll(value, "`", "``") + "`"
}

func validTableName(scope databaseScope, name string) bool {
	prefix := "qixi_crm_a_"
	if scope == scopeBusiness {
		prefix = "qixi_crm_b_"
	}
	return strings.HasPrefix(name, prefix) && tableNamePattern.MatchString(name)
}

func backupFileName(scope databaseScope, tables []string) string {
	return fmt.Sprintf("ecrm_%s_%s_%02d-tables.sql", scope, time.Now().Format("20060102_150405_000000"), len(tables))
}

func isSafeBackupPath(baseDir, path string) bool {
	base, err := filepath.Abs(baseDir)
	if err != nil {
		return false
	}
	target, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	return strings.HasPrefix(target, base+string(os.PathSeparator)) && strings.HasSuffix(target, ".sql")
}

func parsePage(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return page, limit
}
