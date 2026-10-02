package database

import "database/sql"

func InitTable(db *sql.DB) (sql.Result, error) {
	sql := `
CREATE TABLE IF NOT EXISTS user (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    username      TEXT NOT NULL UNIQUE,
    email         TEXT NOT NULL UNIQUE,
    name          TEXT,
    bio           TEXT,
    imagepath     TEXT,
    password_hash TEXT NOT NULL,
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
 
CREATE TABLE IF NOT EXISTS session (
    uuid       TEXT PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);
 
CREATE TABLE IF NOT EXISTS post (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER NOT NULL,
    title      TEXT NOT NULL,
    content    TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);
 
CREATE TABLE IF NOT EXISTS comment (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id           INTEGER NOT NULL,
    parent_post_id    INTEGER NOT NULL,
    parent_comment_id INTEGER,
    content           TEXT NOT NULL,
    created_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME,
    deleted_at        DATETIME,
    FOREIGN KEY (user_id)           REFERENCES user(id)    ON DELETE CASCADE,
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
    is_main     INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id)     REFERENCES post(id)     ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE
);
 
CREATE TABLE IF NOT EXISTS post_like (
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    value   INTEGER NOT NULL,
    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE
);
 
CREATE TABLE IF NOT EXISTS comment_like (
    user_id    INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    value      INTEGER NOT NULL,
    PRIMARY KEY (user_id, comment_id),
    FOREIGN KEY (user_id)    REFERENCES user(id)    ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comment(id) ON DELETE CASCADE
);



CREATE VIRTUAL TABLE IF NOT EXISTS post_fts USING fts5(
    title,
    content,
    content='post',
    content_rowid='id'
);

CREATE VIRTUAL TABLE IF NOT EXISTS comment_fts USING fts5(
    content,
    content='comment',
    content_rowid='id'
);



CREATE TRIGGER IF NOT EXISTS post_ai AFTER INSERT ON post BEGIN
    INSERT INTO post_fts(rowid, title, content)
    VALUES (new.id, new.title, new.content);
END;

CREATE TRIGGER IF NOT EXISTS post_ad AFTER DELETE ON post BEGIN
    INSERT INTO post_fts(post_fts, rowid, title, content)
    VALUES ('delete', old.id, old.title, old.content);
END;

CREATE TRIGGER IF NOT EXISTS post_au AFTER UPDATE ON post BEGIN
    INSERT INTO post_fts(post_fts, rowid, title, content)
    VALUES ('delete', old.id, old.title, old.content);
    
    INSERT INTO post_fts(rowid, title, content)
    VALUES (new.id, new.title, new.content);
END;


CREATE TRIGGER IF NOT EXISTS comment_ai AFTER INSERT ON comment BEGIN
    INSERT INTO comment_fts(rowid, content)
    VALUES (new.id, new.content);
END;

CREATE TRIGGER IF NOT EXISTS comment_ad AFTER DELETE ON comment BEGIN
    INSERT INTO comment_fts(comment_fts, rowid, content)
    VALUES ('delete', old.id, old.content);
END;

CREATE TRIGGER IF NOT EXISTS comment_au AFTER UPDATE ON comment BEGIN
    INSERT INTO comment_fts(comment_fts, rowid, content)
    VALUES ('delete', old.id, old.content);
    
    INSERT INTO comment_fts(rowid, content)
    VALUES (new.id, new.content);
END;`

	return db.Exec(sql)
}
