package integralcate

import "testing"

func TestNormalizeSave(t *testing.T) {
	show := 0
	sort := 12
	name, gotShow, gotSort, pic, pid, ok := normalizeSave(saveInput{
		CateName: "  美妆护肤  ",
		IsShow:   &show,
		Sort:     &sort,
		Pic:      "",
	})
	if !ok || name != "美妆护肤" || gotShow != 0 || gotSort != 12 || pic != "" || pid != 0 {
		t.Fatalf("normalizeSave got name=%q show=%d sort=%d pic=%q pid=%d ok=%v", name, gotShow, gotSort, pic, pid, ok)
	}
	if _, _, _, _, _, ok := normalizeSave(saveInput{CateName: "   "}); ok {
		t.Fatal("empty name should fail")
	}
	bad := 100000
	if _, _, _, _, _, ok := normalizeSave(saveInput{CateName: "数码", Sort: &bad}); ok {
		t.Fatal("sort overflow should fail")
	}
}
