CREATE TABLE user (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    name varchar(255) NOT NULL,
    github_name varchar(255) NOT NULL UNIQUE,
    authority enum("visitor", "user", "admin") NOT NULL
    );
CREATE TABLE game (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    name VARCHAR(255) NOT NULL
    );
CREATE TABLE ai_github (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    game_id int NOT NULL,
    user_id int NOT NULL,
    github VARCHAR(255) NOT NULL,
    branch VARCHAR(255) NOT NULL,
    updating enum("active", "ignore") NOT NULL,
    FOREIGN KEY (game_id) REFERENCES game(id),
    FOREIGN KEY (user_id) REFERENCES user(id)
    );
CREATE TABLE ai (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    state enum("found", "setup", "ready", "failed", "purged") NOT NULL,
    ai_github_id int NOT NULL,
    commit VARCHAR(255) NOT NULL,
    FOREIGN KEY (ai_github_id) REFERENCES ai_github(id)
    );
CREATE TABLE playout (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    state enum("ready", "running", "completed", "failed") NOT NULL,
    game_id int NOT NULL,
    token varchar(255) NOT NULL,
    FOREIGN KEY (game_id) REFERENCES game(id)
    );
CREATE TABLE playout_ai (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    playout_id int NOT NULL,
    ai_id int NOT NULL,
    turn int NOT NULL,
    FOREIGN KEY (ai_id) REFERENCES ai(id),
    FOREIGN KEY (playout_id) REFERENCES playout(id),
    UNIQUE (playout_id, turn)
    );
CREATE TABLE playout_result (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    playout_id int NOT NULL UNIQUE,
    record text NOT NULL,
    exception text NOT NULL,
    FOREIGN KEY (playout_id) REFERENCES playout(id));
CREATE TABLE playout_result_ai (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    turn int NOT NULL,
    playout_id int NOT NULL,
    stderr text NOT NULL,
    result int NOT NULL,
    exception text NOT NULL,
    FOREIGN KEY (playout_id) REFERENCES playout(id),
    UNIQUE (playout_id, turn)
    );
CREATE TABLE rate_ai (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    ai_id int NOT NULL,
    rate double NOT NULL,
    FOREIGN KEY (ai_id) REFERENCES ai(id)
    );