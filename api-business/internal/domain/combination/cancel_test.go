package combination

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"
)

type memStore struct {
	members map[uint]*Member
	buyings map[uint]*Buying
	groups  map[uint]*ProductGroup
	nextMID uint
}

func newMemStore() *memStore {
	return &memStore{
		members: map[uint]*Member{},
		buyings: map[uint]*Buying{},
		groups:  map[uint]*ProductGroup{},
		nextMID: 1,
	}
}

func (m *memStore) GetMemberByOrder(_ context.Context, orderID uint) (*Member, error) {
	for _, row := range m.members {
		if row.OrderID == orderID && row.IsDel == 0 {
			cp := *row
			return &cp, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *memStore) SoftDeleteMember(_ context.Context, id uint) error {
	row, ok := m.members[id]
	if !ok || row.IsDel == 1 {
		return gorm.ErrRecordNotFound
	}
	row.IsDel = 1
	return nil
}

func (m *memStore) ListMembers(_ context.Context, buyingID uint) ([]Member, error) {
	var out []Member
	for _, row := range m.members {
		if row.GroupBuyingID == buyingID && row.IsDel == 0 {
			out = append(out, *row)
		}
	}
	return out, nil
}

func (m *memStore) GetBuying(_ context.Context, id uint) (*Buying, error) {
	row, ok := m.buyings[id]
	if !ok || row.IsDel == 1 {
		return nil, gorm.ErrRecordNotFound
	}
	cp := *row
	return &cp, nil
}

func (m *memStore) SoftDeleteBuying(_ context.Context, id uint) error {
	row, ok := m.buyings[id]
	if !ok || row.IsDel == 1 {
		return gorm.ErrRecordNotFound
	}
	row.IsDel = 1
	row.Status = -1
	return nil
}

// unused Store methods
func (m *memStore) ListGroups(context.Context, *uint, bool, int, int) ([]ProductGroup, int64, error) {
	return nil, 0, errors.New("n/a")
}
func (m *memStore) GetGroup(_ context.Context, id uint) (*ProductGroup, error) {
	g, ok := m.groups[id]
	if !ok || g.IsDel == 1 {
		return nil, gorm.ErrRecordNotFound
	}
	cp := *g
	return &cp, nil
}
func (m *memStore) CreateGroup(context.Context, *ProductGroup) error { return errors.New("n/a") }
func (m *memStore) UpdateGroup(context.Context, *ProductGroup) error { return errors.New("n/a") }
func (m *memStore) SoftDeleteGroup(context.Context, uint) error      { return errors.New("n/a") }
func (m *memStore) LoadProductMeta(context.Context, uint) (string, string, string, float64, float64, uint, error) {
	return "demo", "", "mer", 99, 10, 1, nil
}
func (m *memStore) CreateBuying(context.Context, *Buying) error { return errors.New("n/a") }
func (m *memStore) UpdateBuying(context.Context, *Buying) error { return errors.New("n/a") }
func (m *memStore) ListOpenBuyings(context.Context, uint, int) ([]Buying, error) {
	return nil, errors.New("n/a")
}
func (m *memStore) CreateMember(_ context.Context, row *Member) error {
	if row.ID == 0 {
		row.ID = m.nextMID
		m.nextMID++
	}
	cp := *row
	m.members[cp.ID] = &cp
	return nil
}
func (m *memStore) FindMember(_ context.Context, buyingID, uid uint) (*Member, error) {
	for _, row := range m.members {
		if row.GroupBuyingID == buyingID && row.UID == uid && row.IsDel == 0 {
			cp := *row
			return &cp, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *memStore) UpdateMember(context.Context, *Member) error { return errors.New("n/a") }
func (m *memStore) ListOrderIDsByBuying(context.Context, uint) ([]uint, error) {
	return nil, errors.New("n/a")
}
func (m *memStore) BumpGroupSuccess(context.Context, uint) error { return errors.New("n/a") }

func TestCancelUnpaid_LeaderClosesBuying(t *testing.T) {
	st := newMemStore()
	st.buyings[1] = &Buying{GroupBuyingID: 1, Status: 0}
	st.members[10] = &Member{ID: 10, GroupBuyingID: 1, OrderID: 100, Status: 0, IsLeader: 1}
	svc := NewService(st)
	if err := svc.CancelUnpaid(context.Background(), 100); err != nil {
		t.Fatal(err)
	}
	if st.members[10].IsDel != 1 {
		t.Fatalf("member not deleted")
	}
	if st.buyings[1].IsDel != 1 || st.buyings[1].Status != -1 {
		t.Fatalf("buying not closed: %+v", st.buyings[1])
	}
}

func TestCancelUnpaid_JoinerKeepsBuying(t *testing.T) {
	st := newMemStore()
	st.buyings[1] = &Buying{GroupBuyingID: 1, Status: 0}
	st.members[10] = &Member{ID: 10, GroupBuyingID: 1, OrderID: 100, Status: 0, IsLeader: 1}
	st.members[11] = &Member{ID: 11, GroupBuyingID: 1, OrderID: 101, Status: 0, IsLeader: 0}
	svc := NewService(st)
	if err := svc.CancelUnpaid(context.Background(), 101); err != nil {
		t.Fatal(err)
	}
	if st.members[11].IsDel != 1 {
		t.Fatalf("joiner not deleted")
	}
	if st.buyings[1].IsDel != 0 {
		t.Fatalf("buying should stay open")
	}
}

func TestCancelUnpaid_SkipPaidMember(t *testing.T) {
	st := newMemStore()
	st.buyings[1] = &Buying{GroupBuyingID: 1, Status: 0}
	st.members[10] = &Member{ID: 10, GroupBuyingID: 1, OrderID: 100, Status: 1, IsLeader: 1}
	svc := NewService(st)
	if err := svc.CancelUnpaid(context.Background(), 100); err != nil {
		t.Fatal(err)
	}
	if st.members[10].IsDel != 0 {
		t.Fatalf("paid member should not be deleted")
	}
}

func TestBeginJoin_UnpaidSeatsBlockAndRelease(t *testing.T) {
	st := newMemStore()
	st.groups[5] = &ProductGroup{
		ProductGroupID: 5, ProductID: 1, MerID: 1, Price: 9.9,
		BuyingCountNum: 2, BuyingNum: 1, Time: 24,
		Status: 1, IsShow: 1, ActionStatus: 1, ProductStatus: 1,
	}
	st.buyings[1] = &Buying{
		GroupBuyingID: 1, ProductGroupID: 5, Status: 0,
		BuyingCountNum: 2, YetBuyingNum: 0,
		EndTime: 1<<62,
	}
	st.members[10] = &Member{ID: 10, GroupBuyingID: 1, UID: 1, OrderID: 100, Status: 0, IsLeader: 1}
	st.members[11] = &Member{ID: 11, GroupBuyingID: 1, UID: 2, OrderID: 101, Status: 0, IsLeader: 0}
	st.nextMID = 12
	svc := NewService(st)
	ctx := context.Background()

	if _, _, err := svc.BeginJoin(ctx, 3, 5, 1, "u3"); !errors.Is(err, ErrBuyingFull) {
		t.Fatalf("want full, got %v", err)
	}
	if err := svc.CancelUnpaid(ctx, 101); err != nil {
		t.Fatal(err)
	}
	bid, leader, err := svc.BeginJoin(ctx, 3, 5, 1, "u3")
	if err != nil || leader || bid != 1 {
		t.Fatalf("rejoin after release: bid=%d leader=%v err=%v", bid, leader, err)
	}
	if err := svc.AttachMember(ctx, 1, 5, 3, 102, false, "u3"); err != nil {
		t.Fatal(err)
	}
	if err := svc.AttachMember(ctx, 1, 5, 4, 103, false, "u4"); !errors.Is(err, ErrBuyingFull) {
		t.Fatalf("attach over capacity: %v", err)
	}
}
