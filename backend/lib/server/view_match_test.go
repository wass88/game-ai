package server

import "testing"

func TestViewMatches(t *testing.T) {
	db := mockPlayoutDB()
	var of int = 0;
	res, err := db.GetMatches(1, &of)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Matches) == 0 {
		t.Fatalf("No response")
	}
}

func TestViewMatch(t *testing.T) {
	db := mockPlayoutDB()
	res, err := db.GetMatch(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", res)
}
