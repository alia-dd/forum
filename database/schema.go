package database

import "database/sql"

func InitTable(db *sql.DB) (sql.Result, error) {
	sql := `
CREATE TABLE IF NOT EXISTS users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    username      TEXT NOT NULL UNIQUE,
    name          TEXT,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME
);
 
CREATE TABLE IF NOT EXISTS session (
    id         TEXT PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
 
CREATE TABLE IF NOT EXISTS post (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER NOT NULL,
    title      TEXT NOT NULL,
    content    TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
 
CREATE TABLE IF NOT EXISTS comment (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id           INTEGER NOT NULL,
    parent_post_id    INTEGER NOT NULL,
    parent_comment_id INTEGER,
    content           TEXT NOT NULL,
    created_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME,
    FOREIGN KEY (user_id)           REFERENCES users(id)    ON DELETE CASCADE,
    FOREIGN KEY (parent_post_id)    REFERENCES post(id)    ON DELETE CASCADE,
    FOREIGN KEY (parent_comment_id) REFERENCES comment(id) ON DELETE CASCADE
);
 
CREATE TABLE IF NOT EXISTS category (
    id   	   INTEGER PRIMARY KEY AUTOINCREMENT,
    name 	   TEXT NOT NULL UNIQUE,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
 
CREATE TABLE IF NOT EXISTS post_category (
    post_id     INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id)     REFERENCES post(id)     ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE
);
 
CREATE TABLE IF NOT EXISTS post_like (
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    value   INTEGER NOT NULL,
    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE
);
 
CREATE TABLE IF NOT EXISTS comment_like (
    user_id    INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    value      INTEGER NOT NULL,
    PRIMARY KEY (user_id, comment_id),
    FOREIGN KEY (user_id)    REFERENCES users(id)    ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comment(id) ON DELETE CASCADE
);`

	return db.Exec(sql)
}
