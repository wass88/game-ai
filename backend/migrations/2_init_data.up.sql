INSERT INTO game (id, name)
VALUES (1, "reversi"), (2, "game27");

INSERT INTO user (id, name, github_name, authority)
VALUES (1, "wass88", "wass88", "user");

INSERT INTO ai_github (game_id, user_id, github, branch, updating)
VALUES
  (2, 1, "wass88/game-27-ai", "master", "active"),
  (2, 1, "wass88/game-27-ai", "master", "active"),
  (2, 1, "wass88/game-27-python", "master", "active");
