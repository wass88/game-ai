package server

import "testing"


func TestRandomPlayoutConfig(t *testing.T) {
	db := mockPlayoutDB()
	p := NewAutoPlayout(db)
	_, err := db.DB.Exec(`INSERT ai (state, ai_github_id, commit)
		VALUES ("ready", 1, "000003")`)
	if err != nil {
		t.Fatal(err)
	}
	c, err := p.RandomPlayoutConfig()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", c)
}

func TestNewPlayoutFromConfig(t *testing.T) {
	db := mockPlayoutDB()
	p := NewAutoPlayout(db)
	c := &PlayoutConfig{1, []AIID{1, 2}}
	err := p.CreatePlayoutFromConfig(c)
	if err != nil {
		t.Fatal(err)
	}
}

func TestKickPlayout(t *testing.T) {
	db := mockPlayoutDB()
	p := NewAutoPlayout(db)
	_, err := db.DB.Exec(`INSERT ai (state, ai_github_id, commit)
		VALUES ("ready", 1, "000003")`)
	if err != nil {
		t.Fatal(err)
	}
	err = p.Kick()
	t.Logf("%v", err)
}