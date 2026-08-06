package listquery

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminFilter is the common merchant list query shape used by admin grids.
type AdminFilter struct {
	Keyword  string
	Status   *int8
	DateFrom string
	DateTo   string
}

func ParseAdminFilter(c *gin.Context) AdminFilter {
	filter := AdminFilter{
		Keyword:  strings.TrimSpace(c.Query("keyword")),
		DateFrom: strings.TrimSpace(c.Query("date_from")),
		DateTo:   strings.TrimSpace(c.Query("date_to")),
	}
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		if value, err := strconv.Atoi(raw); err == nil && (value == 0 || value == 1) {
			v := int8(value)
			filter.Status = &v
		}
	}
	return filter
}

func ApplyTimeColumnDateRange(q *gorm.DB, column, from, to string) *gorm.DB {
	if from != "" {
		if t, err := time.ParseInLocation("2006-01-02", from, time.Local); err == nil {
			q = q.Where(column+" >= ?", t)
		}
	}
	if to != "" {
		if t, err := time.ParseInLocation("2006-01-02", to, time.Local); err == nil {
			q = q.Where(column+" < ?", t.AddDate(0, 0, 1))
		}
	}
	return q
}

func ApplyUnixColumnDateRange(q *gorm.DB, column, from, to string) *gorm.DB {
	if from != "" {
		if t, err := time.ParseInLocation("2006-01-02", from, time.Local); err == nil {
			q = q.Where(column+" >= ?", t.Unix())
		}
	}
	if to != "" {
		if t, err := time.ParseInLocation("2006-01-02", to, time.Local); err == nil {
			q = q.Where(column+" < ?", t.AddDate(0, 0, 1).Unix())
		}
	}
	return q
}

func ApplyKeywordLike(q *gorm.DB, keyword string, columns ...string) *gorm.DB {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" || len(columns) == 0 {
		return q
	}
	like := "%" + keyword + "%"
	parts := make([]string, 0, len(columns))
	args := make([]any, 0, len(columns))
	for _, column := range columns {
		parts = append(parts, column+" LIKE ?")
		args = append(args, like)
	}
	return q.Where(strings.Join(parts, " OR "), args...)
}
