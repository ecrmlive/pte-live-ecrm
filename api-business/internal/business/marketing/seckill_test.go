package marketing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"time"
)

func TestSlotActiveUsesFourHourWindow(t *testing.T) {
	base := time.Date(2026, time.July, 31, 14, 30, 0, 0, time.Local)
	if !slotActive("14:00", base) {
		t.Fatal("14:00 场次在 14:30 应为进行中")
	}
	if slotActive("19:00", base) {
		t.Fatal("19:00 场次在 14:30 不应为进行中")
	}
	if slotActive("invalid", base) {
		t.Fatal("非法场次不能激活")
	}
}

func TestActivityActiveRequiresDateAndSlot(t *testing.T) {
	now := time.Date(2026, time.July, 31, 14, 30, 0, 0, time.Local)
	start := now.Add(-24 * time.Hour)
	end := now.Add(24 * time.Hour)
	activity := seckillActivityView{StartsAt: &start, EndsAt: &end}
	if !activityActive(activity, seckillRules{TimeSlots: []string{"14:00"}}, now) {
		t.Fatal("有效日期和场次应激活活动")
	}
	if activityActive(activity, seckillRules{TimeSlots: []string{"19:00"}}, now) {
		t.Fatal("非当前场次不应激活活动")
	}
	past := now.Add(-time.Hour)
	activity.EndsAt = &past
	if activityActive(activity, seckillRules{TimeSlots: []string{"14:00"}}, now) {
		t.Fatal("过期活动不应激活")
	}
}

func TestSeckillGetRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewSeckillHandler(nil).Register(router)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/seckill/actives/0", nil))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("invalid seckill activity ID status=%d, want %d", recorder.Code, http.StatusBadRequest)
	}
}
