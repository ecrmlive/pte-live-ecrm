package queryfilter

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ApplyCreatedAtRange filters q by date_from/date_to query params (YYYY-MM-DD, inclusive).
func ApplyCreatedAtRange(q *gorm.DB, c *gin.Context, column string) *gorm.DB {
	if from := strings.TrimSpace(c.Query("date_from")); from != "" {
		if t, err := time.ParseInLocation("2006-01-02", from, time.Local); err == nil {
			q = q.Where(column+" >= ?", t)
		}
	}
	if to := strings.TrimSpace(c.Query("date_to")); to != "" {
		if t, err := time.ParseInLocation("2006-01-02", to, time.Local); err == nil {
			q = q.Where(column+" < ?", t.AddDate(0, 0, 1))
		}
	}
	return q
}

// Page reads page/limit from query with sane defaults.
func Page(c *gin.Context) (int, int) {
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
