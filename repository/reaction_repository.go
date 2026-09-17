package repository

import (
	"context"
	"database/sql"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
)

const (
	fetchPostReactionQurrey    = ` SELECT  user_id, post_id, value  FROM post_like WHERE post_id = ?`
	fetchCommentReactionQurrey = ` SELECT  user_id, post_id, value  FROM comment_like WHERE post_id = ?`
)

type ReactionRepository struct {
	db *sql.DB
}

// what if i fetch the reaction and save it to a map by postid as key
// then it would be easy to assign a reaction to post/comment
// but that would couse a problem
// if you chang it in post does it to fetch all rection but that might be in the billions

// next idea the boring way get getReaction per post/comment
// will that couse an edge case of some sort?

// when page is loaded eg by goin to the forum or searching a post/comment/username
// it will get each post/comment and add the reaction
// if  the user add a reaction/ remves it,
// it will check if the react i either 1 or -1 if not remove the row

func NewReactionRepositry(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{db: db}
}

func (r *ReactionRepository) GetPostReaction(cx context.Context, PostID int) ([]models.Reaction, error) {
	var postReactions []models.Reaction
	rows, fetchErr := r.db.QueryContext(cx, fetchPostReactionQurrey, PostID)
	if fetchErr != nil {
		// do i want to return server error that look a bit too much for reaction
	}
	for rows.Next() {
		var postReaction models.Reaction
		var value int
		rows.Scan(&postReaction.UseId, &postReaction.ReactionId, &value)
		if value > 0 {
			postReaction.LikeCount = 1
		} else if value < 0 {
			postReaction.LikeCount = 1
		}
		postReactions = append(postReactions, postReaction)
	}
	if err := rows.Err(); err != nil {
		return nil, customerrors.ErrInternalError
	}
	return postReactions, nil

}
func (r *ReactionRepository) GetCommentReaction(cx context.Context, CommentID int) ([]models.Reaction, error) {
	var commentReactions []models.Reaction
	rows, fetchErr := r.db.QueryContext(cx, fetchCommentReactionQurrey, CommentID)
	if fetchErr != nil {
		// do i want to return server error that look a bit too much for reaction
	}
	for rows.Next() {
		var commentReaction models.Reaction
		var value int
		rows.Scan(&commentReaction.UseId, &commentReaction.ReactionId, &value)
		if value > 0 {
			commentReaction.LikeCount = 1
		} else if value < 0 {
			commentReaction.DislikeCount = 1
		}
		commentReactions = append(commentReactions, commentReaction)
	}
	if err := rows.Err(); err != nil {
		return nil, customerrors.ErrInternalError
	}
	return commentReactions, nil
}
