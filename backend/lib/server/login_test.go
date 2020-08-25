package server

import "testing"

func TestNewUser(t *testing.T) {
	db := getDB()
	dropMock(db)
	id, err := db.NewUser("wass")
	if err != nil {
		t.Fatal(err)
	}
	db.DB.MustExec(`DELETE FROM user WHERE id = ?`, id)
}

func TestNewUserIfN(t *testing.T) {
	db := getDB()
	dropMock(db)
	id1, err := db.NewUserIfNotExist("wass")
	if err != nil {
		t.Fatal(err)
	}
	id2, err := db.NewUserIfNotExist("wass")
	if err != nil {
		t.Fatal(err)
	}
	if id1 != id2 {
		t.Fatalf("id1:%d != id2:%d", id1, id2)
	}
	db.DB.MustExec(`DELETE FROM user WHERE id = ?`, id1)
}
