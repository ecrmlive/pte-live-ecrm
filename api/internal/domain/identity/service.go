package identity

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Store interface {
	FindPlatformAdminByAccount(ctx context.Context, account string) (*SystemAdmin, error)
	FindPlatformAdminByID(ctx context.Context, id uint) (*SystemAdmin, error)
	TouchPlatformLogin(ctx context.Context, id uint, ip string) error
	UpdatePlatformPassword(ctx context.Context, id uint, hash string) error

	FindMerchantAdminByAccount(ctx context.Context, account string) (*MerchantAdmin, error)
	FindMerchantAdminByID(ctx context.Context, id uint) (*MerchantAdmin, error)
	FindMerchant(ctx context.Context, merID uint) (*Merchant, error)
	TouchMerchantLogin(ctx context.Context, id uint, ip string) error
	UpdateMerchantPassword(ctx context.Context, id uint, hash string) error

	FindUserByAccount(ctx context.Context, account string) (*User, error)
	FindUserByID(ctx context.Context, id uint) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	TouchUserLogin(ctx context.Context, id uint, ip string) error
	ListUsers(ctx context.Context, page, limit int) ([]User, int64, error)
	UpdateUserSvip(ctx context.Context, u *User) error

	FindStoreServiceByAccount(ctx context.Context, account string) (*StoreService, error)
	FindStoreServiceByID(ctx context.Context, id uint) (*StoreService, error)
	ListStoreServices(ctx context.Context, merID uint, page, limit int) ([]StoreService, int64, error)
	CreateStoreService(ctx context.Context, row *StoreService) error
	UpdateStoreService(ctx context.Context, row *StoreService) error

	ListMenusByIDs(ctx context.Context, isMer uint8, ids []uint) ([]SystemMenu, error)
	ListButtonMenusByIDs(ctx context.Context, isMer uint8, ids []uint) ([]SystemMenu, error)
	ListMenusManage(ctx context.Context, isMer uint8) ([]SystemMenu, error)
	GetMenu(ctx context.Context, id uint) (*SystemMenu, error)
	UpdateMenu(ctx context.Context, row *SystemMenu) error
	ListRoleRules(ctx context.Context, roleIDs []uint) ([]string, error)
	ListRoles(ctx context.Context, merID uint, page, limit int) ([]SystemRole, int64, error)
	ListMerchantRoles(ctx context.Context, merID uint, page, limit int) ([]SystemRole, int64, error)
	GetRole(ctx context.Context, id uint) (*SystemRole, error)
	CreateRole(ctx context.Context, row *SystemRole) error
	UpdateRole(ctx context.Context, row *SystemRole) error

	ListPlatformAdmins(ctx context.Context, page, limit int) ([]SystemAdmin, int64, error)
	CreatePlatformAdmin(ctx context.Context, row *SystemAdmin) error
	UpdatePlatformAdmin(ctx context.Context, row *SystemAdmin) error

	ListMerchantAdmins(ctx context.Context, merID uint, page, limit int) ([]MerchantAdmin, int64, error)
	CreateMerchantAdmin(ctx context.Context, row *MerchantAdmin) error
	UpdateMerchantAdmin(ctx context.Context, row *MerchantAdmin) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

type PlatformProfile struct {
	AdminID  uint   `json:"admin_id"`
	Account  string `json:"account"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Roles    string `json:"roles"`
	Level    uint8  `json:"level"`
}

type MerchantProfile struct {
	MerchantAdminID uint   `json:"merchant_admin_id"`
	MerID           uint   `json:"mer_id"`
	MerName         string `json:"mer_name"`
	Account         string `json:"account"`
	RealName        string `json:"real_name"`
	Phone           string `json:"phone"`
	Roles           string `json:"roles"`
	Level           uint8  `json:"level"`
}

type AppProfile struct {
	UID         uint       `json:"uid"`
	Account     string     `json:"account"`
	Nickname    string     `json:"nickname"`
	Avatar      string     `json:"avatar"`
	Phone       string     `json:"phone"`
	Integral    int        `json:"integral"`
	NowMoney    float64    `json:"now_money"`
	IsSvip      int8       `json:"is_svip"`
	SvipEndtime *time.Time `json:"svip_endtime,omitempty"`
	IsSvipActive bool      `json:"is_svip_active"`
}

type StoreServiceProfile struct {
	ServiceID uint   `json:"service_id"`
	MerID     uint   `json:"mer_id"`
	MerName   string `json:"mer_name"`
	Account   string `json:"account"`
	Nickname  string `json:"nickname"`
	IsVerify  int8   `json:"is_verify"`
	IsGoods   int8   `json:"is_goods"`  // 1=可发货
	Customer  int8   `json:"customer"` // 1=可客服
}

func (s *Service) LoginPlatform(ctx context.Context, account, password, ip string) (*SystemAdmin, error) {
	admin, err := s.store.FindPlatformAdminByAccount(ctx, strings.TrimSpace(account))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if admin.Status != 1 {
		return nil, ErrAccountDisabled
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Pwd), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	_ = s.store.TouchPlatformLogin(ctx, admin.AdminID, ip)
	return admin, nil
}

func (s *Service) LoginMerchant(ctx context.Context, account, password, ip string) (*MerchantAdmin, *Merchant, error) {
	admin, err := s.store.FindMerchantAdminByAccount(ctx, strings.TrimSpace(account))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, err
	}
	if admin.Status != 1 {
		return nil, nil, ErrAccountDisabled
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Pwd), []byte(password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}
	mer, err := s.store.FindMerchant(ctx, admin.MerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrMerchantDisabled
		}
		return nil, nil, err
	}
	if mer.Status != 1 || mer.MerState != 1 {
		return nil, nil, ErrMerchantDisabled
	}
	_ = s.store.TouchMerchantLogin(ctx, admin.MerchantAdminID, ip)
	return admin, mer, nil
}

func (s *Service) PlatformProfile(ctx context.Context, adminID uint) (*PlatformProfile, error) {
	admin, err := s.store.FindPlatformAdminByID(ctx, adminID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &PlatformProfile{
		AdminID:  admin.AdminID,
		Account:  admin.Account,
		RealName: admin.RealName,
		Phone:    admin.Phone,
		Roles:    admin.Roles,
		Level:    admin.Level,
	}, nil
}

func (s *Service) MerchantProfile(ctx context.Context, adminID uint) (*MerchantProfile, error) {
	admin, err := s.store.FindMerchantAdminByID(ctx, adminID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	mer, err := s.store.FindMerchant(ctx, admin.MerID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	merName := ""
	if mer != nil {
		merName = mer.MerName
	}
	return &MerchantProfile{
		MerchantAdminID: admin.MerchantAdminID,
		MerID:           admin.MerID,
		MerName:         merName,
		Account:         admin.Account,
		RealName:        admin.RealName,
		Phone:           admin.Phone,
		Roles:           admin.Roles,
		Level:           admin.Level,
	}, nil
}

func (s *Service) ChangePlatformPassword(ctx context.Context, adminID uint, oldPwd, newPwd string) error {
	if len(newPwd) < 6 {
		return ErrWeakPassword
	}
	admin, err := s.store.FindPlatformAdminByID(ctx, adminID)
	if err != nil {
		return ErrNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Pwd), []byte(oldPwd)); err != nil {
		return ErrBadPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.store.UpdatePlatformPassword(ctx, adminID, string(hash))
}

func (s *Service) ChangeMerchantPassword(ctx context.Context, adminID uint, oldPwd, newPwd string) error {
	if len(newPwd) < 6 {
		return ErrWeakPassword
	}
	admin, err := s.store.FindMerchantAdminByID(ctx, adminID)
	if err != nil {
		return ErrNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Pwd), []byte(oldPwd)); err != nil {
		return ErrBadPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.store.UpdateMerchantPassword(ctx, adminID, string(hash))
}

func (s *Service) LoginApp(ctx context.Context, account, password, ip string) (*User, error) {
	user, err := s.store.FindUserByAccount(ctx, strings.TrimSpace(account))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if user.Status != 1 {
		return nil, ErrAccountDisabled
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Pwd), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	_ = s.store.TouchUserLogin(ctx, user.UID, ip)
	return user, nil
}

func (s *Service) RegisterApp(ctx context.Context, account, password, nickname, ip string) (*User, error) {
	account = strings.TrimSpace(account)
	nickname = strings.TrimSpace(nickname)
	if account == "" || len(password) < 6 {
		return nil, ErrWeakPassword
	}
	if nickname == "" {
		nickname = account
	}
	if len([]rune(nickname)) > 16 {
		nickname = string([]rune(nickname)[:16])
	}
	if _, err := s.store.FindUserByAccount(ctx, account); err == nil {
		return nil, ErrAccountExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &User{
		Account:  account,
		Pwd:      string(hash),
		Nickname: nickname,
		Avatar:   "",
		Status:   1,
		UserType: "h5",
		LastIP:   ip,
	}
	if err := s.store.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	_ = s.store.TouchUserLogin(ctx, user.UID, ip)
	return user, nil
}

func (s *Service) AppProfile(ctx context.Context, uid uint) (*AppProfile, error) {
	user, err := s.store.FindUserByID(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &AppProfile{
		UID:          user.UID,
		Account:      user.Account,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		Phone:        user.Phone,
		Integral:     user.Integral,
		NowMoney:     user.NowMoney,
		IsSvip:       user.IsSvip,
		SvipEndtime:  user.SvipEndtime,
		IsSvipActive: UserSvipActive(user),
	}, nil
}

func (s *Service) LoginStoreService(ctx context.Context, account, password string) (*StoreService, *Merchant, error) {
	row, err := s.store.FindStoreServiceByAccount(ctx, strings.TrimSpace(account))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, err
	}
	if row.IsDel == 1 || row.Status != 1 || row.IsOpen != 1 {
		return nil, nil, ErrAccountDisabled
	}
	if err := bcrypt.CompareHashAndPassword([]byte(row.Pwd), []byte(password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}
	mer, err := s.store.FindMerchant(ctx, row.MerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrMerchantDisabled
		}
		return nil, nil, err
	}
	if mer.Status != 1 || mer.MerState != 1 {
		return nil, nil, ErrMerchantDisabled
	}
	return row, mer, nil
}

// LoginCustomerService 客服门户：店员账号且 customer=1。
func (s *Service) LoginCustomerService(ctx context.Context, account, password string) (*StoreService, *Merchant, error) {
	row, mer, err := s.LoginStoreService(ctx, account, password)
	if err != nil {
		return nil, nil, err
	}
	if row.Customer != 1 {
		return nil, nil, ErrNoCustomerPerm
	}
	return row, mer, nil
}

func (s *Service) StoreServiceProfile(ctx context.Context, serviceID uint) (*StoreServiceProfile, error) {
	row, err := s.store.FindStoreServiceByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if row.IsDel == 1 || row.Status != 1 {
		return nil, ErrAccountDisabled
	}
	merName := ""
	if mer, err := s.store.FindMerchant(ctx, row.MerID); err == nil && mer != nil {
		merName = mer.MerName
	}
	return &StoreServiceProfile{
		ServiceID: row.ServiceID,
		MerID:     row.MerID,
		MerName:   merName,
		Account:   row.Account,
		Nickname:  row.Nickname,
		IsVerify:  row.IsVerify,
		IsGoods:   row.IsGoods,
		Customer:  row.Customer,
	}, nil
}

func (s *Service) RequireStoreVerifyPerm(ctx context.Context, serviceID uint) error {
	profile, err := s.StoreServiceProfile(ctx, serviceID)
	if err != nil {
		return err
	}
	if profile.IsVerify != 1 {
		return ErrNoVerifyPerm
	}
	return nil
}

// RequireStoreDeliverPerm 店员发货权（is_goods=1）。
func (s *Service) RequireStoreDeliverPerm(ctx context.Context, serviceID uint) error {
	profile, err := s.StoreServiceProfile(ctx, serviceID)
	if err != nil {
		return err
	}
	if profile.IsGoods != 1 {
		return ErrNoDeliverPerm
	}
	return nil
}

func (s *Service) MenusForPlatform(ctx context.Context, rolesCSV string, level uint8) ([]*MenuNode, error) {
	return s.menus(ctx, 1, rolesCSV, level == 0)
}

func (s *Service) MenusForMerchant(ctx context.Context, rolesCSV string, level uint8) ([]*MenuNode, error) {
	return s.menus(ctx, 2, rolesCSV, level == 0)
}

func (s *Service) menus(ctx context.Context, isMer uint8, rolesCSV string, isSuper bool) ([]*MenuNode, error) {
	var menuIDs []uint
	if !isSuper {
		roleIDs := parseIDs(rolesCSV)
		rules, err := s.store.ListRoleRules(ctx, roleIDs)
		if err != nil {
			return nil, err
		}
		menuIDs = parseIDs(strings.Join(rules, ","))
		if len(menuIDs) == 0 {
			return []*MenuNode{}, nil
		}
	}
	rows, err := s.store.ListMenusByIDs(ctx, isMer, menuIDs)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(rows), nil
}

func parseIDs(csv string) []uint {
	parts := strings.Split(csv, ",")
	out := make([]uint, 0, len(parts))
	seen := map[uint]struct{}{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.ParseUint(p, 10, 64)
		if err != nil {
			continue
		}
		id := uint(n)
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func buildMenuTree(rows []SystemMenu) []*MenuNode {
	nodes := make(map[uint]*MenuNode, len(rows))
	order := make([]uint, 0, len(rows))
	for _, row := range rows {
		nodes[row.MenuID] = &MenuNode{
			MenuID:   row.MenuID,
			PID:      row.PID,
			Path:     row.Path,
			Icon:     row.Icon,
			MenuName: row.MenuName,
			Route:    row.Route,
			Sort:     row.Sort,
			IsMenu:   row.IsMenu,
		}
		order = append(order, row.MenuID)
	}
	roots := make([]*MenuNode, 0)
	for _, id := range order {
		node := nodes[id]
		if node.PID == 0 {
			roots = append(roots, node)
			continue
		}
		parent, ok := nodes[node.PID]
		if !ok {
			roots = append(roots, node)
			continue
		}
		parent.Children = append(parent.Children, node)
	}
	var sortNodes func(list []*MenuNode)
	sortNodes = func(list []*MenuNode) {
		sort.SliceStable(list, func(i, j int) bool {
			if list[i].Sort == list[j].Sort {
				return list[i].MenuID < list[j].MenuID
			}
			return list[i].Sort > list[j].Sort
		})
		for _, n := range list {
			if len(n.Children) > 0 {
				sortNodes(n.Children)
			}
		}
	}
	sortNodes(roots)
	return roots
}
