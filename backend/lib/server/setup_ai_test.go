package server

import (
	"testing"
)

func mockCreateAIGithub(t *testing.T, db *DB) AIGithubID {
	aig := AIGithubA{
		UserID: 1,
		GameID: 1,
		Github: "wass88/reversi-random",
		Branch: "master",
	}
	id, err := db.CreateAIGithub(&aig)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func TestCreateAIGithub(t *testing.T) {
	db := mockGameUser()
	id := mockCreateAIGithub(t, db)
	if id < 0 {
		t.Fatalf("%d >= 0", id)
	}
	t.Logf("id=%d", id)
}

func TestGetAIGithub(t *testing.T) {
	db := mockGameUser()
	id := mockCreateAIGithub(t, db)

	list, err := db.GetActiveAI()
	if err != nil {
		t.Fatal(err)
	}
	ok := false
	for _, aia := range ([]AIGithubAct)(list) {
		if aia.ID == id {
			ok = true
		}
	}
	if !ok {
		t.Fatalf("Missing %d\nIn %v", id, list)
	}
}

func TestFetchCommit(t *testing.T) {
	commit, err := FetchCommitFromGithub("wass88/reversi-random", "master")
	if err != nil {
		t.Fatal(err)
	}
	if commit == "" {
		t.Fatal("Commit is empty")
	}
	t.Logf("Commit = %s", commit)
}

func TestFindAINeedUpdate(t *testing.T) {
	db := mockGameUser()
	_ = mockCreateAIGithub(t, db)
	res, err := db.FindAIGithubNeedUpdate()
	if err != nil {
		t.Fatal(err)
	}
	ok := false
	for _, u := range res {
		t.Logf("need Update %v", u)
		if u.Github == "wass88/reversi-random" {
			ok = true
		}
	}
	if !ok {
		t.Fatalf("Missing Update")
	}
}

func createAI(db *DB, t *testing.T) AIID {
	id := mockCreateAIGithub(t, db)
	u := AIGithubNeedUpdate{
		Github: "wass88/reversi-random",
		Branch: "master",
		Commit: "b3cd1a475dded156758005866761de51ee690607",
		ID:     id,
	}
	aiID, err := db.CreateAI(&u)
	if err != nil {
		t.Fatal(err)
	}
	if aiID < 0 {
		t.Fatalf("Missing id %d", aiID)
	}
	return aiID
}

func TestCreateAI(t *testing.T) {
	db := mockGameUser()
	aiID := createAI(db, t)
	t.Logf("id = %d", aiID)
}

func TestNeedSetupAI(t *testing.T) {
	db := mockGameUser()
	aiID := createAI(db, t)
	needs, err := db.GetNeedSetupAI()
	if err != nil {
		t.Fatal(err)
	}
	ok := false
	for _, n := range needs {
		t.Logf("Setup := %v", n)
		if n.ID == aiID {
			ok = true
		}
	}
	if !ok {
		t.Fatalf("Missing Commit %d", aiID)
	}
}

func TestUpdateStateAI(t *testing.T) {
	db := mockGameUser()
	aiID := createAI(db, t)
	err := aiID.UpdateState(db, AISetup)
	if err != nil {
		t.Fatal(err)
	}
}
func TestReadyStateAI(t *testing.T) {
	db := mockGameUser()
	_ = createAI(db, t)
	err := db.ReadyContianersByCommit("git_addr", "master", "000001")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetupCmd(t *testing.T) {
	a := AINeedSetup{
		ID:     1,
		Github: "wass88/reversi-random",
		Branch: "master",
		Commit: "cccccc",
	}
	conf := AIRunnerConf{
		API: "http://api",
		Dir: "../../.data",
		Cmd: "../../target/container",
	}
	cmd := a.SetupCmd(conf)
	exp := []string{
		"../../target/container",
		"-api", "http://api",
		"-dir", "../../.data",
		"-github", "wass88/reversi-random",
		"-branch", "master",
		"-commit", "cccccc", "setup",
	}
	for i, c := range cmd.Args {
		if c != exp[i] {
			t.Fatalf("Wrong Cmd [%d]  %s != %s \n%v != %v", i, c, exp[i], cmd, exp)
		}
	}
}

func TestKickAI(t *testing.T) {
	db := mockGameUser()
	_ = createAI(db, t)
	db.Config = &Config{}
	db.Config.AIRunner = AIRunnerConf{
		API: "http://api",
		Dir: "../../.data",
		Cmd: "../../target/container",
	}
	_ = createAI(db, t)
	err := db.KickSetupAI()
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
}
