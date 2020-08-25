package server

func getDB() *DB {
	dbname := `root:goodpassword@tcp(127.0.0.1:13306)/dev`
	conf := Config{DBName: dbname}
	return conf.NewDB()
}

func mockGameUser() *DB {
	db := getDB()
	dropMock(db)
	setupGame(db)
	setupUser(db)
	return db
}
func mockPlayoutDB() *DB {
	db := getDB()
	dropMock(db)
	setupGame(db)
	setupUser(db)
	setupAI(db)
	setupPlayout(db)
	return db
}
func mockPlayoutResultDB() *DB {
	db := getDB()
	dropMock(db)
	setupGame(db)
	setupUser(db)
	setupAI(db)
	setupPlayout(db)
	setupPlayoutResult(db)
	return db
}

func dropMock(db *DB) {
	db.DB.MustExec(`DELETE FROM playout_result
		WHERE playout_id IN (SELECT id FROM playout WHERE game_id = 1)`)
	db.DB.MustExec(`DELETE FROM playout_result_ai
		WHERE playout_id IN (SELECT id FROM playout WHERE game_id = 1)`)
	db.DB.MustExec(`DELETE FROM playout_ai
	    WHERE playout_id IN
			(SELECT id FROM playout WHERE game_id = 1)
		`)
	db.DB.MustExec(`DELETE FROM playout WHERE game_id = 1`)
	db.DB.MustExec(`DELETE FROM ai
		WHERE ai_github_id
		IN (SELECT id FROM ai_github WHERE game_id = 1)`)
	db.DB.MustExec(`DELETE FROM ai_github
		WHERE game_id = 1`)
	db.DB.MustExec(`DELETE FROM ai_github WHERE id = 2`)
	db.DB.MustExec(`DELETE FROM user WHERE id = 1`)
	db.DB.MustExec(`DELETE FROM game WHERE id = 1`)
}
func setupGame(db *DB) {
	db.DB.MustExec(`
		INSERT INTO game (id, name)
		VALUE (1, "reversi")`)
}
func setupUser(db *DB) {
	db.DB.MustExec(`
		INSERT INTO user (id, name, github_name, authority)
		VALUE (1, "test_user", "test_github", "admin")`)
}

func setupAI(db *DB) {
	db.DB.MustExec(`
		INSERT INTO ai_github (id, game_id, user_id, github, branch)
		VALUE (1, 1, 1, "git_addr", "master")`)
	db.DB.MustExec(`
		INSERT INTO ai_github (id, game_id, user_id, github, branch)
		VALUE (2, 1, 1, "git_addr", "super")`)
	db.DB.MustExec(`
		INSERT INTO ai (id, ai_github_id, commit)
		VALUE (1, 1, "000001")`)
	db.DB.MustExec(`
		INSERT INTO ai (id, ai_github_id, commit)
		VALUE (2, 2, "000002")`)
}

func setupPlayout(db *DB) {
	db.DB.MustExec(`
		INSERT INTO playout (id, state, game_id, token)
		VALUES (1, "ready", 1, "TOKEN")
	`)
	db.DB.MustExec(`
		INSERT INTO playout_ai (id, ai_id, playout_id, turn)
		VALUES (1, 1, 1, 0)
	`)
	db.DB.MustExec(`
		INSERT INTO playout_ai (id, ai_id, playout_id, turn)
		VALUES (2, 2, 1, 1)
	`)
	db.DB.MustExec(`
		INSERT INTO playout (id, state, game_id, token)
		VALUES (2, "completed", 1, "TOKEN2")
	`)
	db.DB.MustExec(`
		INSERT INTO playout_ai (id, ai_id, playout_id, turn)
		VALUES (3, 1, 2, 1)
	`)
	db.DB.MustExec(`
		INSERT INTO playout_ai (id, ai_id, playout_id, turn)
		VALUES (4, 2, 2, 0)
	`)
}

func setupPlayoutResult(db *DB) {
	db.DB.MustExec(`
		INSERT INTO playout_result (playout_id, record, exception)
		VALUES (1, "rec", "exp")
	`)
	db.DB.MustExec(`
		INSERT INTO playout_result_ai (turn, playout_id, result, stderr, exception)
		VALUES (0, 1, 1, "stdout1", "exp1")
	`)
	db.DB.MustExec(`
		INSERT INTO playout_result_ai (turn, playout_id, result, stderr, exception)
		VALUES (1, 1, -1, "stdout1", "exp1")
	`)
}
