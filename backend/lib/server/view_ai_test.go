package server

import "testing"

func TestViewAI(t *testing.T) {
	db := mockPlayoutDB()
	res, err := db.GetAIGithubsByGame(1)
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

func TestViewLatestAI(t *testing.T) {
	db := mockPlayoutDB()
	res, err := db.GetLatestAIByGame(1)
	if err != nil {
		t.Fatal(err)
	}
	ok := false
	for _, ai := range res {
		if ai.ID == 2 {
			ok = true
		}
		t.Logf("%v", ai)
	}
	if !ok {
		t.Fatalf("Missing ID=2\n%v", res)
	}
}
