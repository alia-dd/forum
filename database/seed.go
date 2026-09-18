package database

import (
	"database/sql"
	"fmt"
)

// LLM nonsense(?)
func SeedData(db *sql.DB) error {
	var userCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM user`).Scan(&userCount); err != nil {
		return fmt.Errorf("seed: counting user: %w", err)
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

	usersSeedData := []struct {
		username string
		name     string
		email    string
		password string
	}{
		{"alice", "alice", "alice@example.com", "1234abcd"},
		{"bob", "bob", "bob@example.com", "1234abcd"},
		{"carol", "carol", "carol@example.com", "1234abcd"},
		{"dave", "dave", "dave@example.com", "1234abcd"},
		{"erin", "erin", "erin@example.com", "1234abcd"},
		{"frank", "frank", "frank@example.com", "1234abcd"},
		{"grace", "grace", "grace@example.com", "1234abcd"},
		{"heidi", "heidi", "heidi@example.com", "1234abcd"},
	}
	users := make(map[string]int)
	for _, user := range usersSeedData {
		users[user.name] = int(insert(`INSERT INTO user (username, name, email, password_hash) VALUES (?, ?, ?, ?)`, user.username,
			user.name,
			user.email,
			user.password))
	}

	insert(`INSERT INTO session (uuid, user_id, expires_at)
	        VALUES (?, ?, datetime('now', '+7 days'))`,
		"dev-session-alice", users["alice"])

	// --- categories ---
	categorySeedData := []string{"Science Fiction", "Fantasy", "Classics", "Character Studies", "Mystery & Thriller", "Non-Fiction", "Poetry", "Author Interviews", "Book Recommendations"}
	category := make(map[string]int64)
	for _, cat := range categorySeedData {
		category[cat] = insert(`INSERT INTO category (name) VALUES (?)`, cat)
	}

	// --- posts ---
	p1 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		users["alice"], "Thoughts on Dune, Chapter 5",
		"Paul's arc in this chapter is incredible. The foreshadowing lands so much harder on a reread.")
	p2 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		users["bob"], "Best fantasy series to start with?",
		"New to the genre and looking for a first series that isn't 14 books long. Suggestions?")
	p3 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		users["carol"], "Revisiting Pride and Prejudice",
		"It holds up remarkably well. Curious what everyone thinks of the pacing in the first half.")
	p4 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		users["dave"], "Reread value", "Some books only click on the second pass. Which ones did that for you?")
	p5 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		users["frank"],
		"On the quiet radicalism of slow books: why I've stopped apologising for taking three months to finish a single novel, and what the club taught me about reading without a finish line",
		"I used to treat a to-read pile like a backlog to clear.\n\nThis year I read four books. Four. And I remember all of them — the shape of the sentences, where I was sitting, what I argued about here afterwards.\n\nSlow reading isn't a productivity failure. It's the whole point. Curious whether anyone else made the same shift.")
	p6 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		users["erin"], "The Name of the Wind — gorgeous prose or overrated?",
		"Rothfuss can write a sentence, no question. But is the lyricism carrying a story that's actually a bit thin? Convince me either way.")
	p7 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		users["grace"], "True-crime that reads like fiction (without exploiting anyone)",
		"Looking for well-reported non-fiction with real narrative craft. In Cold Blood is the obvious one — what else belongs on that shelf?")
	p8 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		users["heidi"], "A poem that changed how you read everything after it",
		"For me it was Mary Oliver's 'Wild Geese'. Share yours.")
	p9 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		users["alice"], "Roundup: what authors say about their terrible first drafts",
		"Collecting the best 'my first draft was garbage' quotes from interviews. Post your favourites and I'll compile them.")
	p10 := insert(`INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)`,
		users["bob"], "Unpopular opinion: sometimes the film really is better",
		"Fight me. I'll start: the adaptation tightened a saggy middle act and the book knows it.")

	// --- post <-> category links (is_main = 1 marks each post's main category; sides default to 0) ---
	insert(`INSERT INTO post_category (post_id, category_id, is_main) VALUES (?, ?, 1)`, p1, category["Science Fiction"])
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p1, category["Character Studies"])
	insert(`INSERT INTO post_category (post_id, category_id, is_main) VALUES (?, ?, 1)`, p2, category["Fantasy"])
	insert(`INSERT INTO post_category (post_id, category_id, is_main) VALUES (?, ?, 1)`, p3, category["Classics"])
	insert(`INSERT INTO post_category (post_id, category_id, is_main) VALUES (?, ?, 1)`, p4, category["Classics"])
	insert(`INSERT INTO post_category (post_id, category_id, is_main) VALUES (?, ?, 1)`, p5, category["Science Fiction"])
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p5, category["Character Studies"])
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p5, category["Non-Fiction"])
	insert(`INSERT INTO post_category (post_id, category_id, is_main) VALUES (?, ?, 1)`, p6, category["Fantasy"])
	insert(`INSERT INTO post_category (post_id, category_id, is_main) VALUES (?, ?, 1)`, p7, category["Mystery & Thriller"])
	insert(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, p7, category["Non-Fiction"])
	insert(`INSERT INTO post_category (post_id, category_id, is_main) VALUES (?, ?, 1)`, p8, category["Poetry"])
	insert(`INSERT INTO post_category (post_id, category_id, is_main) VALUES (?, ?, 1)`, p9, category["Author Interviews"])
	insert(`INSERT INTO post_category (post_id, category_id, is_main) VALUES (?, ?, 1)`, p10, category["Classics"])

	// --- comments ---
	// Top-level comments have parent_comment_id = NULL (pass nil).
	c1 := insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		users["bob"], p1, nil, "Totally agree. The Gom Jabbar scene reads completely differently the second time.")
	c2 := insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		users["carol"], p1, nil, "I still find the pacing here a little slow, honestly.")
	// A REPLY to c1: parent_comment_id is set, but it's still anchored to post p1.
	c3 := insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		users["alice"], p1, c1, "Right? And it sets up his whole relationship with fear later on.")
	// A top-level comment on a different post (id not needed later).
	insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		users["alice"], p2, nil, "Start with the first Earthsea book — self-contained and short.")
	cc1 := insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		users["bob"], p6, nil, "Overrated. The framing device is doing all the heavy lifting.")
	cc2 := insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		users["carol"], p6, cc1, "Hard disagree — the framing IS the story. Kvothe's unreliability is the point.")
	cc3 := insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		users["dave"], p6, cc2, "This is the reply-to-a-reply case — good for testing indentation depth 3.")
	insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		users["erin"], p6, nil, "Second top-level comment on the same post, to test sibling threads.")
	insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		users["frank"], p7, nil, "Say Nothing by Patrick Radden Keefe. Reads like a thriller, deeply reported.")
	insert(`INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content) VALUES (?, ?, ?, ?)`,
		users["grace"], p7, nil, "The Devil in the White City, obviously.")

	// --- post likes/dislikes (unique per user+post) ---
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["bob"], p1, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["carol"], p1, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["alice"], p2, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["alice"], p3, -1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["alice"], p6, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["bob"], p6, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["carol"], p6, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["dave"], p6, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["erin"], p6, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["frank"], p6, -1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["alice"], p10, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["grace"], p10, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["bob"], p10, -1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["carol"], p10, -1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["dave"], p10, -1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["heidi"], p7, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["alice"], p7, 1)
	insert(`INSERT INTO post_like (user_id, post_id, value) VALUES (?, ?, ?)`, users["erin"], p4, 1)

	// --- comment likes/dislikes (unique per user+comment) ---
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, users["alice"], c1, 1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, users["carol"], c2, -1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, users["bob"], c3, 1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, users["dave"], cc1, 1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, users["erin"], cc1, 1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, users["alice"], cc1, -1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, users["bob"], cc2, 1)
	insert(`INSERT INTO comment_like (user_id, comment_id, value) VALUES (?, ?, ?)`, users["frank"], cc3, -1)

	if err != nil {
		return fmt.Errorf("seed: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("seed: commit: %w", err)
	}

	fmt.Println("seed works")
	return nil
}
