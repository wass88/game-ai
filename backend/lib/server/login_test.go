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

func TestNewUserIfN(t *testing.T) {
	db := getDB()
	dropMock(db)
	db.DB.MustExec(`DELETE FROM user WHERE name = ?`, "wass")
	id1, err := db.NewUserIfNotExist("wass")
	if err != nil {
		t.Fatal(err)
	}
	id2, err := db.NewUserIfNotExist("wass")
	if err != nil {
		t.Fatal(err)
	}
	if id1 == nil || id2 == nil || id1.ID != id2.ID {
		t.Fatalf("id1:%v != id2:%v", id1, id2)
	}
	db.DB.MustExec(`DELETE FROM user WHERE id = ?`, id1.ID)
}
