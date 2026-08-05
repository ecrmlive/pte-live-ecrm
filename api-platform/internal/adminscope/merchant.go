// Package adminscope resolves unified-admin data ranges from qixi_crm_a_*.
package adminscope

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"gorm.io/gorm"
)

var ErrNotConfigured = errors.New("未配置后台数据范围")

type MerchantScope struct {
	Full        bool
	MerchantIDs []uint64
	RegionIDs   []uint64
}

type scopeRow struct {
	ScopeType  string          `gorm:"column:scope_type"`
	ScopeValue json.RawMessage `gorm:"column:scope_value"`
}

type directMerchantScope struct {
	MerchantIDs []uint64 `json:"merchant_ids"`
}

func ResolveMerchantScope(ctx context.Context, adminDB *gorm.DB, claims *authjwt.Claims) (MerchantScope, error) {
	if claims == nil || claims.AdminID == 0 || adminDB == nil {
		return MerchantScope{}, ErrNotConfigured
	}
	if hasRole(claims.Roles, "platform") {
		return MerchantScope{Full: true}, nil
	}
	merchantRole := hasRole(claims.Roles, "merchant")
	regionRole := hasRole(claims.Roles, "region")
	if !merchantRole && !regionRole {
		return MerchantScope{}, ErrNotConfigured
	}
	var rows []scopeRow
	if err := adminDB.WithContext(ctx).Table("qixi_crm_a_data_scope").
		Select("scope_type,scope_value").
		Where("admin_user_id = ? AND scope_type IN ?", claims.AdminID, []string{"merchant", "region"}).
		Find(&rows).Error; err != nil {
		return MerchantScope{}, err
	}
	scope := decodeMerchantScope(rows, merchantRole, regionRole)
	if len(scope.MerchantIDs) == 0 && len(scope.RegionIDs) == 0 {
		return MerchantScope{}, ErrNotConfigured
	}
	return scope, nil
}

func decodeMerchantScope(rows []scopeRow, merchantRole, regionRole bool) MerchantScope {
	scope := MerchantScope{MerchantIDs: []uint64{}, RegionIDs: []uint64{}}
	for _, row := range rows {
		switch row.ScopeType {
		case "merchant":
			if !merchantRole {
				continue
			}
			var direct directMerchantScope
			if json.Unmarshal(row.ScopeValue, &direct) == nil {
				scope.MerchantIDs = appendUnique(scope.MerchantIDs, direct.MerchantIDs...)
			}
		case "region":
			if !regionRole {
				continue
			}
			var regions []uint64
			if json.Unmarshal(row.ScopeValue, &regions) == nil {
				scope.RegionIDs = appendUnique(scope.RegionIDs, regions...)
			}
		}
	}
	return scope
}

func appendUnique(values []uint64, additions ...uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(values)+len(additions))
	for _, value := range values {
		seen[value] = struct{}{}
	}
	for _, value := range additions {
		if value == 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		values = append(values, value)
	}
	return values
}

func hasRole(roles []string, expected string) bool {
	for _, role := range roles {
		if role == expected {
			return true
		}
	}
	return false
}
