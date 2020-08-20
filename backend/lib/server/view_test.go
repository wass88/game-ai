package server

import "testing"

func TestViewMatches(t *testing.T) {
	db := mockPlayoutDB()
	res, err := db.GetMatches(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) == 0 {
		t.Fatalf("No response")
	}
	for _, r := range res {
		t.Logf("%+v", r)
	}
}
