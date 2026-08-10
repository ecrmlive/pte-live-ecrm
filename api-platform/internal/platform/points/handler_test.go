package points

import "testing"

func TestValidUpdate(t *testing.T) {
	points, zeroPoints, stock, sale := int64(99), int64(0), 0, 1
	if !validUpdate(&productSaveInput{PointsRequired: &points, Version: 1}) ||
		!validUpdate(&productSaveInput{Stock: &stock, Version: 1}) ||
		!validUpdate(&productSaveInput{SaleStatus: &sale, Version: 1}) {
		t.Fatal("valid product update rejected")
	}
	if validUpdate(&productSaveInput{}) || validUpdate(&productSaveInput{PointsRequired: &zeroPoints, Version: 1}) {
		t.Fatal("invalid product update accepted")
	}
}

func TestValidCreate(t *testing.T) {
	points, stock := int64(120), 10
	if !validCreate(&productSaveInput{Title: "真丝方巾", PointsRequired: &points, Stock: &stock}) {
		t.Fatal("valid create rejected")
	}
	if validCreate(&productSaveInput{Title: "", PointsRequired: &points, Stock: &stock}) {
		t.Fatal("empty title accepted")
	}
}

func TestOptionalPositive(t *testing.T) {
	if id, err := optionalPositive("12"); err != nil || id != 12 {
		t.Fatalf("valid id = %d, %v", id, err)
	}
	if _, err := optionalPositive("0"); err == nil {
		t.Fatal("zero id must be rejected")
	}
}
