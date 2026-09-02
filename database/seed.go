package database

import (
	"database/sql"
	"fmt"
)

// LLM nonsense(?)
func SeedData(db *sql.DB) error {
	var userCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM user`).Scan(&userCount); err != nil {
		return fmt.Errorf("seed: counting users: %w", err)
	}
	if userCount > 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("seed: begin tx: %w", err)
	}
	defer tx.Rollback()

	insert := func(query string, args ...any) int64 {
		if err != nil {
			return 0
		}
		var res sql.Result
		res, err = tx.Exec(query, args...)
		if err != nil {
			return 0
		}
		var id int64
		id, err = res.LastInsertId()
		return id
	}

	// --- users ---
	alice := insert(`INSERT INTO user (username, email, password_hash) VALUES (?, ?, ?)`,
		"alice", "alice@example.com", "1234")
	bob := insert(`INSERT INTO user (username, email, password_hash) VALUES (?, ?, ?)`,
		"bob", "bob@example.com", "1234")
	carol := insert(`INSERT INTO user (username, email, password_hash) VALUES (?, ?, ?)`,
		"carol", "carol@example.com", "1234")

	insert(`INSERT INTO session (id, user_id, expires_at)
	        VALUES (?, ?, datetime('now', '+7 days'))`,
		"dev-session-alice", alice)

	// --- categories ---
	scifi := insert(`INSERT INTO category (name) VALUES (?)`, "Science Fiction")
	fantasy := insert(`INSERT INTO category (name) VALUES (?)`, "Fantasy")
	classics := insert(`INSERT INTO category (name) VALUES (?)`, "Classics")
	charStudy := insert(`INSERT INTO category (name) VALUES (?)`, "Character Studies")

	// --- posts ---
	p1 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		alice, "Thoughts on Dune, Chapter 5",
		"Paul's arc in this chapter is incredible. The foreshadowing lands so much harder on a reread.")
	p2 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		bob, "Best fantasy series to start with?",
		"New to the genre and looking for a first series that isn't 14 books long. Suggestions?")
	p3 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		carol, "Revisiting Pride and Prejudice",
		"It holds up remarkably well. Curious what everyone thinks of the pacing in the first half.")

	// --- post <-> category links (p1 gets two, to exercise the many-to-many) ---
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p1, scifi)
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p1, charStudy)
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p2, fantasy)
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p3, classics)

	// --- comments ---
	// Top-level comments have parent_comment_id = NULL (pass nil).
	c1 := insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		bob, p1, nil, "Totally agree. The Gom Jabbar scene reads completely differently the second time.")
	c2 := insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		carol, p1, nil, "I still find the pacing here a little slow, honestly.")
	// A REPLY to c1: parent_comment_id is set, but it's still anchored to post p1.
	c3 := insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		alice, p1, c1, "Right? And it sets up his whole relationship with fear later on.")
	// A top-level comment on a different post (id not needed later).
	insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		alice, p2, nil, "Start with the first Earthsea book — self-contained and short.")

	// --- post likes/dislikes (unique per user+post) ---
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, bob, p1, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, carol, p1, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, alice, p2, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, alice, p3, -1)

	// --- comment likes/dislikes (unique per user+comment) ---
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, alice, c1, 1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, carol, c2, -1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, bob, c3, 1)

	if err != nil {
		return fmt.Errorf("seed: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("seed: commit: %w", err)
	}
	return nil
}
