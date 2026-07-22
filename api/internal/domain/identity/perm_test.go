package identity

import (
	"context"
	"errors"
	"testing"
)

type permStore struct {
	roles   map[uint]string
	buttons []SystemMenu
	admins  map[uint]*MerchantAdmin
}

func (p *permStore) ListRoleRules(_ context.Context, roleIDs []uint) ([]string, error) {
	var out []string
	for _, id := range roleIDs {
		if r, ok := p.roles[id]; ok {
			out = append(out, r)
		}
	}
	return out, nil
}

func (p *permStore) ListButtonMenusByIDs(_ context.Context, _ uint8, ids []uint) ([]SystemMenu, error) {
	if len(ids) == 0 {
		return p.buttons, nil
	}
	want := map[uint]struct{}{}
	for _, id := range ids {
		want[id] = struct{}{}
	}
	var out []SystemMenu
	for _, b := range p.buttons {
		if _, ok := want[b.MenuID]; ok {
			out = append(out, b)
		}
	}
	return out, nil
}

func (p *permStore) FindMerchantAdminByID(_ context.Context, id uint) (*MerchantAdmin, error) {
	a, ok := p.admins[id]
	if !ok {
		return nil, errors.New("not found")
	}
	cp := *a
	return &cp, nil
}

func (p *permStore) FindMerchant(_ context.Context, merID uint) (*Merchant, error) {
	return &Merchant{MerID: merID, MerName: "demo", Status: 1, MerState: 1}, nil
}

func TestRequireMerchantMenu_SubAccount(t *testing.T) {
	st := &permStore{
		roles: map[uint]string{3: "101,105,106,127"},
		buttons: []SystemMenu{
			{MenuID: 127, Path: "order/deliver", IsMenu: 2, IsMer: 2},
			{MenuID: 128, Path: "order/verify", IsMenu: 2, IsMer: 2},
		},
		admins: map[uint]*MerchantAdmin{
			3: {MerchantAdminID: 3, MerID: 1, Roles: "3", Level: 1, Status: 1},
			1: {MerchantAdminID: 1, MerID: 1, Roles: "2", Level: 0, Status: 1},
		},
	}
	// stub remaining Store via panic wrapper — only used methods above
	svc := NewService(&permStub{permStore: st})
	ctx := context.Background()

	if err := svc.RequireMerchantMenu(ctx, 3, MerPermOrderDeliver); err != nil {
		t.Fatalf("deliver should pass: %v", err)
	}
	if err := svc.RequireMerchantMenu(ctx, 3, MerPermOrderVerify); !errors.Is(err, ErrNoPerm) {
		t.Fatalf("verify should deny: %v", err)
	}
	if err := svc.RequireMerchantMenu(ctx, 1, MerPermOrderVerify); err != nil {
		t.Fatalf("level0 should pass: %v", err)
	}
	paths, err := svc.MerchantPermissionPaths(ctx, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 1 || paths[0] != "order/deliver" {
		t.Fatalf("paths=%v", paths)
	}
}

// permStub embeds unused Store methods as no-ops for compile.
type permStub struct{ *permStore }

func (p *permStub) FindPlatformAdminByAccount(context.Context, string) (*SystemAdmin, error) {
	return nil, errors.New("n/a")
}
func (p *permStub) FindPlatformAdminByID(context.Context, uint) (*SystemAdmin, error) {
	return nil, errors.New("n/a")
}
func (p *permStub) TouchPlatformLogin(context.Context, uint, string) error { return errors.New("n/a") }
func (p *permStub) UpdatePlatformPassword(context.Context, uint, string) error {
	return errors.New("n/a")
}
func (p *permStub) FindMerchantAdminByAccount(context.Context, string) (*MerchantAdmin, error) {
	return nil, errors.New("n/a")
}
func (p *permStub) TouchMerchantLogin(context.Context, uint, string) error { return errors.New("n/a") }
func (p *permStub) UpdateMerchantPassword(context.Context, uint, string) error {
	return errors.New("n/a")
}
func (p *permStub) FindUserByAccount(context.Context, string) (*User, error) {
	return nil, errors.New("n/a")
}
func (p *permStub) FindUserByID(context.Context, uint) (*User, error) { return nil, errors.New("n/a") }
func (p *permStub) CreateUser(context.Context, *User) error           { return errors.New("n/a") }
func (p *permStub) TouchUserLogin(context.Context, uint, string) error { return errors.New("n/a") }
func (p *permStub) ListUsers(context.Context, int, int) ([]User, int64, error) {
	return nil, 0, errors.New("n/a")
}
func (p *permStub) UpdateUserSvip(context.Context, *User) error { return errors.New("n/a") }
func (p *permStub) FindStoreServiceByAccount(context.Context, string) (*StoreService, error) {
	return nil, errors.New("n/a")
}
func (p *permStub) FindStoreServiceByID(context.Context, uint) (*StoreService, error) {
	return nil, errors.New("n/a")
}
func (p *permStub) ListStoreServices(context.Context, uint, int, int) ([]StoreService, int64, error) {
	return nil, 0, errors.New("n/a")
}
func (p *permStub) CreateStoreService(context.Context, *StoreService) error {
	return errors.New("n/a")
}
func (p *permStub) UpdateStoreService(context.Context, *StoreService) error {
	return errors.New("n/a")
}
func (p *permStub) ListMenusByIDs(context.Context, uint8, []uint) ([]SystemMenu, error) {
	return nil, errors.New("n/a")
}
func (p *permStub) ListMenusManage(context.Context, uint8) ([]SystemMenu, error) {
	return nil, errors.New("n/a")
}
func (p *permStub) GetMenu(context.Context, uint) (*SystemMenu, error) { return nil, errors.New("n/a") }
func (p *permStub) UpdateMenu(context.Context, *SystemMenu) error      { return errors.New("n/a") }
func (p *permStub) ListRoles(context.Context, uint, int, int) ([]SystemRole, int64, error) {
	return nil, 0, errors.New("n/a")
}
func (p *permStub) ListMerchantRoles(context.Context, uint, int, int) ([]SystemRole, int64, error) {
	return nil, 0, errors.New("n/a")
}
func (p *permStub) GetRole(context.Context, uint) (*SystemRole, error) { return nil, errors.New("n/a") }
func (p *permStub) CreateRole(context.Context, *SystemRole) error      { return errors.New("n/a") }
func (p *permStub) UpdateRole(context.Context, *SystemRole) error      { return errors.New("n/a") }
func (p *permStub) ListPlatformAdmins(context.Context, int, int) ([]SystemAdmin, int64, error) {
	return nil, 0, errors.New("n/a")
}
func (p *permStub) CreatePlatformAdmin(context.Context, *SystemAdmin) error {
	return errors.New("n/a")
}
func (p *permStub) UpdatePlatformAdmin(context.Context, *SystemAdmin) error {
	return errors.New("n/a")
}
func (p *permStub) ListMerchantAdmins(context.Context, uint, int, int) ([]MerchantAdmin, int64, error) {
	return nil, 0, errors.New("n/a")
}
func (p *permStub) CreateMerchantAdmin(context.Context, *MerchantAdmin) error {
	return errors.New("n/a")
}
func (p *permStub) UpdateMerchantAdmin(context.Context, *MerchantAdmin) error {
	return errors.New("n/a")
}
