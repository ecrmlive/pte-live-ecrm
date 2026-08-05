package dashboard

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestAppendUnique(t *testing.T) {
	values := appendUnique([]uint64{2, 5}, 5)
	values = appendUnique(values, 8)
	values = appendUnique(values, 0)
	if len(values) != 3 || values[0] != 2 || values[1] != 5 || values[2] != 8 {
		t.Fatalf("unexpected store scope IDs: %#v", values)
	}
}

func TestHasRole(t *testing.T) {
	if !hasRole([]string{"merchant", "region"}, "region") {
		t.Fatal("region role should be recognized")
	}
	if hasRole([]string{"operations"}, "customer_service") {
		t.Fatal("unrelated role must not pass service scope")
	}
}

func TestSummaryEmptyRankingSerializesAsArray(t *testing.T) {
	row := Summary{StoreSalesRank: []StoreSalesRank{}}
	raw, err := json.Marshal(row)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(raw, []byte(`"store_sales_rank":[]`)) {
		t.Fatalf("empty ranking must be an array, got %s", raw)
	}
}
