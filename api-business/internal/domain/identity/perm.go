package identity

import (
	"context"
	"strings"
)

// 商户按钮权限 menu_id（is_menu=2）；与 sql/024–039 对齐。
const (
	MerPermOrderDeliver      uint = 127
	MerPermOrderVerify       uint = 128
	MerPermRefundApprove     uint = 129
	MerPermRefundReject      uint = 130
	MerPermProductShow       uint = 132
	MerPermProductStock      uint = 133
	MerPermProductCreate     uint = 134
	MerPermProductDelete     uint = 135
	MerPermSeckillToggle     uint = 136
	MerPermCombinationToggle uint = 137
	MerPermCouponToggle      uint = 138
	MerPermPresellToggle     uint = 139
	MerPermAssistToggle      uint = 140
	MerPermCouponCreate      uint = 141
	MerPermCouponDelete      uint = 142
	MerPermSeckillCreate     uint = 143
	MerPermSeckillDelete     uint = 144
	MerPermCombinationCreate uint = 145
	MerPermCombinationDelete uint = 146
	MerPermPresellCreate     uint = 147
	MerPermPresellDelete     uint = 148
	MerPermAssistCreate      uint = 149
	MerPermAssistDelete      uint = 150
	MerPermBroadcastCreate   uint = 151
	MerPermBroadcastDelete   uint = 152
	MerPermReservationConfig uint = 153
	MerPermSvipUpdate        uint = 154
	MerPermBroadcastLive     uint = 155
	MerPermBroadcastGoods    uint = 156
	MerPermAttachmentUpload  uint = 157
	MerPermAttachmentDelete  uint = 158
	MerPermShopUpdate        uint = 159
	MerPermStaffWrite        uint = 160
	MerPermAdminsWrite       uint = 161
	MerPermRolesWrite        uint = 162
	MerPermReplyWrite        uint = 164
	MerPermCommunityCreate   uint = 165
	MerPermCommunityUpdate   uint = 166
	MerPermCommunityDelete   uint = 167
	MerPermCouponSend        uint = 8667
)

// 平台按钮权限 menu_id（is_menu=2）；与 sql/027、029–039 对齐。
const (
	PlatPermRefundApprove    uint = 34
	PlatPermRefundReject     uint = 35
	PlatPermWithdrawApprove  uint = 36
	PlatPermWithdrawReject   uint = 37
	PlatPermCommunityAudit   uint = 38
	PlatPermCommunityDelete  uint = 39
	PlatPermBroadcastAudit   uint = 40
	PlatPermCouponToggle     uint = 41
	PlatPermCouponCreate     uint = 42
	PlatPermCouponDelete     uint = 43
	PlatPermSvipUpdate       uint = 44
	PlatPermAttachmentUpload uint = 45
	PlatPermAttachmentDelete uint = 46
	PlatPermDiyCreate        uint = 47
	PlatPermDiyUpdate        uint = 48
	PlatPermDiyDelete        uint = 49
	PlatPermDiyActive        uint = 50
	PlatPermDiyPick          uint = 51
	PlatPermAgreementUpdate  uint = 53
	PlatPermCloudConfigWrite uint = 20902
	PlatPermProductAudit     uint = 20903
	PlatPermCategoryManage   uint = 20904
	PlatPermBrandManage      uint = 20905
	PlatPermNoticeManage     uint = 20906
	PlatPermUserLabelManage  uint = 20907
	PlatPermUserGroupManage  uint = 20908
	// 区域代理（CRMEB business-zones）权限节点。
	PlatPermCircleManage      uint = 9924
	PlatPermCircleAgentManage uint = 9926
	PlatPermCircleAgentReview uint = 9928
)

// HasMerchantMenu 校验商户管理员是否拥有指定菜单/按钮（level=0 主账号放行）。
func (s *Service) HasMerchantMenu(ctx context.Context, adminID, menuID uint) (bool, error) {
	if menuID == 0 {
		return false, ErrBadParam
	}
	profile, err := s.MerchantProfile(ctx, adminID)
	if err != nil {
		return false, err
	}
	return s.hasMenu(ctx, 2, profile.Roles, profile.Level == 0, menuID)
}

// RequireMerchantMenu 无权限时返回 ErrNoPerm。
func (s *Service) RequireMerchantMenu(ctx context.Context, adminID, menuID uint) error {
	ok, err := s.HasMerchantMenu(ctx, adminID, menuID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNoPerm
	}
	return nil
}

// MerchantPermissionPaths 返回当前账号拥有的按钮 path 列表（供前端隐藏操作）。
func (s *Service) MerchantPermissionPaths(ctx context.Context, adminID uint) ([]string, error) {
	profile, err := s.MerchantProfile(ctx, adminID)
	if err != nil {
		return nil, err
	}
	return s.permissionPaths(ctx, 2, profile.Roles, profile.Level == 0)
}

// HasPlatformMenu 校验平台管理员是否拥有指定菜单/按钮（level=0 超管放行）。
func (s *Service) HasPlatformMenu(ctx context.Context, adminID, menuID uint) (bool, error) {
	if menuID == 0 {
		return false, ErrBadParam
	}
	profile, err := s.PlatformProfile(ctx, adminID)
	if err != nil {
		return false, err
	}
	return s.hasMenu(ctx, 1, profile.Roles, profile.Level == 0, menuID)
}

// RequirePlatformMenu 无权限时返回 ErrNoPerm。
func (s *Service) RequirePlatformMenu(ctx context.Context, adminID, menuID uint) error {
	ok, err := s.HasPlatformMenu(ctx, adminID, menuID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNoPerm
	}
	return nil
}

// PlatformPermissionPaths 返回平台按钮 path 列表。
func (s *Service) PlatformPermissionPaths(ctx context.Context, adminID uint) ([]string, error) {
	profile, err := s.PlatformProfile(ctx, adminID)
	if err != nil {
		return nil, err
	}
	return s.permissionPaths(ctx, 1, profile.Roles, profile.Level == 0)
}

func (s *Service) hasMenu(ctx context.Context, isMer uint8, rolesCSV string, isSuper bool, menuID uint) (bool, error) {
	if isSuper {
		return true, nil
	}
	ids, err := s.roleMenuIDs(ctx, rolesCSV)
	if err != nil {
		return false, err
	}
	for _, id := range ids {
		if id == menuID {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) permissionPaths(ctx context.Context, isMer uint8, rolesCSV string, isSuper bool) ([]string, error) {
	var ids []uint
	if !isSuper {
		var err error
		ids, err = s.roleMenuIDs(ctx, rolesCSV)
		if err != nil {
			return nil, err
		}
		if len(ids) == 0 {
			return []string{}, nil
		}
	}
	rows, err := s.store.ListButtonMenusByIDs(ctx, isMer, ids)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	seen := map[string]struct{}{}
	for _, row := range rows {
		p := strings.TrimSpace(row.Path)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out, nil
}

func (s *Service) roleMenuIDs(ctx context.Context, rolesCSV string) ([]uint, error) {
	roleIDs := parseIDs(rolesCSV)
	rules, err := s.store.ListRoleRules(ctx, roleIDs)
	if err != nil {
		return nil, err
	}
	return parseIDs(strings.Join(rules, ",")), nil
}
