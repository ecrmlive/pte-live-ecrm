package order

import (
	"testing"
	"time"
)

func TestOrderResponsesExposePointsOrderFields(t *testing.T) {
	group := groupResponse(groupRow{ID: 41, OrderNo: "PG-POINTS-41", PayAmount: 0, TotalQuantity: 2, PayStatus: "paid", ActivityType: 20, PointsAmount: 240, CreatedAt: time.Date(2026, 8, 3, 15, 4, 5, 0, time.Local)})
	if got := group["activity_type"]; got != 20 {
		t.Fatalf("group activity_type=%v, want 20", got)
	}
	if got := group["points_amount"]; got != int64(240) {
		t.Fatalf("group points_amount=%v, want 240", got)
	}
	if got := group["paid"]; got != 1 {
		t.Fatalf("group paid=%v, want 1", got)
	}
	if got := group["create_time"]; got != "2026-08-03 15:04:05" {
		t.Fatalf("group create_time=%v, want 2026-08-03 15:04:05", got)
	}

	child := orderResponse(orderRow{ID: 42, OrderNo: "PS-POINTS-42", MerchantID: 1, MerchantNameSnapshot: "七禧服饰旗舰店", TotalQuantity: 2, Status: "paid", ActivityType: 20, PointsAmount: 240}, nil)
	if got := child["activity_type"]; got != 20 {
		t.Fatalf("order activity_type=%v, want 20", got)
	}
	if got := child["points_amount"]; got != int64(240) {
		t.Fatalf("order points_amount=%v, want 240", got)
	}
}

func TestOrderItemResponseIncludesCommentState(t *testing.T) {
	response := orderItemResponse(orderItemRow{ID: 208, TitleSnapshot: "轻奢羊绒针织衫", Commented: true})
	if got := response["commented"]; got != true {
		t.Fatalf("commented = %v, want true", got)
	}
}
