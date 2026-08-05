package points

import "testing"

func TestValidUpdate(t *testing.T) {
	points, zeroPoints, stock, sale := int64(99), int64(0), 0, 1
	if !validUpdate(&productUpdateInput{PointsRequired: &points, Version: 1}) || !validUpdate(&productUpdateInput{Stock: &stock, Version: 1}) || !validUpdate(&productUpdateInput{SaleStatus: &sale, Version: 1}) {
		t.Fatal("valid product update rejected")
	}
	if validUpdate(&productUpdateInput{Version: 0}) || validUpdate(&productUpdateInput{PointsRequired: &zeroPoints, Version: 1}) {
		t.Fatal("invalid product update accepted")
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
