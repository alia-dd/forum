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
	dave := insert(`INSERT INTO user (username, email, password_hash) VALUES (?, ?, ?)`,
		"dave", "dave@example.com", "1234")
	erin := insert(`INSERT INTO user (username, email, password_hash) VALUES (?, ?, ?)`,
		"erin", "erin@example.com", "1234")
	frank := insert(`INSERT INTO user (username, email, password_hash) VALUES (?, ?, ?)`,
		"frank", "frank@example.com", "1234")
	grace := insert(`INSERT INTO user (username, email, password_hash) VALUES (?, ?, ?)`,
		"grace", "grace@example.com", "1234")
	heidi := insert(`INSERT INTO user (username, email, password_hash) VALUES (?, ?, ?)`,
		"heidi", "heidi@example.com", "1234")

	insert(`INSERT INTO session (id, user_id, expires_at)
	        VALUES (?, ?, datetime('now', '+7 days'))`,
		"dev-session-alice", alice)

	// --- categories ---
	scifi := insert(`INSERT INTO category (name) VALUES (?)`, "Science Fiction")
	fantasy := insert(`INSERT INTO category (name) VALUES (?)`, "Fantasy")
	classics := insert(`INSERT INTO category (name) VALUES (?)`, "Classics")
	charStudy := insert(`INSERT INTO category (name) VALUES (?)`, "Character Studies")
	mystery := insert(`INSERT INTO category (name) VALUES (?)`, "Mystery & Thriller")
	nonfiction := insert(`INSERT INTO category (name) VALUES (?)`, "Non-Fiction")
	poetry := insert(`INSERT INTO category (name) VALUES (?)`, "Poetry")
	authorItv := insert(`INSERT INTO category (name) VALUES (?)`, "Author Interviews")
	_ = insert(`INSERT INTO category (name) VALUES (?)`, "Book Recommendations")

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
	p4 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		dave, "Reread value", "Some books only click on the second pass. Which ones did that for you?")
	p5 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		frank,
		"On the quiet radicalism of slow books: why I've stopped apologising for taking three months to finish a single novel, and what the club taught me about reading without a finish line",
		"I used to treat a to-read pile like a backlog to clear.\n\nThis year I read four books. Four. And I remember all of them — the shape of the sentences, where I was sitting, what I argued about here afterwards.\n\nSlow reading isn't a productivity failure. It's the whole point. Curious whether anyone else made the same shift.")
	p6 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		erin, "The Name of the Wind — gorgeous prose or overrated?",
		"Rothfuss can write a sentence, no question. But is the lyricism carrying a story that's actually a bit thin? Convince me either way.")
	p7 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		grace, "True-crime that reads like fiction (without exploiting anyone)",
		"Looking for well-reported non-fiction with real narrative craft. In Cold Blood is the obvious one — what else belongs on that shelf?")
	p8 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		heidi, "A poem that changed how you read everything after it",
		"For me it was Mary Oliver's 'Wild Geese'. Share yours.")
	p9 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		alice, "Roundup: what authors say about their terrible first drafts",
		"Collecting the best 'my first draft was garbage' quotes from interviews. Post your favourites and I'll compile them.")
	p10 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		bob, "Unpopular opinion: sometimes the film really is better",
		"Fight me. I'll start: the adaptation tightened a saggy middle act and the book knows it.")

	// --- post <-> category links (p1 gets two, to exercise the many-to-many) ---
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p1, scifi)
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p1, charStudy)
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p2, fantasy)
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p3, classics)
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p4, classics)
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p5, scifi)
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p5, charStudy)
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p5, nonfiction)
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p6, fantasy)
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p7, mystery)
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p7, nonfiction)
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p8, poetry)
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p9, authorItv)

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
	cc1 := insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		bob, p6, nil, "Overrated. The framing device is doing all the heavy lifting.")
	cc2 := insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		carol, p6, cc1, "Hard disagree — the framing IS the story. Kvothe's unreliability is the point.")
	cc3 := insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		dave, p6, cc2, "This is the reply-to-a-reply case — good for testing indentation depth 3.")
	insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		erin, p6, nil, "Second top-level comment on the same post, to test sibling threads.")
	insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		frank, p7, nil, "Say Nothing by Patrick Radden Keefe. Reads like a thriller, deeply reported.")
	insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		grace, p7, nil, "The Devil in the White City, obviously.")

	// --- post likes/dislikes (unique per user+post) ---
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, bob, p1, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, carol, p1, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, alice, p2, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, alice, p3, -1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, alice, p6, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, bob, p6, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, carol, p6, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, dave, p6, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, erin, p6, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, frank, p6, -1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, alice, p10, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, grace, p10, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, bob, p10, -1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, carol, p10, -1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, dave, p10, -1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, heidi, p7, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, alice, p7, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, erin, p4, 1)

	// --- comment likes/dislikes (unique per user+comment) ---
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, alice, c1, 1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, carol, c2, -1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, bob, c3, 1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, dave, cc1, 1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, erin, cc1, 1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, alice, cc1, -1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, bob, cc2, 1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, frank, cc3, -1)

	if err != nil {
		return fmt.Errorf("seed: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("seed: commit: %w", err)
	}
	return nil
}
