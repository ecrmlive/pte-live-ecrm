package invoice

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"
)

type fakeStore struct {
	profiles map[uint64]InvoiceProfile
	invoices map[uint64]Invoice
	orders   map[uint64]OrderMeta
	nextID   uint64
}

func (f *fakeStore) ListProfiles(_ context.Context, userID uint64) ([]InvoiceProfile, error) {
	rows := []InvoiceProfile{}
	for _, row := range f.profiles {
		if row.UserID == userID {
			rows = append(rows, row)
		}
	}
	return rows, nil
}
func (f *fakeStore) GetProfile(_ context.Context, userID, id uint64) (*InvoiceProfile, error) {
	row, ok := f.profiles[id]
	if !ok || row.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}
func (f *fakeStore) CreateProfile(_ context.Context, row *InvoiceProfile) error {
	f.nextID++
	row.ID = f.nextID
	f.profiles[row.ID] = *row
	return nil
}
func (f *fakeStore) UpdateProfile(_ context.Context, row *InvoiceProfile) error {
	if _, ok := f.profiles[row.ID]; !ok {
		return gorm.ErrRecordNotFound
	}
	f.profiles[row.ID] = *row
	return nil
}
func (f *fakeStore) DeleteProfile(_ context.Context, userID, id uint64) error {
	row, ok := f.profiles[id]
	if !ok || row.UserID != userID {
		return gorm.ErrRecordNotFound
	}
	delete(f.profiles, id)
	return nil
}
func (f *fakeStore) SetDefaultProfile(_ context.Context, userID, id uint64) error {
	row, ok := f.profiles[id]
	if !ok || row.UserID != userID {
		return gorm.ErrRecordNotFound
	}
	for key, item := range f.profiles {
		if item.UserID == userID {
			item.IsDefault = key == id
			f.profiles[key] = item
		}
	}
	return nil
}
func (f *fakeStore) ListByUID(_ context.Context, userID uint64, page, limit int) ([]Invoice, int64, error) {
	rows := []Invoice{}
	for _, row := range f.invoices {
		if order, ok := f.orders[row.OrderID]; ok && order.UserID == userID {
			rows = append(rows, row)
		}
	}
	return rows, int64(len(rows)), nil
}
func (f *fakeStore) GetByUID(_ context.Context, userID, id uint64) (*Invoice, error) {
	row, ok := f.invoices[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	order, ok := f.orders[row.OrderID]
	if !ok || order.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}
func (f *fakeStore) FindByOrder(_ context.Context, orderID uint64) (*Invoice, error) {
	for _, row := range f.invoices {
		if row.OrderID == orderID {
			return &row, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeStore) Create(_ context.Context, row *Invoice) error {
	f.nextID++
	row.ID = f.nextID
	f.invoices[row.ID] = *row
	return nil
}
func (f *fakeStore) LoadOrder(_ context.Context, orderID uint64) (*OrderMeta, error) {
	row, ok := f.orders[orderID]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		profiles: map[uint64]InvoiceProfile{9: {ID: 9, UserID: 101, Type: ProfileEnterprise, Title: "七禧云创科技有限公司", TaxNo: "91310000DEMO12345X", Email: "finance@example.invalid", IsDefault: true}},
		invoices: map[uint64]Invoice{},
		orders:   map[uint64]OrderMeta{501: {OrderID: 501, UserID: 101, Status: "paid"}, 502: {OrderID: 502, UserID: 202, Status: "paid"}, 503: {OrderID: 503, UserID: 101, Status: "pending_pay"}},
		nextID:   20,
	}
}

func TestApplyUsesOwnedPaidOrderAndSnapshotsChineseProfile(t *testing.T) {
	store := newFakeStore()
	svc := NewService(store)
	got, err := svc.Apply(context.Background(), 101, ApplyInput{OrderID: 501, InvoiceProfileID: 9})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if got.Status != StatusRequested || got.Title != "七禧云创科技有限公司" || got.TaxNo != "91310000DEMO12345X" || got.ID == 0 {
		t.Fatalf("unexpected invoice: %#v", got)
	}
	if store.profiles[9].Title != got.Title {
		t.Fatal("invoice must retain profile snapshot")
	}
}

func TestApplyRejectsForeignAndUnpaidOrders(t *testing.T) {
	svc := NewService(newFakeStore())
	for _, orderID := range []uint64{502, 503} {
		_, err := svc.Apply(context.Background(), 101, ApplyInput{OrderID: orderID, InvoiceProfileID: 9})
		if !errors.Is(err, ErrOrder) {
			t.Fatalf("order %d err=%v", orderID, err)
		}
	}
}

func TestCreateProfileRequiresEnterpriseTaxNumberAndSetsDefault(t *testing.T) {
	store := newFakeStore()
	svc := NewService(store)
	if _, err := svc.CreateProfile(context.Background(), 101, ProfileInput{Type: ProfileEnterprise, Title: "七禧演示企业"}); !errors.Is(err, ErrBadParam) {
		t.Fatalf("enterprise without tax number err=%v", err)
	}
	profile, err := svc.CreateProfile(context.Background(), 101, ProfileInput{Type: ProfilePersonal, Title: "七禧体验用户", IsDefault: true})
	if err != nil || !profile.IsDefault {
		t.Fatalf("profile=%#v err=%v", profile, err)
	}
	if store.profiles[9].IsDefault {
		t.Fatal("previous default profile must be cleared")
	}
}

func TestGetMineDoesNotLeakOtherUsersInvoice(t *testing.T) {
	store := newFakeStore()
	store.invoices[18] = Invoice{ID: 18, OrderID: 502, Status: StatusRequested}
	_, err := NewService(store).GetMine(context.Background(), 101, 18)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err=%v", err)
	}
}
