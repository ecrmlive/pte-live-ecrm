package identitypersist

import (
	"context"
	"strings"
	"time"

	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/domain/identity"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) FindPlatformAdminByAccount(ctx context.Context, account string) (*identity.SystemAdmin, error) {
	var row identity.SystemAdmin
	err := r.db.WithContext(ctx).
		Where("account = ? AND is_del = 0", account).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) FindPlatformAdminByID(ctx context.Context, id uint) (*identity.SystemAdmin, error) {
	var row identity.SystemAdmin
	err := r.db.WithContext(ctx).
		Where("admin_id = ? AND is_del = 0", id).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) TouchPlatformLogin(ctx context.Context, id uint, ip string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&identity.SystemAdmin{}).
		Where("admin_id = ?", id).
		Updates(map[string]interface{}{
			"last_ip":     ip,
			"last_time":   now,
			"login_count": gorm.Expr("login_count + 1"),
		}).Error
}

func (r *Repo) UpdatePlatformPassword(ctx context.Context, id uint, hash string) error {
	return r.db.WithContext(ctx).Model(&identity.SystemAdmin{}).
		Where("admin_id = ?", id).
		Update("pwd", hash).Error
}

func (r *Repo) FindMerchantAdminByAccount(ctx context.Context, account string) (*identity.MerchantAdmin, error) {
	var row identity.MerchantAdmin
	err := r.db.WithContext(ctx).
		Where("account = ? AND is_del = 0", account).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) FindMerchantAdminByID(ctx context.Context, id uint) (*identity.MerchantAdmin, error) {
	var row identity.MerchantAdmin
	err := r.db.WithContext(ctx).
		Where("merchant_admin_id = ? AND is_del = 0", id).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) FindMerchant(ctx context.Context, merID uint) (*identity.Merchant, error) {
	var row identity.Merchant
	err := r.db.WithContext(ctx).
		Where("mer_id = ? AND is_del = 0", merID).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) TouchMerchantLogin(ctx context.Context, id uint, ip string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&identity.MerchantAdmin{}).
		Where("merchant_admin_id = ?", id).
		Updates(map[string]interface{}{
			"last_ip":     ip,
			"last_time":   now,
			"login_count": gorm.Expr("login_count + 1"),
		}).Error
}

func (r *Repo) UpdateMerchantPassword(ctx context.Context, id uint, hash string) error {
	return r.db.WithContext(ctx).Model(&identity.MerchantAdmin{}).
		Where("merchant_admin_id = ?", id).
		Update("pwd", hash).Error
}

func (r *Repo) FindUserByAccount(ctx context.Context, account string) (*identity.User, error) {
	var row identity.User
	err := r.db.WithContext(ctx).
		Where("account = ? AND cancel_time IS NULL", account).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) FindUserByID(ctx context.Context, id uint) (*identity.User, error) {
	var row identity.User
	err := r.db.WithContext(ctx).
		Where("uid = ? AND cancel_time IS NULL", id).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateUser(ctx context.Context, user *identity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repo) TouchUserLogin(ctx context.Context, id uint, ip string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&identity.User{}).
		Where("uid = ?", id).
		Updates(map[string]interface{}{
			"last_ip":   ip,
			"last_time": now,
		}).Error
}

