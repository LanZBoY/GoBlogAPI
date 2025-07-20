package basemodel

import "testing"

func TestNewDefaultQuery(t *testing.T) {
	q := NewDefaultQuery()
	if q.Skip != 0 || q.Limit != 10 {
		t.Fatalf("unexpected defaults: %+v", q)
	}
}
