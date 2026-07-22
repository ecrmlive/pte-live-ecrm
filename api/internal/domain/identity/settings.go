package identity

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

type AdminSaveInput struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Roles    string `json:"roles"`
	Status   *uint8 `json:"status"`
}

type RoleSaveInput struct {
	RoleName string  `json:"role_name"`
	Rules    string  `json:"rules"`
	MenuIDs  *[]uint `json:"menu_ids"`
	Status   *uint8  `json:"status"`
}

type MenuSaveInput struct {
	MenuName string `json:"menu_name"`
	Sort     *int8  `json:"sort"`
	IsShow   *uint8 `json:"is_show"`
}

type StaffSaveInput struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Status   *int8  `json:"status"`
	IsOpen   *int8  `json:"is_open"`
	IsVerify *int8  `json:"is_verify"`
	IsGoods  *int8  `json:"is_goods"` // 1=可发货
}

func (s *Service) ListPlatformAdmins(ctx context.Context, page, limit int) (*PageResult[SystemAdmin], error) {
	page, limit = normalizePage(page, limit)
	list, total, err := s.store.ListPlatformAdmins(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[SystemAdmin]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) CreatePlatformAdmin(ctx context.Context, in AdminSaveInput) (*SystemAdmin, error) {
	account := strings.TrimSpace(in.Account)
	if account == "" {
		return nil, ErrBadParam
	}
	pwd := strings.TrimSpace(in.Password)
	if len(pwd) < 6 {
		return nil, ErrWeakPassword
	}
	if _, err := s.store.FindPlatformAdminByAccount(ctx, account); err == nil {
		return nil, ErrAccountExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	status := uint8(1)
	if in.Status != nil {
		status = *in.Status
	}
	roles := strings.TrimSpace(in.Roles)
	if roles == "" {
		roles = "1"
	}
	row := &SystemAdmin{
		Account: account, Pwd: string(hash),
		RealName: strings.TrimSpace(in.RealName), Phone: strings.TrimSpace(in.Phone),
		Roles: roles, Status: status, Level: 1,
	}
	if err := s.store.CreatePlatformAdmin(ctx, row); err != nil {
		return nil, err
	}
	return s.store.FindPlatformAdminByID(ctx, row.AdminID)
}

func (s *Service) UpdatePlatformAdmin(ctx context.Context, id, actorID uint, in AdminSaveInput) (*SystemAdmin, error) {
	row, err := s.store.FindPlatformAdminByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if name := strings.TrimSpace(in.RealName); name != "" {
		row.RealName = name
	}
	if phone := strings.TrimSpace(in.Phone); phone != "" {
		row.Phone = phone
	}
	if roles := strings.TrimSpace(in.Roles); roles != "" {
		row.Roles = roles
	}
	if in.Status != nil {
		if id == actorID && *in.Status != 1 {
			return nil, ErrBadParam
		}
		row.Status = *in.Status
	}
	if pwd := strings.TrimSpace(in.Password); pwd != "" {
		if len(pwd) < 6 {
			return nil, ErrWeakPassword
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		row.Pwd = string(hash)
	}
	if err := s.store.UpdatePlatformAdmin(ctx, row); err != nil {
		return nil, err
	}
	return s.store.FindPlatformAdminByID(ctx, id)
}

func (s *Service) ListMerchantAdmins(ctx context.Context, merID uint, page, limit int) (*PageResult[MerchantAdmin], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalizePage(page, limit)
	list, total, err := s.store.ListMerchantAdmins(ctx, merID, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[MerchantAdmin]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) CreateMerchantAdmin(ctx context.Context, merID uint, in AdminSaveInput) (*MerchantAdmin, error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	account := strings.TrimSpace(in.Account)
	if account == "" {
		return nil, ErrBadParam
	}
	pwd := strings.TrimSpace(in.Password)
	if len(pwd) < 6 {
		return nil, ErrWeakPassword
	}
	if _, err := s.store.FindMerchantAdminByAccount(ctx, account); err == nil {
		return nil, ErrAccountExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	roles := strings.TrimSpace(in.Roles)
	if roles == "" {
		roles = "2"
	}
	if err := s.validateMerchantRoleIDs(ctx, merID, roles); err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	status := uint8(1)
	if in.Status != nil {
		status = *in.Status
	}
	row := &MerchantAdmin{
		MerID: merID, Account: account, Pwd: string(hash),
		RealName: strings.TrimSpace(in.RealName), Phone: strings.TrimSpace(in.Phone),
		Roles: roles, Status: status, Level: 1,
	}
	if err := s.store.CreateMerchantAdmin(ctx, row); err != nil {
		return nil, err
	}
	return s.store.FindMerchantAdminByID(ctx, row.MerchantAdminID)
}

func (s *Service) UpdateMerchantAdmin(ctx context.Context, merID, id, actorID uint, in AdminSaveInput) (*MerchantAdmin, error) {
	if merID == 0 || id == 0 {
		return nil, ErrBadParam
	}
	row, err := s.store.FindMerchantAdminByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if row.MerID != merID {
		return nil, ErrNotFound
	}
	if name := strings.TrimSpace(in.RealName); name != "" {
		row.RealName = name
	}
	if phone := strings.TrimSpace(in.Phone); phone != "" {
		row.Phone = phone
	}
	if roles := strings.TrimSpace(in.Roles); roles != "" {
		if err := s.validateMerchantRoleIDs(ctx, merID, roles); err != nil {
			return nil, err
		}
		row.Roles = roles
	}
	if in.Status != nil {
		if id == actorID && *in.Status != 1 {
			return nil, ErrBadParam
		}
		// 主账号不可被禁用
		if row.Level == 0 && *in.Status != 1 {
			return nil, ErrBadParam
		}
		row.Status = *in.Status
	}
	if pwd := strings.TrimSpace(in.Password); pwd != "" {
		if len(pwd) < 6 {
			return nil, ErrWeakPassword
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		row.Pwd = string(hash)
	}
	if err := s.store.UpdateMerchantAdmin(ctx, row); err != nil {
		return nil, err
	}
	return s.store.FindMerchantAdminByID(ctx, id)
}

func (s *Service) validateMerchantRoleIDs(ctx context.Context, merID uint, rolesCSV string) error {
	ids := parseIDs(rolesCSV)
	if len(ids) == 0 {
		return ErrBadParam
	}
	for _, id := range ids {
		role, err := s.store.GetRole(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBadParam
			}
			return err
		}
		if role.Status != 1 {
			return ErrBadParam
		}
		ok := (role.MerID == 0 && role.IsAgent == 2) || role.MerID == merID
		if !ok {
			return ErrBadParam
		}
	}
	return nil
}

func (s *Service) ListRoles(ctx context.Context, merID uint, page, limit int) (*PageResult[SystemRole], error) {
	page, limit = normalizePage(page, limit)
	list, total, err := s.store.ListRoles(ctx, merID, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[SystemRole]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListMerchantRoles(ctx context.Context, merID uint, page, limit int) (*PageResult[SystemRole], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalizePage(page, limit)
	list, total, err := s.store.ListMerchantRoles(ctx, merID, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[SystemRole]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) CreateRole(ctx context.Context, merID uint, isAgent int8, in RoleSaveInput) (*SystemRole, error) {
	name := strings.TrimSpace(in.RoleName)
	if name == "" {
		return nil, ErrBadParam
	}
	rules := strings.TrimSpace(in.Rules)
	if in.MenuIDs != nil {
		rules = joinMenuIDs(*in.MenuIDs)
	}
	status := uint8(1)
	if in.Status != nil {
		status = *in.Status
	}
	row := &SystemRole{
		RoleName: name, Rules: rules, Status: status,
		MerID: merID, IsAgent: isAgent,
	}
	if err := s.store.CreateRole(ctx, row); err != nil {
		return nil, err
	}
	return s.store.GetRole(ctx, row.RoleID)
}

func (s *Service) UpdateRole(ctx context.Context, id uint, ownerMerID *uint, in RoleSaveInput) (*SystemRole, error) {
	row, err := s.store.GetRole(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if ownerMerID != nil {
		// 商户只能改本店角色，不可改共享模板（mer_id=0）
		if row.MerID == 0 || row.MerID != *ownerMerID {
			return nil, ErrNotFound
		}
	}
	if name := strings.TrimSpace(in.RoleName); name != "" {
		row.RoleName = name
	}
	if in.MenuIDs != nil {
		row.Rules = joinMenuIDs(*in.MenuIDs)
	} else if in.Rules != "" {
		row.Rules = strings.TrimSpace(in.Rules)
	}
	if in.Status != nil {
		row.Status = *in.Status
	}
	if err := s.store.UpdateRole(ctx, row); err != nil {
		return nil, err
	}
	return s.store.GetRole(ctx, id)
}

func (s *Service) ListMenusManage(ctx context.Context, isMer uint8) ([]SystemMenu, error) {
	return s.store.ListMenusManage(ctx, isMer)
}

func (s *Service) MenuTree(ctx context.Context, isMer uint8) ([]*MenuNode, error) {
	rows, err := s.store.ListMenusManage(ctx, isMer)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(rows), nil
}

func joinMenuIDs(ids []uint) string {
	if len(ids) == 0 {
		return ""
	}
	parts := make([]string, 0, len(ids))
	seen := map[uint]struct{}{}
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		parts = append(parts, strconv.FormatUint(uint64(id), 10))
	}
	return strings.Join(parts, ",")
}

func (s *Service) UpdateMenu(ctx context.Context, id uint, in MenuSaveInput) (*SystemMenu, error) {
	row, err := s.store.GetMenu(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if name := strings.TrimSpace(in.MenuName); name != "" {
		row.MenuName = name
	}
	if in.Sort != nil {
		row.Sort = *in.Sort
	}
	if in.IsShow != nil {
		row.IsShow = *in.IsShow
	}
	if err := s.store.UpdateMenu(ctx, row); err != nil {
		return nil, err
	}
	return s.store.GetMenu(ctx, id)
}

func (s *Service) ListStaff(ctx context.Context, merID uint, page, limit int) (*PageResult[StoreService], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalizePage(page, limit)
	list, total, err := s.store.ListStoreServices(ctx, merID, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[StoreService]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) CreateStaff(ctx context.Context, merID uint, in StaffSaveInput) (*StoreService, error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	account := strings.TrimSpace(in.Account)
	if account == "" {
		return nil, ErrBadParam
	}
	pwd := strings.TrimSpace(in.Password)
	if len(pwd) < 6 {
		return nil, ErrWeakPassword
	}
	if _, err := s.store.FindStoreServiceByAccount(ctx, account); err == nil {
		return nil, ErrAccountExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	nick := strings.TrimSpace(in.Nickname)
	if nick == "" {
		nick = account
	}
	row := &StoreService{
		MerID: merID, Account: account, Pwd: string(hash), Nickname: nick,
		Phone: strings.TrimSpace(in.Phone), Status: 1, IsOpen: 1, IsVerify: 1, IsGoods: 1,
		Customer: 0, CreateTime: time.Now(),
	}
	if in.Status != nil {
		row.Status = *in.Status
	}
	if in.IsOpen != nil {
		row.IsOpen = *in.IsOpen
	}
	if in.IsVerify != nil {
		row.IsVerify = *in.IsVerify
	}
	if in.IsGoods != nil {
		row.IsGoods = *in.IsGoods
	}
	if err := s.store.CreateStoreService(ctx, row); err != nil {
		return nil, err
	}
	return s.store.FindStoreServiceByID(ctx, row.ServiceID)
}

func (s *Service) UpdateStaff(ctx context.Context, merID, id uint, in StaffSaveInput) (*StoreService, error) {
	row, err := s.store.FindStoreServiceByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if row.MerID != merID {
		return nil, ErrNotFound
	}
	if nick := strings.TrimSpace(in.Nickname); nick != "" {
		row.Nickname = nick
	}
	if phone := strings.TrimSpace(in.Phone); phone != "" {
		row.Phone = phone
	}
	if in.Status != nil {
		row.Status = *in.Status
	}
	if in.IsOpen != nil {
		row.IsOpen = *in.IsOpen
	}
	if in.IsVerify != nil {
		row.IsVerify = *in.IsVerify
	}
	if in.IsGoods != nil {
		row.IsGoods = *in.IsGoods
	}
	if pwd := strings.TrimSpace(in.Password); pwd != "" {
		if len(pwd) < 6 {
			return nil, ErrWeakPassword
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		row.Pwd = string(hash)
	}
	if err := s.store.UpdateStoreService(ctx, row); err != nil {
		return nil, err
	}
	return s.store.FindStoreServiceByID(ctx, id)
}

func normalizePage(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}
