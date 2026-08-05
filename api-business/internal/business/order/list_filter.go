package order

import (
	"errors"
	"strings"
)

var ErrOrderListStatus = errors.New("订单状态筛选参数错误")
var ErrOrderListKeyword = errors.New("订单搜索关键词不合法")
var ErrOrderListFulfillmentStatus = errors.New("订单履约状态筛选参数错误")

// NormalizeGroupPayStatus turns the public order-list filter into the small,
// explicit group-payment state machine.  An empty value represents all orders.
func NormalizeGroupPayStatus(value string) (string, error) {
	status := strings.TrimSpace(value)
	if status == "" || status == "all" {
		return "", nil
	}
	switch status {
	case "pending", "paid", "closed":
		return status, nil
	default:
		return "", ErrOrderListStatus
	}
}

func NormalizeGroupFulfillmentStatus(value string) (string, error) {
	status := strings.TrimSpace(value)
	if status == "" || status == "all" {
		return "", nil
	}
	switch status {
	case "awaiting_fulfillment", "awaiting_receipt", "awaiting_comment":
		return status, nil
	default:
		return "", ErrOrderListFulfillmentStatus
	}
}

func NormalizeOrderListKeyword(value string) (string, error) {
	keyword := strings.TrimSpace(value)
	if len([]rune(keyword)) > 64 {
		return "", ErrOrderListKeyword
	}
	return keyword, nil
}
