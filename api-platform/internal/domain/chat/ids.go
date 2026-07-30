package chat

// 数值 IM user_id 分段，避免与其他业务租户冲突（mall app_id=30001）。
const (
	IMUserAppBase     int64 = 1_000_000_000
	IMUserServiceBase int64 = 2_000_000_000
)

func NumericIMUserID(portal string, localID uint) int64 {
	switch portal {
	case "service":
		return IMUserServiceBase + int64(localID)
	default:
		return IMUserAppBase + int64(localID)
	}
}
