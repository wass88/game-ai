package server

import "testing"

func TestNewUser(t *testing.T) {
	db := getDB()
	dropMock(db)
	db.DB.MustExec(`DELETE FROM user WHERE name = ?`, "wass")
	id, err := db.NewUser("wass")
	if err != nil {
		t.Fatal(err)
	}
	db.DB.MustExec(`DELETE FROM user WHERE id = ?`, id)
}

func TestNewUserIfNotExist(t *testing.T) {
	db := getDB()
	dropMock(db)
	db.DB.MustExec(`DELETE FROM user WHERE name = ?`, "username")
	id1, err := db.NewUserIfNotExist("username")
	if err != nil {
		t.Fatal(err)
	}
	id2, err := db.NewUserIfNotExist("username")
	if err != nil {
		t.Fatal(err)
	}
	if id1 == nil || id2 == nil || id1.ID != id2.ID {
		t.Fatalf("id1 != id2\nid1: %v\nid2: %v", id1, id2)
	}
	db.DB.MustExec(`DELETE FROM user WHERE id = ?`, id1.ID)
}