func (r *Repo) ListUsers(ctx context.Context, page, limit int) ([]identity.User, int64, error) {
	q := r.db.WithContext(ctx).Model(&identity.User{}).Where("cancel_time IS NULL")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []identity.User
	err := q.Order("uid DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) UpdateUserSvip(ctx context.Context, u *identity.User) error {
	return r.db.WithContext(ctx).Model(&identity.User{}).Where("uid = ?", u.UID).
		Updates(map[string]interface{}{
			"is_svip": u.IsSvip, "svip_endtime": u.SvipEndtime,
		}).Error
}

func (r *Repo) FindStoreServiceByAccount(ctx context.Context, account string) (*identity.StoreService, error) {
	var row identity.StoreService
	err := r.db.WithContext(ctx).
		Where("account = ? AND is_del = 0", account).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) FindStoreServiceByID(ctx context.Context, id uint) (*identity.StoreService, error) {
	var row identity.StoreService
	err := r.db.WithContext(ctx).
		Where("service_id = ? AND is_del = 0", id).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) ListMenusByIDs(ctx context.Context, isMer uint8, ids []uint) ([]identity.SystemMenu, error) {
	var rows []identity.SystemMenu
	q := r.db.WithContext(ctx).
		Where("is_mer = ? AND is_show = 1 AND is_menu = 1", isMer)
	if len(ids) > 0 {
		q = q.Where("menu_id IN ?", ids)
	}
	err := q.Order("sort DESC, menu_id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) ListButtonMenusByIDs(ctx context.Context, isMer uint8, ids []uint) ([]identity.SystemMenu, error) {
	var rows []identity.SystemMenu
	q := r.db.WithContext(ctx).
		Where("is_mer = ? AND is_menu = 2", isMer)
	if len(ids) > 0 {
		q = q.Where("menu_id IN ?", ids)
	}
	err := q.Order("sort DESC, menu_id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) ListMenusManage(ctx context.Context, isMer uint8) ([]identity.SystemMenu, error) {
	// 角色勾选树：菜单 + 按钮权限
	var rows []identity.SystemMenu
	err := r.db.WithContext(ctx).
		Where("is_mer = ?", isMer).
		Where("is_menu IN ?", []uint8{1, 2}).
		Order("sort DESC, menu_id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) GetMenu(ctx context.Context, id uint) (*identity.SystemMenu, error) {
	var row identity.SystemMenu
	err := r.db.WithContext(ctx).Where("menu_id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) UpdateMenu(ctx context.Context, row *identity.SystemMenu) error {
	return r.db.WithContext(ctx).Model(row).Where("menu_id = ?", row.MenuID).Updates(map[string]interface{}{
		"menu_name": row.MenuName, "sort": row.Sort, "is_show": row.IsShow,
	}).Error
}

func (r *Repo) ListRoleRules(ctx context.Context, roleIDs []uint) ([]string, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}
	var roles []identity.SystemRole
	err := r.db.WithContext(ctx).
		Where("role_id IN ? AND status = 1", roleIDs).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(roles))
	for _, role := range roles {
		if strings.TrimSpace(role.Rules) != "" {
			out = append(out, role.Rules)
		}
	}
	return out, nil
}

func (r *Repo) ListRoles(ctx context.Context, merID uint, page, limit int) ([]identity.SystemRole, int64, error) {
	q := r.db.WithContext(ctx).Model(&identity.SystemRole{}).Where("mer_id = ?", merID)
	if merID == 0 {
		q = q.Where("is_agent IN ?", []int8{0, 1})
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []identity.SystemRole
	err := q.Order("role_id ASC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) ListMerchantRoles(ctx context.Context, merID uint, page, limit int) ([]identity.SystemRole, int64, error) {
	// 共享商户模板（mer_id=0,is_agent=2）+ 本店角色
	q := r.db.WithContext(ctx).Model(&identity.SystemRole{}).
		Where("(mer_id = 0 AND is_agent = 2) OR mer_id = ?", merID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []identity.SystemRole
	err := q.Order("mer_id ASC, role_id ASC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetRole(ctx context.Context, id uint) (*identity.SystemRole, error) {
	var row identity.SystemRole
	err := r.db.WithContext(ctx).Where("role_id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateRole(ctx context.Context, row *identity.SystemRole) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateRole(ctx context.Context, row *identity.SystemRole) error {
	return r.db.WithContext(ctx).Model(row).Where("role_id = ?", row.RoleID).Updates(map[string]interface{}{
		"role_name": row.RoleName, "rules": row.Rules, "status": row.Status,
		"is_agent": row.IsAgent, "circle_id": row.CircleID,
	}).Error
}

func (r *Repo) ListPlatformAdmins(ctx context.Context, page, limit int) ([]identity.SystemAdmin, int64, error) {
	q := r.db.WithContext(ctx).Model(&identity.SystemAdmin{}).Where("is_del = 0")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []identity.SystemAdmin
	err := q.Order("admin_id ASC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) CreatePlatformAdmin(ctx context.Context, row *identity.SystemAdmin) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdatePlatformAdmin(ctx context.Context, row *identity.SystemAdmin) error {
	return r.db.WithContext(ctx).Model(row).Where("admin_id = ?", row.AdminID).Updates(map[string]interface{}{
		"real_name": row.RealName, "phone": row.Phone, "roles": row.Roles,
		"status": row.Status, "pwd": row.Pwd, "is_agent": row.IsAgent,
		"region_ids": row.RegionIDs, "circle_agent_id": row.CircleAgentID,
	}).Error
}

func (r *Repo) ListMerchantAdmins(ctx context.Context, merID uint, page, limit int) ([]identity.MerchantAdmin, int64, error) {
	q := r.db.WithContext(ctx).Model(&identity.MerchantAdmin{}).Where("mer_id = ? AND is_del = 0", merID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []identity.MerchantAdmin
	err := q.Order("merchant_admin_id ASC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) CreateMerchantAdmin(ctx context.Context, row *identity.MerchantAdmin) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateMerchantAdmin(ctx context.Context, row *identity.MerchantAdmin) error {
	return r.db.WithContext(ctx).Model(row).Where("merchant_admin_id = ?", row.MerchantAdminID).Updates(map[string]interface{}{
		"real_name": row.RealName, "phone": row.Phone, "roles": row.Roles,
		"status": row.Status, "pwd": row.Pwd,
	}).Error
}

func (r *Repo) ListStoreServices(ctx context.Context, merID uint, page, limit int) ([]identity.StoreService, int64, error) {
	q := r.db.WithContext(ctx).Model(&identity.StoreService{}).Where("mer_id = ? AND is_del = 0", merID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []identity.StoreService
	err := q.Order("service_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) CreateStoreService(ctx context.Context, row *identity.StoreService) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateStoreService(ctx context.Context, row *identity.StoreService) error {
	return r.db.WithContext(ctx).Model(row).Where("service_id = ?", row.ServiceID).Updates(map[string]interface{}{
		"nickname": row.Nickname, "phone": row.Phone, "status": row.Status,
		"is_open": row.IsOpen, "is_verify": row.IsVerify, "is_goods": row.IsGoods, "pwd": row.Pwd,
	}).Error
}
